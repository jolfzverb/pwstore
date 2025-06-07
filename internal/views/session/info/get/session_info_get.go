package sessioninfoget

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jolfzverb/pwstore/internal/components/storages/sessions"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/generated/api"
	"github.com/jolfzverb/pwstore/internal/generated/api/models"
)

type Handler struct {
	Deps *dependencies.Collection
}

func (h *Handler) HandleInfo(ctx context.Context, r *models.InfoRequest) (*models.InfoResponse, error) {
	if len(r.Headers.Authorization) <= len("Bearer ") {
		slog.Warn("Invalid token format")
		return api.Info400Response(), nil
	}
	token := (r.Headers.Authorization)[len("Bearer "):]
	if len(token) == 0 {
		slog.Warn("Invalid token format")
		return api.Info400Response(), nil
	}
	if len(r.Query.SessionID) == 0 {
		slog.Warn("Invalid session_id format")
		return api.Info400Response(), nil
	}

	session, err := h.Deps.SessionsStorage.SelectSession(ctx, r.Query.SessionID, token)
	if errors.Is(err, sessions.ErrSessionNotFound) {
		slog.Warn(fmt.Sprintf("Session not found for (session_id, token): %v", err))
		return api.Info401Response(), nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return api.Info200Response(
		models.SessionInfo{
			Email: session.Email,
		},
	), nil
}
