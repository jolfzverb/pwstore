package sessionnewpost

import (
	"context"
	"database/sql"
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
	nonce            string
	state            string
}

func createSession(ctx context.Context, db *sql.DB, idempotencyToken string) (*newSession, error) {
	stmt, err := db.PrepareContext(ctx, insertNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := newSession{}
	err = stmt.QueryRowContext(ctx, idempotencyToken).
		Scan(&session.idempotencyToken, &session.sessionID, &session.nonce, &session.state)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &session, nil
}

func PostSessionNew(
	ctx context.Context,
	deps dependencies.Collection,
	request api.PostSessionNewRequestObject,
) (api.PostSessionNewResponseObject, error) {
	session, err := createSession(ctx, deps.DB, request.Params.XIdempotencyToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	return api.PostSessionNew200JSONResponse(api.NewSessionResponse{
		SessionId:             session.sessionID,
		AuthorizationEndpoint: deps.Config.OpenIDSettings.AuthorizationEndpoint,
		ResponseType:          deps.Config.OpenIDSettings.ResponseType,
		ClientId:              deps.Config.OpenIDSettings.ClientID,
		Scope:                 deps.Config.OpenIDSettings.Scope,
		RedirectUri:           deps.Config.OpenIDSettings.RedirectURI,
		State:                 session.state,
		Nonce:                 session.nonce,
	}), nil
}
