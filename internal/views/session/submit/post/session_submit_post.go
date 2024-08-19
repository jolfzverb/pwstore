package sessionsubmitpost

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/jolfzverb/pwstore/internal/api"
	googleopenid "github.com/jolfzverb/pwstore/internal/clients/google_open_id"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
	"github.com/jolfzverb/pwstore/internal/contextkey"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

type GoogleOpenIDClaims struct {
	AuthorizedPresenter string `json:"azp"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"` //nolint:tagliatelle
	Nonce               string `json:"nonce"`
	jwt.RegisteredClaims
}

func PostSessionSubmit(
	ctx context.Context,
	request api.PostSessionSubmitRequestObject,
) (api.PostSessionSubmitResponseObject, error) {
	deps := ctx.Value(contextkey.Deps).(*dependencies.Collection)
	session, err := deps.PendingSessionsStorage.FetchPendingSession(ctx, request.Body.SessionId)
	if errors.Is(err, pendingsessions.ErrSessionNotFound) {
		return api.PostSessionSubmit404Response{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	tokenRequest := googleopenid.PostTokenFormdataRequestBody{
		Code:         request.Body.Code,
		ClientId:     deps.Config.OpenIDSettings.ClientID,
		ClientSecret: deps.Secrets.OpenIDSettings.ClientSecret,
		RedirectUri:  deps.Config.OpenIDSettings.RedirectURI,
		GrantType:    deps.Config.OpenIDSettings.GrantType,
	}
	tokenResponse, err := deps.GoogleOpenIDClient.PostTokenWithFormdataBodyWithResponse(ctx, tokenRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to request token: %w", err)
	}
	slog.Debug("Finished request to /token", slog.String("body", string(tokenResponse.Body)))

	if tokenResponse.JSON400 != nil {
		errorDescription := ""
		if tokenResponse.JSON400.ErrorDescription != nil {
			errorDescription = *tokenResponse.JSON400.ErrorDescription
		}
		slog.Warn(fmt.Sprintf("Error on getting token %s: %s", tokenResponse.JSON400.Error, errorDescription))
		return api.PostSessionSubmit401Response{}, nil
	}

	if tokenResponse.JSON200 == nil {
		return nil, fmt.Errorf("token response is not OK: %d!=200, body=%s",
			tokenResponse.StatusCode(), string(tokenResponse.Body))
	}

	idToken := tokenResponse.JSON200.IdToken

	parsedToken, _, err := jwt.NewParser().ParseUnverified(idToken, &GoogleOpenIDClaims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
	tokenClaims, ok := parsedToken.Claims.(*GoogleOpenIDClaims)
	if !ok {
		return nil, fmt.Errorf("failed to extract claims")
	}

	if tokenClaims.Issuer != deps.Config.OpenIDSettings.Issuer {
		return nil, fmt.Errorf("token issuer is invalid: %s", tokenClaims.Issuer)
	}
	if len(tokenClaims.Audience) != 1 || tokenClaims.Audience[0] != deps.Config.OpenIDSettings.ClientID {
		return nil, fmt.Errorf("token audience is invalid: %s", strings.Join(tokenClaims.Audience, ","))
	}
	if tokenClaims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token already expired: %s", tokenClaims.ExpiresAt.String())
	}
	if tokenClaims.Nonce != session.Nonce {
		return nil, fmt.Errorf("token nonce does not match: %s != %s", tokenClaims.Nonce, session.Nonce)
	}
	if !tokenClaims.EmailVerified {
		return nil, fmt.Errorf("email is not verified")
	}

	newSession, err := deps.SessionsStorage.InsertSession(
		ctx, session.SessionID, tokenClaims.Subject, tokenClaims.Email, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return api.PostSessionSubmit200JSONResponse{
		Token: newSession.Token,
	}, nil
}
