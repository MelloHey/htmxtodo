package main

import (
	"fmt"
	"log"

	"github.com/MelloHey/htmxtodo/api"
	"github.com/MelloHey/htmxtodo/store"
)

func main() {
	fmt.Println("Hello ZOZO")

	store, err := store.NewPostgresStore()
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":4000", store)
	server.Run()
}
