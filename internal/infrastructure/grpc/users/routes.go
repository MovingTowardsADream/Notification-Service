package users

import (
	"google.golang.org/grpc"

	notifyv1 "Notification_Service/api/gen/go/notify"
)

type userRoutes struct {
	notifyv1.UnimplementedUsersServer
	userInfo UserInfo
}

func Users(gRPC *grpc.Server, userInfo UserInfo) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{userInfo: userInfo})
}
