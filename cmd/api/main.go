package main

import (
	"fmt"

	"github.com/MotiurRahmanSany/student-management-api/internal/config"
)

func main() {
	fmt.Println("Loading configuration...")
	config := config.GetConfig()
	fmt.Printf("Starting %s version %s on port %d\n", config.AppName, config.Version, config.HttpPort)
	fmt.Printf("Database host: %s, port: %d, user: %s\n", config.Db.Host, config.Db.Port, config.Db.User)
			
	serve(config)
}
