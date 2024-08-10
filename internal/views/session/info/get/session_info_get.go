package sessioninfoget

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

func GetSessionInfo(
	ctx context.Context,
	deps dependencies.Collection,
	request api.GetSessionInfoRequestObject,
) (api.GetSessionInfoResponseObject, error) {
	splitToken := strings.Split(request.Params.Authorization, "Bearer ")
	if len(splitToken) != 1 || len(splitToken[0]) == 0 {
		slog.Warn("Invalid token format")
		return api.GetSessionInfo400Response{}, nil
	}

	session, err := deps.SessionsStorage.SelectSession(ctx, request.Body.SessionId, splitToken[0])
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	response := api.GetSessionInfo200JSONResponse{
		Email: session.Email,
	}
	return response, nil
}
