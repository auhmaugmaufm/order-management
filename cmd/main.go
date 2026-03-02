package main

import (
	"fmt"

	"github.com/auhmaugmaufm/event-driven-order/pkg/config"
)

func main() {

	config.Load()
	cfg := config.Get()

	// db := database.NewPostgresDB(cfg)

	fmt.Println("Running on port:", cfg.AppPort)

}
