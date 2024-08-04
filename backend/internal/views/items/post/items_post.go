package itemspost

import (
	"context"
	_ "embed"
	"fmt"

	"backend/internal/api"
	"backend/internal/dependencies"
)

//go:embed queries/insert_items.sql
var insertItemsSQL string

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func PostItems(ctx context.Context, deps dependencies.Collection, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	stmt, err := deps.DB.PrepareContext(ctx, insertItemsSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	newItem := api.Item{}
	err = stmt.QueryRowContext(ctx, request.Body.Name, request.Body.Price).Scan(&newItem.Id, &newItem.Name, &newItem.Price)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	return api.PostItems200JSONResponse(newItem), nil
}
