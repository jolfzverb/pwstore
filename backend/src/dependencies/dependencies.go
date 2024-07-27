package dependencies

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Collection struct {
	Db     *sql.DB
	Config *ConfigType
}

var (
	deps Collection
)

func InitializeDependencies(configFile string) (*Collection, error) {
	var err error
	deps.Config, err = GetConfig(configFile)
	if err != nil {
		return nil, err
	}

	deps.Db, err = sql.Open("postgres", deps.Config.Database.ConnectionString)
	if err != nil {
		return nil, err
	}
	err = deps.Db.Ping()
	if err != nil {
		return nil, err
	}

	return &deps, nil
}

func GetDependencies() Collection {
	return deps
}
