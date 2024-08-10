package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestCreateNewSession(t *testing.T) {
	c := Prepare(t)
	defer Finalize(c)
	t.Run("create new session", func(t *testing.T) {
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
				"client_id": "client_id",
				"scope": ["openid", "email"],
				"redirect_uri": "https://localhost",
				"state": "state",
				"nonce": "nonce"
			}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
	})
}
