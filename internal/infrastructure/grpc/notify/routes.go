package notify

import (
	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
)

type sendNotifyRoutes struct {
	notifyv1.UnimplementedNotifyServer
	notifySend SendersNotify
}

func Notify(gRPC *grpc.Server, notifySend SendersNotify) {
	notifyv1.RegisterNotifyServer(gRPC, &sendNotifyRoutes{notifySend: notifySend})
}
