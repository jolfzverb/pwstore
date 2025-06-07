package sessioncreatepost

import (
	"context"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/generated/api"
	"github.com/jolfzverb/pwstore/internal/generated/api/models"
)

type Handler struct {
	Deps *dependencies.Collection
}

func (h *Handler) HandleCreate(ctx context.Context, r *models.CreateRequest) (*models.CreateResponse, error) {
	session, err := h.Deps.PendingSessionsStorage.CreatePendingSession(ctx, r.Headers.IdempotencyKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return api.Create200Response(models.NewSessionResponse{
		SessionID:             session.SessionID,
		AuthorizationEndpoint: h.Deps.Config.OpenIDSettings.AuthorizationEndpoint,
		ResponseType:          h.Deps.Config.OpenIDSettings.ResponseType,
		ClientID:              h.Deps.Config.OpenIDSettings.ClientID,
		Scope:                 models.NewSessionResponseScope(h.Deps.Config.OpenIDSettings.Scope),
		RedirectURI:           h.Deps.Config.OpenIDSettings.RedirectURI,
		State:                 session.State,
		Nonce:                 session.Nonce,
	}), nil
}
