package sessionsubmitpost

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
)

//go:embed queries/select_new_session.sql
var selectNewSessionSQL string

//go:embed queries/insert_permanent_session.sql
var insertPermanentSession string

type newSession struct {
	idempotencyToken string
	sessionID        string
}

type ololoSession struct {
	token     string
	sessionID string
}

func querySession(ctx context.Context, db *sql.DB, sessionID string) (*newSession, error) {
	stmt, err := db.PrepareContext(ctx, selectNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := newSession{
		sessionID: sessionID,
	}
	err = stmt.QueryRowContext(ctx, sessionID).Scan(&session.sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	return &session, nil
}

func insertSession(ctx context.Context, db *sql.DB, sessionID string) (*ololoSession, error) {
	stmt, err := db.PrepareContext(ctx, selectNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := ololoSession{
		sessionID: sessionID,
	}
	err = stmt.QueryRowContext(ctx, sessionID).Scan(&session.token)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	return &session, nil
}

func PostSessionSubmit(ctx context.Context, deps dependencies.Collection, request api.PostSessionSubmitRequestObject) (api.PostSessionSubmitResponseObject, error) {
	session, err := querySession(ctx, deps.DB, request.Body.SessionId)
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	// query google api for id_token

	newSession, err := insertSession(ctx, deps.DB, session.sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to query session: %w", err)
	}

	return api.PostSessionSubmit200JSONResponse{
		Token: newSession.token,
	}, nil
}
