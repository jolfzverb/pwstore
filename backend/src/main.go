package main

import (
	"backend/src/dependencies"
	"backend/src/endpoints"
	"flag"
	"log"
)

func main() {
	configFile := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	_, err := dependencies.InitializeDependencies(*configFile)
	if err != nil {
		log.Fatal(err)
	}
	server, err := endpoints.InitializeServer()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Server started at :8080")
	log.Fatal(server.ListenAndServe())
}
