package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
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
		Use https://jwt.io/ to generate token with header and payload above
	*/
	return `eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwczovL2FjY291bnRzLmdvb2dsZS5jb20iLCJzdWIiOiJzdWJqZWN0Iiw` +
		`iYXVkIjoiY2xpZW50X2lkIiwiZXhwIjoxMjE0NzQ4MzY0NywiaWF0IjoxNzIzMjE3MzMyLCJhenAiOiJjbGllbnRfaWQiLCJlbWFpbCI6Im5` +
		`vcmVwbHlAY29tcGFueS5jb20iLCJlbWFpbF92ZXJpZmllZCI6dHJ1ZSwibm9uY2UiOiJub25jZSJ9.YrViMwqc1bEVg1HHOIHPVCI5hVGgri` +
		`DMSXSBBeSj1Kg`
}

func TokenResponse() string {
	return fmt.Sprintf(`{
			"id_token":"%s",
			"scope":[]
		}`, GetToken())
}

func TokenRequest() string {
	return "client_id=client_id&" +
		"client_secret=client_secret&" +
		"code=auth_code&" +
		"grant_type=authorization_code&" +
		"redirect_uri=https%3A%2F%2Flocalhost"
}

func TestSubmitSession(t *testing.T) {
	c := Prepare(t)
	defer Finalize(c)
	t.Run("submit session", func(t *testing.T) {
		// should create separate mocks for each request for each collection,
		// but this works for now
		c.PGMock(PGKeys{"idempotency_token", "session_id", "nonce", "state"},
			PGValues{{"idempotency_token", "session_id", "nonce", "state"}})
		c.PGMock(PGKeys{"session_id", "subject", "email", "id_token", "token"},
			PGValues{{"session_id", "subject", "noreply@company.com", GetToken(), "token"}})

		//nolint:bodyclose
		c.googleOpenIDMock.Add("/token", c.MockHelper(TokenRequest(), TokenResponse(), http.StatusOK))

		body := `{"code":"auth_code"}`
		code, response := c.MakeRequest(http.MethodPost, "/session/submit", &body, nil)

		require.EqualValues(t, 200, code)
		require.JSONEq(t, `{"token":"token"}`, response)
	})
}
