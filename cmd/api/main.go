package main

import (
	"fmt"

	"github.com/MotiurRahmanSany/student-management-api/internal/config"
	"github.com/MotiurRahmanSany/student-management-api/internal/database"
)

func main() {
	config := config.GetConfig()
	pool, err := database.NewConnection(config.Db)
	if err != nil {
		fmt.Printf("Error connecting to database: %v\n", err)
		return
	} 
	defer pool.Close()
	
	

	fmt.Printf("Starting %s version %s on port %d\n", config.AppName, config.Version, config.HttpPort)
	fmt.Printf("Database host: %s, port: %d, user: %s\n", config.Db.Host, config.Db.Port, config.Db.User)

	
}
