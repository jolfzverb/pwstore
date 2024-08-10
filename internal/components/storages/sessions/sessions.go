package sessions

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/jolfzverb/pwstore/internal/components/postgres"
)

type Session struct {
	SessionID string
	Subject   string
	Email     string
	IDToken   string
	Token     string
}

//go:embed queries/insert_new_session.sql
var insertNewSessionSQL string

//go:embed queries/select_session_by_session_id_and_token.sql
var selectSessionBySessionIDAndTokenSQL string

type Storage struct {
	db *postgres.Postgres
}

func CreateStorage(db *postgres.Postgres) *Storage {
	return &Storage{db}
}

func (s Storage) InsertSession(
	ctx context.Context,
	sessionID string,
	subject string,
	email string,
	idToken string,
) (*Session, error) {
	stmt, err := s.db.PrepareContext(ctx, insertNewSessionSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := Session{}
	err = stmt.QueryRowContext(ctx, sessionID, subject, email, idToken).Scan(
		&session.SessionID,
		&session.Subject,
		&session.Email,
		&session.IDToken,
		&session.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	if session.SessionID != sessionID || session.Subject != subject ||
		session.Email != email || session.IDToken != idToken {
		return nil, fmt.Errorf("session mismatch error")
	}

	return &session, nil
}

func (s Storage) SelectSession(ctx context.Context, sessionID string, token string) (*Session, error) {
	stmt, err := s.db.PrepareContext(ctx, selectSessionBySessionIDAndTokenSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	session := Session{}
	err = stmt.QueryRowContext(ctx, sessionID, token).Scan(
		&session.SessionID,
		&session.Subject,
		&session.Email,
		&session.IDToken,
		&session.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return &session, nil
}
