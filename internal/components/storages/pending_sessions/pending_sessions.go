package pendingsessions

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/go-faster/errors"

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
		return nil, errors.New("failed to prepare statement: " + fmt.Sprint(err))
	}
	defer func() { _ = stmt.Close() }()

	session := PendingSession{} //nolint:exhaustruct
	err = stmt.QueryRowContext(ctx, idempotencyToken).
		Scan(&session.IdempotencyToken, &session.SessionID, &session.Nonce, &session.State)
	if err != nil {
		return nil, errors.New("failed to execute statement: " + fmt.Sprint(err))
	}

	return &session, nil
}

func (s Storage) FetchPendingSession(ctx context.Context, sessionID string) (*PendingSession, error) {
	stmt, err := s.db.PrepareContext(ctx, selectSessionSQL)
	if err != nil {
		return nil, errors.New("failed to prepare statement: " + fmt.Sprint(err))
	}
	defer func() { _ = stmt.Close() }()

	sessions := make([]PendingSession, 0, 1)
	rows, err := stmt.QueryContext(ctx, sessionID)
	if err != nil {
		return nil, errors.New("failed to execute query: " + fmt.Sprint(err))
	}
	for rows.Next() {
		var session PendingSession
		err = rows.Scan(&session.IdempotencyToken, &session.SessionID, &session.Nonce, &session.State)
		if err != nil {
			return nil, errors.New("failed to parse session row: " + fmt.Sprint(err))
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.New("failed to scan sessions: " + fmt.Sprint(err))
	}

	if len(sessions) == 0 {
		return nil, ErrSessionNotFound
	}
	if len(sessions) > 1 {
		return nil, errors.New("multiple sessions found")
	}

	return &sessions[0], nil
}
