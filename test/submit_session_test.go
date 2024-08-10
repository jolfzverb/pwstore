package tests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func GetToken() string {
	/*
		{
		  "alg": "HS256"
		}
		{
			"iss":"https://accounts.google.com",
			"sub":"subject",
			"aud":"client_id",
			"exp":12147483647,
			"iat":1723217332,
			"azp":"client_id",
			"email":"noreply@company.com",
			"email_verified":true,
			"nonce":"nonce"
		}
		Use https://jwt.io/ to generate token with above header and payload
	*/
	return `eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiJzdWJqZWN0Iiw` +
		`iYXVkIjoiY2xpZW50X2lkIiwiZXhwIjoxMjE0NzQ4MzY0NywiaWF0IjoxNzIzMjE3MzMyLCJhenAiOiJjbGllbnRfaWQiLCJlbWFpbCI6Im5` +
		`vcmVwbHlAY29tcGFueS5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibm9uY2UiOiJub25jZSJ9.YrViMwqc1bEVg1HHOIHPVCI5hVGgri` +
		`DMSXSBBeSj1Kg`
}

func GetTokenResponse() string {
	return `{"id_token":"` + GetToken() + `","scope":[]}`
}

func TestSubmitSession(t *testing.T) {
	c := Prepare(t)
	defer Finalize(c)
	t.Run("submit session", func(t *testing.T) {
		c.pgMock.ExpectPrepare("SELECT idempotency_token, session_id, nonce, state FROM sessions_tmp " +
			"WHERE session_id=\\$1").ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"idempotency_token", "session_id", "nonce", "state"}).
					AddRow("idempotency_token", "session_id", "nonce", "state"))

		c.pgMock.ExpectPrepare("INSERT INTO sessions \\( session_id, subject, email, id_token \\) " +
			"VALUES \\( \\$1, \\$2, \\$3, \\$4 \\) ON CONFLICT \\(session_id\\) DO UPDATE " +
			"SET idempotency_token = \\$1 RETURNING session_id, subject, email, id_token, token").ExpectQuery().
			WillReturnRows(
				sqlmock.NewRows([]string{"session_id", "subject", "email", "id_token", "token"}).
					AddRow("session_id", "subject", "noreply@company.com", GetToken(), "token"))

		c.googleOpenIDMock.Add("/token", func(r *http.Request) (*http.Response, error) {
			bytes, err := io.ReadAll(r.Body)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body")
			}
			require.Equal(t, string(bytes), `client_id=client_id&client_secret=client_secret&code=auth_code&`+
				`grant_type=authorization_code&redirect_uri=https%3A%2F%2Flocalhost`)
			body := strings.NewReader(GetTokenResponse())
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(body),
				Header: http.Header{
					"Content-Type": []string{"application/json"},
				},
			}, nil
		})

		body := `{"code":"auth_code"}`
		request, _ := http.NewRequestWithContext(c.ctx, http.MethodPost, "/session/submit", strings.NewReader(body))
		request.Header.Add("X-Idempotency-Token", "idempotency_token")
		response := httptest.NewRecorder()
		c.handler.ServeHTTP(response, request)

		expectedResponse := `{"token":"token"}`

		require.EqualValues(t, 200, response.Code)
		require.JSONEq(t, expectedResponse, response.Body.String())
	})
}
