package items_post

import "C"
import (
	"backend/src/api"
	"backend/src/dependencies"
	"context"
	_ "embed"
	"log"
)

//go:embed queries/insert_items.sql
var insertItemsSql string

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

func PostItems(deps dependencies.Collection, ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	log.Println("Start handling POST /items")
	stmt, err := deps.Db.Prepare(insertItemsSql)
	if err != nil {
		log.Printf("Failed to prepare statement: %s\n", err.Error())
		return nil, err
	}
	defer stmt.Close()

	newItem := api.Item{}
	err = stmt.QueryRow(request.Body.Name, request.Body.Price).Scan(&newItem.Id, &newItem.Name, &newItem.Price)

	if err != nil {
		log.Printf("Failed to execute statement: %s\n", err.Error())
		return nil, err
	}

	return api.PostItems200JSONResponse(newItem), nil
}
