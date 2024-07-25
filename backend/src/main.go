package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type Item struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var (
	items  = []Item{}
	nextID = 1
	mu     sync.Mutex
)

func main() {

	http.Handle("/", http.NotFoundHandler())

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})
	http.HandleFunc("/items", itemsHandler)
	http.HandleFunc("/items/", itemHandler)
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func itemsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("itemsHandler")
	switch r.Method {
	case http.MethodGet:
		log.Println("GET")
		getItems(w, r)
	case http.MethodPost:
		log.Println("POST")
		addItem(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func itemHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("itemHandler")
	// Add additional handlers for individual items if needed (e.g., PUT, DELETE)
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

func getItems(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func addItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	if err := json.NewDecoder(r.Body).Decode(&newItem); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	mu.Lock()
	newItem.ID = nextID
	nextID++
	items = append(items, newItem)
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newItem)
}
