package sessions

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/go-faster/errors"

	"github.com/jolfzverb/pwstore/internal/components/postgres"
)

type Session struct {
	SessionID string
	Subject   string
	Email     string
	IDToken   string
	Token     string
}

var ErrSessionNotFound = errors.New("session not found")

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
		return nil, errors.New("failed to prepare statement: " + fmt.Sprint(err))
	}
	defer func() { _ = stmt.Close() }()

	session := Session{} //nolint:exhaustruct
	err = stmt.QueryRowContext(ctx, sessionID, subject, email, idToken).Scan(
		&session.SessionID,
		&session.Subject,
		&session.Email,
		&session.IDToken,
		&session.Token)
	if err != nil {
		return nil, errors.New("failed to execute statement: " + fmt.Sprint(err))
	}

	if session.SessionID != sessionID || session.Subject != subject ||
		session.Email != email || session.IDToken != idToken {
		return nil, errors.New("session mismatch error")
	}

	return &session, nil
}

func (s Storage) SelectSession(ctx context.Context, sessionID string, token string) (*Session, error) {
	stmt, err := s.db.PrepareContext(ctx, selectSessionBySessionIDAndTokenSQL)
	if err != nil {
		return nil, errors.New("failed to prepare statement: " + fmt.Sprint(err))
	}
	defer func() { _ = stmt.Close() }()

	sessions := make([]Session, 0, 1)
	rows, err := stmt.QueryContext(ctx, sessionID, token)
	if err != nil {
		return nil, errors.New("failed to execute query: " + fmt.Sprint(err))
	}
	for rows.Next() {
		var session Session
		err = rows.Scan(&session.SessionID, &session.Subject, &session.Email, &session.IDToken, &session.Token)
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
