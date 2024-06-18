package grpc_send_notify

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type sendNotifyRoutes struct {
	notifyv1.UnimplementedSendNotifyServer
}

func Register(gRPC *grpc.Server) {
	notifyv1.RegisterSendNotifyServer(gRPC, &sendNotifyRoutes{})
}

func (s *sendNotifyRoutes) SendMessage(ctx context.Context, req *notifyv1.SendMessageRequest) (*notifyv1.SendMessageResponse, error) {
	fmt.Println("implement this")
	panic("implement this")
}
