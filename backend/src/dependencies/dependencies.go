package dependencies

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Collection struct {
	DB     *sql.DB
	Config *ConfigType
}

func CreateDependencies(configFile string) (*Collection, error) {
	var err error
	var deps Collection
	deps.Config, err = GetConfig(configFile)
	if err != nil {
		return nil, err
	}

	deps.DB, err = sql.Open("postgres", deps.Config.Database.ConnectionString)
	if err != nil {
		return nil, err
	}
	err = deps.DB.Ping()
	if err != nil {
		return nil, err
	}

	return &deps, nil
}
