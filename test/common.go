package tests

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	googleopenid "github.com/jolfzverb/pwstore/internal/clients/google_open_id"
	"github.com/jolfzverb/pwstore/internal/components/config"
	"github.com/jolfzverb/pwstore/internal/components/secrets"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
	"github.com/jolfzverb/pwstore/internal/components/storages/sessions"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/endpoints"
)

type GoogleOpenIDMock struct {
	handlers map[string]func(*http.Request) (*http.Response, error)
}

func (m GoogleOpenIDMock) Do(request *http.Request) (*http.Response, error) {
	path := request.URL.Path
	handler, ok := m.handlers[path]

	if !ok {
		return nil, fmt.Errorf("GoogleOpenIdMock handler is not set for path %s", path)
	}

	return handler(request)
}

func (m GoogleOpenIDMock) Add(path string, handler func(*http.Request) (*http.Response, error)) {
	m.handlers[path] = handler
}

type TestContext struct {
	db               *sql.DB
	handler          http.Handler
	pgMock           sqlmock.Sqlmock
	ctx              context.Context
	googleOpenIDMock *GoogleOpenIDMock
}

func Prepare(t *testing.T) TestContext {
	t.Helper()

	var testContext TestContext
	var err error
	testContext.db, testContext.pgMock, err = sqlmock.New()
	if err != nil {
		t.Errorf("failed to set up SQL mock: %v", err)
	}

	config, err := config.GetConfig("../configs/config_tests.yaml")
	if err != nil {
		t.Errorf("failed to read config file: %v", err)
	}

	secrets, err := secrets.GetConfig("../configs/secrets_tests.yaml")
	if err != nil {
		t.Errorf("failed to read config file: %v", err)
	}
	googleOpenIDMock := &GoogleOpenIDMock{handlers: make(map[string]func(*http.Request) (*http.Response, error))}
	googleOpenIDClient, err := googleopenid.NewClientWithResponses("", googleopenid.WithHTTPClient(*googleOpenIDMock))
	if err != nil {
		t.Errorf("failed to create GoogleOpenIDMock: %v", err)
	}

	deps := dependencies.Collection{
		DB:                     testContext.db,
		Config:                 config,
		Secrets:                secrets,
		PendingSessionsStorage: pendingsessions.CreateStorage(testContext.db),
		SessionsStorage:        sessions.CreateStorage(testContext.db),
		GoogleOpenIDClient:     googleOpenIDClient,
	}

	testContext.handler = endpoints.GetHandler(deps)
	testContext.ctx = context.Background()
	testContext.googleOpenIDMock = googleOpenIDMock
	return testContext
}

func Finalize(c TestContext) {
	c.db.Close()
}
