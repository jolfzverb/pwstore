package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateNewSession(t *testing.T) {
	c := Prepare(t)
	defer Finalize(c)
	t.Run("create new session", func(t *testing.T) {
		// should create separate mocks for each request for each collection,
		// but this works for now
		c.PGMock(PGKeys{"idempotency_token", "session_id", "nonce", "state"},
			PGValues{{"idempotency_token", "session_id", "nonce", "state"}})

		code, response := c.MakeRequest(http.MethodPost, "/session/new", nil,
			&Headers{"Idempotency-Key": "idempotency_token"})

		expectedResponse := `{"session_id":"session_id","authorization_endpoint": "https://accounts.google.com/o/oauth2/v2/auth","response_type": "code",
				"client_id": "client_id",
				"scope": ["openid", "email"],
				"redirect_uri": "https://localhost",
				"state": "state",
				"nonce": "nonce"
			}`

		require.EqualValues(t, 200, code)
		require.JSONEq(t, expectedResponse, response)
	})
}
