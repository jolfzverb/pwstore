package main

import (
	"backend/src/api"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v2"

	_ "backend/src/api"
)

type Config struct {
	Database struct {
		ConnectionString string `yaml:"connection_string"`
	} `yaml:"database"`
}

type Item struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}

var (
	db *sql.DB
)

type Server struct{}

func (Server) GetItems(ctx context.Context, request api.GetItemsRequestObject) (api.GetItemsResponseObject, error) {
	rows, err := db.Query("SELECT id, name, price FROM items")
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

func (Server) PostItems(ctx context.Context, request api.PostItemsRequestObject) (api.PostItemsResponseObject, error) {
	newItem := api.Item{
		Name:  request.Body.Name,
		Price: request.Body.Price,
	}
	err := db.QueryRow(
		"INSERT INTO items (name, price) VALUES ($1, $2) RETURNING id",
		request.Body.Name, request.Body.Price).Scan(&newItem.Id)

	if err != nil {
		return nil, err
	}
	return api.PostItems200JSONResponse(newItem), nil
}

func main() {
	var err error
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	db, err = sql.Open("postgres", config.Database.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	server := Server{}
	h := api.Handler(api.NewStrictHandler(server, nil))

	s := &http.Server{
		Addr:    ":8080",
		Handler: h,
	}
	log.Println("Server started at :8080")
	log.Fatal(s.ListenAndServe())
}

func loadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}
