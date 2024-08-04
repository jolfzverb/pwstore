package endpoints

import (
	"backend/src/api"
	"backend/src/dependencies"
	"backend/src/views/items/get"
	"backend/src/views/items/post"
	"context"
	"encoding/json"
	"fmt"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"log/slog"
	"net/http"
	"strings"
)

var (
	server *http.Server
)

type Handlers struct {
	deps dependencies.Collection
}

func (h Handlers) GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	return itemsget.GetItems(h.deps, ctx, request)
}

func (h Handlers) PostItems(ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	return itemspost.PostItems(h.deps, ctx, request)
}

func logRequestAndResponse(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	var ff = func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		requestJson, _ := json.Marshal(request)
		slog.Info(fmt.Sprintf("Start handling %s %s", r.Method, r.URL), slog.Any("body", requestJson))
		result, err := f(ctx, w, r, request)
		var s strings.Builder
		json.NewEncoder(&s).Encode(result)
		slog.Info(fmt.Sprintf("Finish handling %s %s", r.Method, r.URL), slog.Any("body", s.String()))
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
		Addr:    ":8080",
		Handler: h,
	}
	return server, nil
}
