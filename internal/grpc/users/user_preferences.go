package grpc_users

import (
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"fmt"
	"google.golang.org/grpc"
)

type userRoutes struct {
	notifyv1.UnimplementedUsersServer
}

func Register(gRPC *grpc.Server) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{})
}

func (s *userRoutes) UserPreferences(ctx context.Context, req *notifyv1.UserPreferencesRequest) (*notifyv1.UserPreferencesResponse, error) {
	fmt.Println("implement this")
	panic("implement this")
}
