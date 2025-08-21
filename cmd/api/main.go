package main

import (
	"github.com/gabrielolivrp/pastebin-api/internal/server"
	"github.com/gabrielolivrp/pastebin-api/pkg/config"
)

func main() {
	config, err := config.LoadConfig(".env")
	if err != nil {
		panic(err)
	}

	if err := server.Start(config); err != nil {
		panic(err)
	}
}
