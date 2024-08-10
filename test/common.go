package tests

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
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
	t                *testing.T
	db               *sql.DB
	handler          http.Handler
	pgMock           sqlmock.Sqlmock
	ctx              context.Context
	googleOpenIDMock *GoogleOpenIDMock
}

type (
	PGKeys   = []string
	PGValues = [][]driver.Value
)

func (c TestContext) PGMock(keys PGKeys, values PGValues) {
	rows := sqlmock.NewRows(keys)
	for i := 0; i < len(values); i++ {
		rows = rows.AddRow(values[i]...)
	}
	c.pgMock.ExpectPrepare(".*").ExpectQuery().
		WillReturnRows(rows)
}

func (c TestContext) MockHelper(requestBody string, responseBody string, status int,
) func(*http.Request) (*http.Response, error) {
	return func(r *http.Request) (*http.Response, error) {
		bytes, err := io.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read request body")
		}
		require.Equal(c.t, string(bytes), requestBody)
		body := strings.NewReader(responseBody)
		return &http.Response{
			StatusCode: status,
			Body:       io.NopCloser(body),
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
		}, nil
	}
}

type Headers = map[string]string

func (c TestContext) MakeRequest(method string, path string, body *string, headers *map[string]string) (int, string) {
	requestBody := "{}"
	if body != nil {
		requestBody = *body
	}
	request, _ := http.NewRequestWithContext(c.ctx, method, path, strings.NewReader(requestBody))
	if headers != nil {
		for k, v := range *headers {
			request.Header.Add(k, v)
		}
	}
	recorder := httptest.NewRecorder()
	c.handler.ServeHTTP(recorder, request)
	response := recorder.Result()
	defer response.Body.Close()
	bytes, err := io.ReadAll(response.Body)
	if err != nil {
		c.t.Errorf("Failed to read response body: %v", err)
	}
	return response.StatusCode, string(bytes)
}

func Prepare(t *testing.T) TestContext {
	t.Helper()

	var testContext TestContext
	testContext.t = t

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
