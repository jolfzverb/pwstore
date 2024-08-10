package postgres

import (
	"database/sql"
	"fmt"

	// register postgres api.
	_ "github.com/lib/pq"

	"github.com/jolfzverb/pwstore/internal/components/secrets"
)

type Postgres = sql.DB

func CreateDB(secret secrets.Model) (*Postgres, error) {
	db, err := sql.Open("postgres", secret.Database.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to db: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("db ping failed: %w", err)
	}

	return db, nil
}
