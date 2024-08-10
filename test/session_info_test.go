package tests

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSessionInfo(t *testing.T) {
	c := Prepare(t)
	defer Finalize(c)
	t.Run("session info", func(t *testing.T) {
		// should create separate mocks for each request for each collection,
		// but this works for now
		c.PGMock(PGKeys{"session_id", "subject", "email", "id_token", "token"},
			PGValues{{"session_id", "subject", "noreply@company.com", "id_token", "token"}})

		code, response := c.MakeRequest(http.MethodGet, "/session/info?session_id=session_id", nil,
			&Headers{"Authorization": "Bearer token"})

		require.EqualValues(t, 200, code)
		require.JSONEq(t, `{"email":"noreply@company.com"}`, response)
	})
}
