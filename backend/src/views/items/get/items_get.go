package itemsget

import "C"
import (
	"backend/src/api"
	"backend/src/dependencies"
	"context"
	_ "embed"
	"fmt"
)

//go:embed queries/select_items.sql
var selectItemsSQL string

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func GetItems(deps dependencies.Collection, ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	stmt, err := deps.DB.Prepare(selectItemsSQL)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to scan items: %w", err)
	}
	response := api.GetItems200JSONResponse{}
	for _, item := range items {
		response.Items = append(response.Items, api.Item{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	return response, nil
}
