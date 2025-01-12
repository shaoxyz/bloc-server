package function_task

import (
	"context"

	"github.com/fBloc/bloc-server/infrastructure/log"
	"github.com/fBloc/bloc-server/repository/function"
	mongoFunc "github.com/fBloc/bloc-server/repository/function/mongo"
	"github.com/fBloc/bloc-server/repository/function_run_record"
	mongoFRR "github.com/fBloc/bloc-server/repository/function_run_record/mongo"
)

type FunctionTaskConfiguration func(fs *FunctionTaskService) error

type FunctionTaskService struct {
	logger             *log.Logger
	Functions          function.FunctionRepository
	FunctionRunRecords function_run_record.FunctionRunRecordRepository
}

func WithLogger(logger *log.Logger) FunctionTaskConfiguration {
	return func(us *FunctionTaskService) error {
		us.logger = logger
		return nil
	}
}

func WithMongoFunctionRepository(hosts []string, port int, db, user, password string) FunctionTaskConfiguration {
	return func(fts *FunctionTaskService) error {
		ur, err := mongoFunc.New(
			context.Background(),
			hosts, port, db, user, password, mongoFunc.DefaultCollectionName,
		)
		if err != nil {
			return err
		}
		fts.Functions = ur
		return nil
	}
}

func WithMongoFunctionRunRecordRepository(
	hosts []string, port int, db, user, password string,
) FunctionTaskConfiguration {
	return func(fts *FunctionTaskService) error {
		ur, err := mongoFRR.New(
			context.Background(),
			hosts, port, db, user, password, mongoFunc.DefaultCollectionName,
		)
		if err != nil {
			return err
		}
		fts.FunctionRunRecords = ur
		return nil
	}
}
