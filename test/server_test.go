package tests

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jolfzverb/pwstore/internal/components/secrets"

	"github.com/jolfzverb/pwstore/internal/components/secrets"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/jolfzverb/pwstore/internal/components/config"
	"github.com/jolfzverb/pwstore/internal/components/secrets"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/endpoints"
)

type TestContext struct {
	db      *sql.DB
	handler http.Handler
	pgMock  sqlmock.Sqlmock
	ctx     context.Context
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

	deps := dependencies.Collection{
		DB:                     testContext.db,
		Config:                 config,
		Secrets:                secrets,
		PendingSessionsStorage: pendingsessions.CreateStorage(testContext.db),
	}

	testContext.handler = endpoints.GetHandler(deps)
	testContext.ctx = context.Background()
	return testContext
}

func Finalize(c TestContext) {
	c.db.Close()
}

func TestSessionNew(t *testing.T) {
	t.Run("create new session", func(t *testing.T) {
		c := Prepare(t)
		defer Finalize(c)

		c.pgMock.ExpectPrepare("INSERT INTO pending_sessions \\( idempotency_token \\) VALUES \\( \\$1 \\) " +
			"ON CONFLICT \\(idempotency_token\\) DO UPDATE SET idempotency_token = \\$1 RETURNING idempotency_token, " +
			"session_id, nonce, state").ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"idempotency_token", "session_id", "nonce", "state"}).
					AddRow("idempotency_token", "session_id", "nonce", "state"))

		request, _ := http.NewRequestWithContext(c.ctx, http.MethodPost, "/session/new", nil)
		request.Header.Add("X-Idempotency-Token", "idempotency_token")
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{
				"session_id":"session_id",
				"authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth",
				"response_type": "code",
				"client_id": "Google OpenID client_id",
				"scope": ["openid", "email"],
				"redirect_uri": "https://localhost",
				"state": "state",
				"nonce": "nonce"
			}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
	})
}
