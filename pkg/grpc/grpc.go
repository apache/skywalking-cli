package grpc

import (
	"context"
	"time"

	common "skywalking/network/common/v3"
	event "skywalking/network/event/v3"

	"google.golang.org/grpc"
)

func ReportEvent(addr string, e *event.Event) (*common.Commands, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rpcCtx, rpcCancel := context.WithTimeout(context.Background(), time.Second)
	defer rpcCancel()

	client := event.NewEventServiceClient(conn)
	stream, err := client.Collect(rpcCtx)
	if err != nil {
		return nil, err
	}

	if err := stream.Send(e); err != nil {
		return nil, err
	}
	return stream.CloseAndRecv()
}
