package sessioninfoget

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

func GetSessionInfo(
	ctx context.Context,
	deps dependencies.Collection,
	request api.GetSessionInfoRequestObject,
) (api.GetSessionInfoResponseObject, error) {
	if len(request.Params.Authorization) <= len("Bearer ") {
		slog.Warn("Invalid token format")
		return api.GetSessionInfo400Response{}, nil
	}
	token := request.Params.Authorization[len("Bearer "):]
	if len(token) == 0 {
		slog.Warn("Invalid token format")
		return api.GetSessionInfo400Response{}, nil
	}
	if len(request.Params.SessionId) == 0 {
		slog.Warn("Invalid session_id format")
		return api.GetSessionInfo400Response{}, nil
	}

	session, err := deps.SessionsStorage.SelectSession(ctx, request.Params.SessionId, token)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	response := api.GetSessionInfo200JSONResponse{
		Email: session.Email,
	}
	return response, nil
}
