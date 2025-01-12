package gotests

import (
	"context"
)

func SetupMocks(ctx context.Context, name string) (
	context.Context,
	*MockServer,
	Repository,
	func(),
) {
	return context.Background(), nil, nil, func() {}
}
