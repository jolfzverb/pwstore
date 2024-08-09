package sessionnewpost

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

//go:embed queries/insert_new_session.sql
var insertNewSessionSQL string

type newSession struct {
	idempotencyToken string
	sessionID        string
}

func PostSessionNew(ctx context.Context, deps dependencies.Collection, request api.PostSessionNewRequestObject) (api.PostSessionNewResponseObject, error) {
	stmt, err := deps.DB.PrepareContext(ctx, insertNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := newSession{
		idempotencyToken: request.Params.XIdempotencyToken,
	}
	err = stmt.QueryRowContext(ctx, request.Params.XIdempotencyToken).Scan(&session.sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return api.PostSessionNew200JSONResponse(api.NewSessionResponse{
		SessionId: session.sessionID,
	}), nil
}
