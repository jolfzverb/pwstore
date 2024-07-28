package endpoints

import (
	"backend/src/api"
	"backend/src/dependencies"
	itemsGet "backend/src/views/items/get"
	itemsPost "backend/src/views/items/post"
	"context"
	"net/http"
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

func InitializeServer() (*http.Server, error) {
	handlers := Handlers{
		Deps: dependencies.GetDependencies(),
	}
	h := api.Handler(api.NewStrictHandler(handlers, nil))

	server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	return server, nil
}
