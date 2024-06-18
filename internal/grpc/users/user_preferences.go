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

func Users(gRPC *grpc.Server) {
	notifyv1.RegisterUsersServer(gRPC, &userRoutes{})
}

func (s *userRoutes) UserPreferences(ctx context.Context, req *notifyv1.UserPreferencesRequest) (*notifyv1.UserPreferencesResponse, error) {

	// TODO: Validate request

	fmt.Println(req)

	return nil, nil
}
