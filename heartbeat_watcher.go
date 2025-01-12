package bloc

import (
	"time"

	"github.com/fBloc/bloc-server/aggregate"
	"github.com/fBloc/bloc-server/event"
)

// RePubDeadRuns 重发运行中断的任务
func (blocApp *BlocApp) RePubDeadRuns() {
	logger := blocApp.GetOrCreateConsumerLogger()
	heartBeatRepo := blocApp.GetOrCreateFuncRunHBeatRepository()
	funcRunRecordRepo := blocApp.GetOrCreateFunctionRunRecordRepository()

	ticker := time.NewTicker(aggregate.HeartBeatDeadThreshold)
	defer ticker.Stop()
	for range ticker.C {
		deads, err := heartBeatRepo.AllDeads()
		if err != nil {
			logger.Errorf("heartbeat watcher error: %s", err.Error())
			continue
		}
		if len(deads) <= 0 {
			continue
		}

		for _, d := range deads {
			// 立即进行删除此条信息（利用mongo通过ID删除的原子性保障来确保不会「重复重发」）
			deleteAmount, err := heartBeatRepo.Delete(d.ID)
			if err != nil {
				logger.Errorf("heartBeatRepo.Delete failed, error: %s", err.Error())
				continue
			}
			if deleteAmount != 1 { // 避免并发watch重复发布
				continue
			}

			// 查询对应的function run record是否存在
			funcRunRecord, err := funcRunRecordRepo.GetByID(d.FunctionRunRecordID)
			if err != nil {
				logger.Errorf("funcRunRecordRepo.GetByID failed. error:: %s", err.Error())
				continue
			}
			if funcRunRecord.IsZero() {
				continue
			}
			err = funcRunRecordRepo.ClearProgress(funcRunRecord.ID)
			if err != nil {
				logger.Errorf("funcRunRecordRepo.ClearProgress: %s", err.Error())
			}

			// 再次进行发布
			err = event.PubEvent(&event.FunctionToRun{
				FunctionRunRecordID: funcRunRecord.ID,
			})
			if err != nil {
				logger.Errorf("pub func event failed. error: %s", err.Error())
			} else {
				logger.Infof("re-pub function run record: %s", funcRunRecord.ID)
			}
		}
	}
}
