package items_get

import "C"
import (
	"backend/src/api"
	"backend/src/dependencies"
	"context"
	_ "embed"
)

//go:embed queries/select_items.sql
var selectItemsSql string

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	deps := dependencies.GetDependencies()
	stmt, err := deps.Db.Prepare(selectItemsSql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	response := api.GetItems200JSONResponse{}
	for _, item := range items {
		response = append(response, api.Item{
			Id:    item.ID,
			Name:  item.Name,
			Price: item.Price,
		})
	}
	return response, nil
}
