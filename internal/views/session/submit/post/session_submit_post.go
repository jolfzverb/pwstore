package sessionsubmitpost

import (
	"context"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"

	googleopenid "github.com/jolfzverb/pwstore/internal/clients/google_open_id"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/generated/api"
	"github.com/jolfzverb/pwstore/internal/generated/api/models"
)

type GoogleOpenIDClaims struct {
	AuthorizedPresenter string `json:"azp"`
	Email               string `json:"email"`
	EmailVerified       bool   `json:"email_verified"` //nolint:tagliatelle
	Nonce               string `json:"nonce"`
	jwt.RegisteredClaims
}

type Handler struct {
	Deps *dependencies.Collection
}

func (h *Handler) HandleSubmit(ctx context.Context, r *models.SubmitRequest) (*models.SubmitResponse, error) {
	session, err := h.Deps.PendingSessionsStorage.FetchPendingSession(ctx, r.Body.SessionID)
	if errors.Is(err, pendingsessions.ErrSessionNotFound) {
		return api.Submit404Response(), nil
	}
	if err != nil {
		return nil, errors.New("failed to query session: " + fmt.Sprint(err))
	}

	tokenRequest := googleopenid.PostTokenFormdataRequestBody{
		Code:         r.Body.Code,
		ClientId:     h.Deps.Config.OpenIDSettings.ClientID,
		ClientSecret: h.Deps.Secrets.OpenIDSettings.ClientSecret,
		RedirectUri:  h.Deps.Config.OpenIDSettings.RedirectURI,
		GrantType:    h.Deps.Config.OpenIDSettings.GrantType,
	}
	tokenResponse, err := h.Deps.GoogleOpenIDClient.PostTokenWithFormdataBodyWithResponse(ctx, tokenRequest)
	if err != nil {
		return nil, errors.New("failed to request token: " + fmt.Sprint(err))
	}
	slog.Debug("Finished request to /token", slog.String("body", string(tokenResponse.Body)))

	if tokenResponse.JSON400 != nil {
		errorDescription := ""
		if tokenResponse.JSON400.ErrorDescription != nil {
			errorDescription = *tokenResponse.JSON400.ErrorDescription
		}
		slog.Warn(fmt.Sprintf("Error on getting token %s: %s", tokenResponse.JSON400.Error, errorDescription))

		return api.Submit401Response(), nil
	}

	if tokenResponse.JSON200 == nil {
		return nil, errors.New("token response is not OK: " +
			strconv.Itoa(tokenResponse.StatusCode()) + "!=200, body=" + string(tokenResponse.Body))
	}

	idToken := tokenResponse.JSON200.IdToken

	parsedToken, _, err := jwt.NewParser().ParseUnverified(idToken, &GoogleOpenIDClaims{}) //nolint:exhaustruct
	if err != nil {
		return nil, errors.New("failed to parse token: " + fmt.Sprint(err))
	}
	tokenClaims, ok := parsedToken.Claims.(*GoogleOpenIDClaims)
	if !ok {
		return nil, errors.New("failed to extract claims")
	}

	if tokenClaims.Issuer != h.Deps.Config.OpenIDSettings.Issuer {
		return nil, errors.New("token issuer is invalid: " + tokenClaims.Issuer)
	}
	if len(tokenClaims.Audience) != 1 || tokenClaims.Audience[0] != h.Deps.Config.OpenIDSettings.ClientID {
		return nil, errors.New("token audience is invalid: " + strings.Join(tokenClaims.Audience, ","))
	}
	if tokenClaims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token already expired: " + tokenClaims.ExpiresAt.String())
	}
	if tokenClaims.Nonce != session.Nonce {
		return nil, errors.New("token nonce does not match: " + tokenClaims.Nonce + " != " + session.Nonce)
	}
	if !tokenClaims.EmailVerified {
		return nil, errors.New("email is not verified")
	}

	newSession, err := h.Deps.SessionsStorage.InsertSession(
		ctx, session.SessionID, tokenClaims.Subject, tokenClaims.Email, idToken)
	if err != nil {
		return nil, errors.New("failed to create session: " + fmt.Sprint(err))
	}

	return api.Submit200Response(models.SubmitSessionResponse{
		Token: newSession.Token,
	}), nil
}
