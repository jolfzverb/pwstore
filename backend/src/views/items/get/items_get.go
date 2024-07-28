package items_get

import "C"
import (
	"backend/src/api"
	"backend/src/dependencies"
	"context"
	_ "embed"
	"log"
)

//go:embed queries/select_items.sql
var selectItemsSql string

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func GetItems(deps dependencies.Collection, ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	log.Println("Start handling GET /items")
	stmt, err := deps.Db.Prepare(selectItemsSql)
	if err != nil {
		log.Printf("Failed to prepare statement: %s\n", err.Error())
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Printf("Failed to execute query: %s\n", err.Error())
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Name, &item.Price); err != nil {
			log.Printf("Failed to scan item: %s\n", err.Error())
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Failed to scan items: %s\n", err.Error())
		return nil, err
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
