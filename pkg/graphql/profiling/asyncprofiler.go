package profiling

import (
	"github.com/apache/skywalking-cli/assets"
	"github.com/apache/skywalking-cli/pkg/graphql/client"
	"github.com/machinebox/graphql"
	"github.com/urfave/cli/v2"
	api "skywalking.apache.org/repo/goapi/query"
)

func CreateAsyncProfilerTask(ctx *cli.Context, condition *api.AsyncProfilerTaskCreationRequest) (api.AsyncProfilerTaskCreationResult, error) {
	var response map[string]api.AsyncProfilerTaskCreationResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/CreateTask.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerTaskList(ctx *cli.Context, condition *api.AsyncProfilerTaskListRequest) (api.AsyncProfilerTaskListResult, error) {
	var response map[string]api.AsyncProfilerTaskListResult

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetTaskList.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerTaskProgress(ctx *cli.Context, taskID string) (api.AsyncProfilerTaskProgress, error) {
	var response map[string]api.AsyncProfilerTaskProgress

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetTaskProgress.graphql"))
	request.Var("taskId", taskID)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}

func GetAsyncProfilerAnalyze(ctx *cli.Context, condition *api.AsyncProfilerAnalyzationRequest) (api.AsyncProfilerAnalyzation, error) {
	var response map[string]api.AsyncProfilerAnalyzation

	request := graphql.NewRequest(assets.Read("graphqls/profiling/asyncprofiler/GetAnalysis.graphql"))
	request.Var("condition", condition)

	err := client.ExecuteQuery(ctx, request, &response)

	return response["result"], err
}
