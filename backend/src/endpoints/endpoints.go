package endpoints

import (
	"backend/src/api"
	"backend/src/dependencies"
	itemsGet "backend/src/views/items/get"
	itemsPost "backend/src/views/items/post"
	"context"
	"encoding/json"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"log"
	"net/http"
	"strings"
)

var (
	server *http.Server
)

type Handlers struct {
	Deps dependencies.Collection
}

func (h Handlers) GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	return itemsGet.GetItems(h.Deps, ctx, request)
}

func (h Handlers) PostItems(ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	return itemsPost.PostItems(h.Deps, ctx, request)
}

func logRequestAndResponse(f strictnethttp.StrictHTTPHandlerFunc, operationID string) strictnethttp.StrictHTTPHandlerFunc {
	var ff = func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (response interface{}, err error) {
		requestJson, _ := json.Marshal(request)
		log.Printf("Start handling %s %s: %s", r.Method, r.URL, string(requestJson))
		result, err := f(ctx, w, r, request)
		var s strings.Builder
		json.NewEncoder(&s).Encode(result)
		log.Printf("Finish handling %s %s: %s", r.Method, r.URL, s.String())
		return result, err
	}
	return ff
}

func InitializeServer() (*http.Server, error) {
	handlers := Handlers{
		Deps: dependencies.GetDependencies(),
	}
	h := api.Handler(api.NewStrictHandler(handlers, []api.StrictMiddlewareFunc{logRequestAndResponse}))

	server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	return server, nil
}
