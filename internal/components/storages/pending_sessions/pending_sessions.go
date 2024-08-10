package pendingsessions

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/components/postgres"
)

type PendingSession struct {
	IdempotencyToken string
	SessionID        string
	Nonce            string
	State            string
}

//go:embed queries/insert_new_session.sql
var insertNewSessionSQL string

//go:embed queries/select_session.sql
var selectSessionSQL string

type Storage struct {
	db *postgres.Postgres
}

func CreateStorage(db *postgres.Postgres) *Storage {
	return &Storage{db}
}

func (s Storage) CreatePendingSession(ctx context.Context, idempotencyToken string) (*PendingSession, error) {
	stmt, err := s.db.PrepareContext(ctx, insertNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := PendingSession{}
	err = stmt.QueryRowContext(ctx, idempotencyToken).
		Scan(&session.IdempotencyToken, &session.SessionID, &session.Nonce, &session.State)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &session, nil
}

func (s Storage) FetchPendingSession(ctx context.Context, sessionID string) (*PendingSession, error) {
	stmt, err := s.db.PrepareContext(ctx, selectSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := PendingSession{}
	err = stmt.QueryRowContext(ctx, sessionID).Scan(
		&session.IdempotencyToken,
		&session.SessionID,
		&session.Nonce,
		&session.State)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}
	return &session, nil
}
