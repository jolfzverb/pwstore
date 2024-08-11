package dependencies

import (
	"fmt"

	googleopenid "github.com/jolfzverb/pwstore/internal/clients/google_open_id"
	"github.com/jolfzverb/pwstore/internal/components/config"
	"github.com/jolfzverb/pwstore/internal/components/postgres"
	"github.com/jolfzverb/pwstore/internal/components/secrets"
	pendingsessions "github.com/jolfzverb/pwstore/internal/components/storages/pending_sessions"
	"github.com/jolfzverb/pwstore/internal/components/storages/sessions"
)

type Collection struct {
	DB                     *postgres.Postgres
	Config                 *config.Model
	Secrets                *secrets.Model
	PendingSessionsStorage *pendingsessions.Storage
	SessionsStorage        *sessions.Storage
	GoogleOpenIDClient     *googleopenid.ClientWithResponses
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
	deps.SessionsStorage = sessions.CreateStorage(deps.DB)

	deps.GoogleOpenIDClient, err = googleopenid.NewClientWithResponses(deps.Config.Clients.GoogleOpenID.Host)
	if err != nil {
		return nil, fmt.Errorf("failed to create Google OpenID client: %w", err)
	}

	return &deps, nil
}
