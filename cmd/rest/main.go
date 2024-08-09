package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/endpoints"
)

func main() {
	configFile := flag.String("config", "configs/config_local.yaml", "path to config")
	secretsFile := flag.String("secrets", "configs/secrets_local.yaml", "path to secrets")
	flag.Parse()

	deps, err := dependencies.CreateDependencies(*configFile, *secretsFile)
	if err != nil {
		slog.Error("Failed to initialize dependencies", slog.Any("error", err))
		os.Exit(1)
	}

	server, err := endpoints.InitializeServer(*deps)
	if err != nil {
		slog.Error("Failed to initialize server", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("Starting server at :8080")
	err = server.ListenAndServe()
	slog.Error("Server stopped", slog.Any("error", err))
	os.Exit(1)
}
