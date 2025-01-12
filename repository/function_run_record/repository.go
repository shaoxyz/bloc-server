package function_run_record

import (
	"time"

	"github.com/fBloc/bloc-server/aggregate"
	"github.com/fBloc/bloc-server/infrastructure/object_storage"
	"github.com/fBloc/bloc-server/internal/filter_options"
	"github.com/fBloc/bloc-server/pkg/ipt"
	"github.com/fBloc/bloc-server/pkg/value_type"
	"github.com/fBloc/bloc-server/value_object"
)

type FunctionRunRecordRepository interface {
	// Create
	Create(*aggregate.FunctionRunRecord) error

	// Read
	GetByID(id value_object.UUID) (*aggregate.FunctionRunRecord, error)
	Filter(
		filter value_object.RepositoryFilter,
		filterOption value_object.RepositoryFilterOption,
	) ([]*aggregate.FunctionRunRecord, error)
	Count(filter value_object.RepositoryFilter) (int64, error)
	FilterByFilterOption(
		kv map[string]interface{},
		filterOptions *filter_options.FilterOption,
	) ([]*aggregate.FunctionRunRecord, error)
	FilterByFlowRunRecordID(
		FlowRunRecordID value_object.UUID,
	) ([]*aggregate.FunctionRunRecord, error)

	// Update
	PatchProgress(id value_object.UUID, progress float32) error
	PatchProgressMsg(id value_object.UUID, progressMsg string) error
	PatchStageIndex(
		id value_object.UUID, progressStageIndex int,
	) error
	PatchProgressStages(
		id value_object.UUID, progressStages []string,
	) error
	SetTimeout(id value_object.UUID, timeoutTime time.Time) error
	SaveIptBrief(
		id value_object.UUID,
		iptConfig ipt.IptSlice,
		ipts [][]interface{},
		objectStorageImplement object_storage.ObjectStorage,
	) error

	ClearProgress(id value_object.UUID) error
	SaveSuc(
		id value_object.UUID, desc string,
		keyMapValueType map[string]value_type.ValueType,
		keyMapValueIsArray map[string]bool,
		keyMapObjectStorageKey, keyMapBriefData map[string]string,
		intercepted bool,
	) error
	SaveCancel(id value_object.UUID) error
	SaveFail(id value_object.UUID, errMsg string) error

	// Delete
}
