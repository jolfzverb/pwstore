package pendingsessions

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/components/postgres"
)

type PendingSession struct {
	IdempotencyToken string
	SessionID        string
	Nonce            string
	State            string
}

var ErrSessionNotFound = errors.New("session not found")

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

	sessions := make([]PendingSession, 0, 1)
	rows, err := stmt.QueryContext(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	for rows.Next() {
		var session PendingSession
		err = rows.Scan(&session.IdempotencyToken, &session.SessionID, &session.Nonce, &session.State)
		if err != nil {
			return nil, fmt.Errorf("failed to parse session row: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan sessions: %w", err)
	}

	if len(sessions) == 0 {
		return nil, ErrSessionNotFound
	}
	if len(sessions) > 1 {
		return nil, fmt.Errorf("multiple sessions found")
	}

	return &sessions[0], nil
}
