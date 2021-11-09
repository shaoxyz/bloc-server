package flow_run_record

import (
	"github.com/fBloc/bloc-backend-go/aggregate"
	"github.com/fBloc/bloc-backend-go/services/flow_run_record"
	"github.com/fBloc/bloc-backend-go/value_object"

	"github.com/google/uuid"
)

var fFRService *flow_run_record.FlowRunRecordService

func InjectFlowRunRecordService(
	f *flow_run_record.FlowRunRecordService,
) {
	fFRService = f
}

type FlowFunctionRecord struct {
	ID                           uuid.UUID                            `json:"id"`
	ArrangementID                uuid.UUID                            `json:"arrangement_id"`
	ArrangementFlowID            string                               `json:"arrangement_flow_id"`
	ArrangementRunRecordID       string                               `json:"arrangement_run_record_id"`
	FlowID                       uuid.UUID                            `json:"flow_id"`
	FlowOriginID                 uuid.UUID                            `json:"flow_origin_id"`
	FlowFuncIDMapFuncRunRecordID map[string]uuid.UUID                 `json:"flowFunctionID_map_functionRunRecordID"`
	TriggerType                  value_object.TriggerType             `json:"trigger_type"`
	TriggerKey                   string                               `json:"trigger_key"`
	TriggerSource                value_object.FlowTriggeredSourceType `json:"trigger_source"`
	TriggerUserName              string                               `json:"trigger_user_name"`
	TriggerTime                  value_object.JsonDate                `json:"trigger_time"`
	StartTime                    value_object.JsonDate                `json:"start_time"`
	EndTime                      value_object.JsonDate                `json:"end_time"`
	Status                       value_object.RunState                `json:"status"`
	ErrorMsg                     string                               `json:"error_msg"`
	RetriedAmount                uint16                               `json:"retried_amount"`
	TimeoutCanceled              bool                                 `json:"timeout_canceled"`
	Canceled                     bool                                 `json:"canceled"`
	CancelUserName               string                               `json:"cancel_user_name"`
}

func fromAgg(
	aggFRR *aggregate.FlowRunRecord,
) *FlowFunctionRecord {
	if aggFRR.IsZero() {
		return nil
	}
	retFlow := &FlowFunctionRecord{
		ID:                           aggFRR.ID,
		ArrangementID:                aggFRR.ArrangementID,
		ArrangementFlowID:            aggFRR.ArrangementFlowID,
		ArrangementRunRecordID:       aggFRR.ArrangementRunRecordID,
		FlowID:                       aggFRR.FlowID,
		FlowOriginID:                 aggFRR.FlowOriginID,
		FlowFuncIDMapFuncRunRecordID: aggFRR.FlowFuncIDMapFuncRunRecordID,
		TriggerType:                  aggFRR.TriggerType,
		TriggerKey:                   aggFRR.TriggerKey,
		TriggerSource:                aggFRR.TriggerSource,
		TriggerTime:                  value_object.NewJsonDate(aggFRR.TriggerTime),
		StartTime:                    value_object.NewJsonDate(aggFRR.StartTime),
		EndTime:                      value_object.NewJsonDate(aggFRR.EndTime),
		Status:                       aggFRR.Status,
		ErrorMsg:                     aggFRR.ErrorMsg,
		RetriedAmount:                aggFRR.RetriedAmount,
		TimeoutCanceled:              aggFRR.TimeoutCanceled,
		Canceled:                     aggFRR.Canceled,
	}
	if aggFRR.CancelUserID == uuid.Nil && aggFRR.TriggerUserID == uuid.Nil {
		return retFlow
	}
	if aggFRR.CancelUserID != uuid.Nil {
		user, _ := fFRService.UserCacheService.GetUserByID(aggFRR.CancelUserID)
		if !user.IsZero() {
			retFlow.CancelUserName = user.Name
		}
	}
	if aggFRR.TriggerUserID != uuid.Nil {
		user, _ := fFRService.UserCacheService.GetUserByID(aggFRR.TriggerUserID)
		if !user.IsZero() {
			retFlow.TriggerUserName = user.Name
		}
	}
	return retFlow
}

func fromAggSlice(
	aggFRRSlice []*aggregate.FlowRunRecord,
) []*FlowFunctionRecord {
	if len(aggFRRSlice) <= 0 {
		return []*FlowFunctionRecord{}
	}
	ret := make([]*FlowFunctionRecord, len(aggFRRSlice))
	for i, j := range aggFRRSlice {
		ret[i] = fromAgg(j)
	}
	return ret
}