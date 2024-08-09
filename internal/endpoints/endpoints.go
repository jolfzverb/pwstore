package endpoints

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	"github.com/jolfzverb/pwstore/internal/api"
	"github.com/jolfzverb/pwstore/internal/dependencies"
	itemsget "github.com/jolfzverb/pwstore/internal/views/items/get"
	itemspost "github.com/jolfzverb/pwstore/internal/views/items/post"
)

var server *http.Server

type Handlers struct {
	deps dependencies.Collection
}

func (h Handlers) GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	return itemsget.GetItems(ctx, h.deps, request)
}

func (h Handlers) PostItems(ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	return itemspost.PostItems(ctx, h.deps, request)
}

func logRequestAndResponse(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	ff := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		requestJSON, err := json.Marshal(request)
		if err == nil {
			slog.Info(fmt.Sprintf("Start handling %s %s", r.Method, r.URL), slog.Any("body", requestJSON), slog.String("operation_id", operationID))
		} else {
			slog.Info(fmt.Sprintf("Start handling %s %s", r.Method, r.URL), slog.String("operation_id", operationID))
		}

		result, err := f(ctx, w, r, request)

		var s strings.Builder
		err2 := json.NewEncoder(&s).Encode(result)
		if err2 == nil {
			slog.Info(fmt.Sprintf("Finish handling %s %s", r.Method, r.URL), slog.Any("body", s.String()), slog.String("operation_id", operationID))
		} else {
			slog.Info(fmt.Sprintf("Finish handling %s %s", r.Method, r.URL), slog.String("operation_id", operationID))
		}

		return result, err
	}
	return ff
}

func GetHandler(deps dependencies.Collection) http.Handler {
	handlers := Handlers{
		deps: deps,
	}
	return api.Handler(api.NewStrictHandler(handlers, []api.StrictMiddlewareFunc{logRequestAndResponse}))
}

func InitializeServer(deps dependencies.Collection) (*http.Server, error) {
	h := GetHandler(deps)

	server = &http.Server{
		Addr:              ":8080",
		Handler:           h,
		ReadHeaderTimeout: 2 * time.Second,
	}
	return server, nil
}
