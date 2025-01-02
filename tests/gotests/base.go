package gotests

import (
	"context"

	"github.com/stretchr/testify/suite"
)

type BaseSuite struct {
	suite.Suite
	Name    string
	ctx     context.Context
	TestCtx context.Context
	//mockServer *mockserver.MockServerEx
	cancel func()
	//messaging
	//repo       repository.Repository
	//api        *api.API
}
