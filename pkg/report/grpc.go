package report

import (
	"context"
	"google.golang.org/grpc"
	"time"

	common "skywalking/network/common/v3"
	event "skywalking/network/event/v3"
)

func ReportEvent(addr string, e *event.Event) (*common.Commands, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), time.Second)
	defer rpcCancel()

	client, err := event.NewEventServiceClient(conn).Collect(rpcCtx)
	if err != nil {
		return nil, err
	}

	if err := client.Send(e); err != nil {
		return nil, err
	}
	return client.CloseAndRecv()
}
