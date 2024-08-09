package sessioninfoget

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

//go:embed queries/select_session_by_token.sql
var selectSessionByTokenSQL string

type sessionInfo struct {
	sessionID string
}

func GetSessionInfo(
	ctx context.Context,
	deps dependencies.Collection,
	request api.GetSessionInfoRequestObject,
) (api.GetSessionInfoResponseObject, error) {
	stmt, err := deps.DB.PrepareContext(ctx, selectSessionByTokenSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	splitToken := strings.Split(request.Params.Authorization, "Bearer ")
	if len(splitToken) != 1 || len(splitToken[0]) == 0 {
		slog.Warn("Invalid token format")
		return api.GetSessionInfo400Response{}, nil
	}

	var session sessionInfo
	err = stmt.QueryRowContext(ctx, splitToken[0]).Scan(&session.sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	response := api.GetSessionInfo200JSONResponse{}
	return response, nil
}
