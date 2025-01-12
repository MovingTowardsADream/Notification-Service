package gotests

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"Notification_Service/internal/infrastructure/config"
	"Notification_Service/pkg/logger"
)

type TestContextKey string

const TestSessionIDHeader TestContextKey = "Test-Session-Id"

type BaseSuite struct {
	suite.Suite
	Name       string
	ctx        context.Context
	TestCtx    context.Context
	mockServer *MockServer
	repo       Repository
	cancel     func()
}

func (s *BaseSuite) SetupSuite() {
	cfg := config.MustLoad()
	log, err := logger.Setup(cfg.Log.Level, cfg.Log.Path)
	assert.NoError(s.T(), err)

	_ = log
}

func (s *BaseSuite) TearDownSuite() {
	s.cancel()
}

func (s *BaseSuite) GetContext() context.Context {
	return s.ctx
}

// clear all stuck data/caches before next test
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
