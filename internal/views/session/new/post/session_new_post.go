package sessionnewpost

import (
	"context"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

func PostSessionNew(
	ctx context.Context,
	deps dependencies.Collection,
	request api.PostSessionNewRequestObject,
) (api.PostSessionNewResponseObject, error) {
	session, err := deps.PendingSessionsStorage.CreatePendingSession(ctx, request.Params.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return api.PostSessionNew200JSONResponse(api.NewSessionResponse{
		SessionId:             session.SessionID,
		AuthorizationEndpoint: deps.Config.OpenIDSettings.AuthorizationEndpoint,
		ResponseType:          deps.Config.OpenIDSettings.ResponseType,
		ClientId:              deps.Config.OpenIDSettings.ClientID,
		Scope:                 deps.Config.OpenIDSettings.Scope,
		RedirectUri:           deps.Config.OpenIDSettings.RedirectURI,
		State:                 session.State,
		Nonce:                 session.Nonce,
	}), nil
}
