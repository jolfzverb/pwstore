package endpoints

import (
	"backend/src/api"
	itemsGet "backend/src/views/items/get"
	itemsPost "backend/src/views/items/post"
	"context"
	"net/http"
)

var (
	server *http.Server
)

type Handlers struct{}

func (Handlers) GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	return itemsGet.GetItems(ctx, request)
}

func (Handlers) PostItems(ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	return itemsPost.PostItems(ctx, request)
}

func InitializeServer() (*http.Server, error) {
	handlers := Handlers{}
	h := api.Handler(api.NewStrictHandler(handlers, nil))

	server = &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	return server, nil
}
