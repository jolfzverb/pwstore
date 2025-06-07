package endpoints

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/jolfzverb/pwstore/internal/dependencies"
	"github.com/jolfzverb/pwstore/internal/generated/api"
	sessioninfoget "github.com/jolfzverb/pwstore/internal/views/session/info/get"
	sessionnewpost "github.com/jolfzverb/pwstore/internal/views/session/new/post"
	sessionsubmitpost "github.com/jolfzverb/pwstore/internal/views/session/submit/post"
)

const (
	defaultTimeout = 2 * time.Second
)

func InitializeServer(deps *dependencies.Collection) (*http.Server, error) {
	r := chi.NewRouter()
	r.Use(middleware.StripSlashes)

	api.NewHandler(&sessionsubmitpost.Handler{Deps: deps}, &sessionnewpost.Handler{Deps: deps},
		&sessioninfoget.Handler{Deps: deps},
	).AddRoutes(r)

	server := http.Server{ //nolint:exhaustruct
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: defaultTimeout,
		ReadTimeout:       defaultTimeout,
		WriteTimeout:      defaultTimeout,
	}

	return &server, nil
}
