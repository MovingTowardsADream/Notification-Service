package gotests

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/stretchr/testify/suite"
)

type TestContextKey string

const TestSessionIDHeader TestContextKey = "Test-Session-Id"

type BaseSuite struct {
	suite.Suite
	Name       string
	ctx        context.Context
	TestCtx    context.Context
	Repo       Repository
	mockServer *MockServer
	cancel     func()
}

func (s *BaseSuite) SetupSuite() {
	s.ctx, s.Repo, s.mockServer, s.cancel = SetupMocks(context.Background(), s.Name, s.T())
}

func (s *BaseSuite) TearDownSuite() {
	s.cancel()
}

func (s *BaseSuite) GetContext() context.Context {
	return s.ctx
}

// clear all stuck data/caches before next test.
func (s *BaseSuite) clear() {}

func (s *BaseSuite) NewTestContext() {
	s.clear()
	s.TestCtx = s.newContext(s.ctx)
}

func (s *BaseSuite) newContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, TestSessionIDHeader, s.generateContextName())
}

func (s *BaseSuite) generateContextName() string {
	return fmt.Sprintf("%s-%d-%d", s.Name, time.Now().Nanosecond(), rand.Uint32())
}
