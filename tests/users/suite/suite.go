package suite

import (
	"Notification_Service/configs"
	notifyv1 "Notification_Service/protos/gen/go/notify"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
)

const _default_gRPC_Host = "localhost"

type Suite struct {
	*testing.T
	Cfg         *configs.Config
	UsersClient notifyv1.UsersClient
}

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := configs.MustLoadPath("../../configs/config.yaml", "../../.env")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

	// Called when the context terminates
	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(context.Background(),
		_default_gRPC_Host+cfg.Port,
		grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

	return ctx, &Suite{
		T:           t,
		Cfg:         cfg,
		UsersClient: notifyv1.NewUsersClient(cc),
	}
}
