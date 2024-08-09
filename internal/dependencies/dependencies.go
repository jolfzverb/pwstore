package dependencies

import (
	"fmt"


	"github.com/jolfzverb/pwstore/internal/components/config"
	"github.com/jolfzverb/pwstore/internal/components/postgres"
	"github.com/jolfzverb/pwstore/internal/components/secrets"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
)

type Collection struct {
	DB                     *postgres.Postgres
	Config                 *config.Model
	Secrets                *secrets.Model
	PendingSessionsStorage *pendingsessions.Storage
}

func CreateDependencies(configFile string, secretsFile string) (*Collection, error) {
	var err error
	var deps Collection
	deps.Config, err = config.GetConfig(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get config: %w", err)
	}

	deps.Secrets, err = secrets.GetConfig(secretsFile)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets: %w", err)
	}

	deps.DB, err = postgres.CreateDB(*deps.Secrets)
	if err != nil {
		return nil, fmt.Errorf("failed to create db: %w", err)
	}

	deps.PendingSessionsStorage = pendingsessions.CreateStorage(deps.DB)

	return &deps, nil
}
