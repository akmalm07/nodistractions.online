// migrate applies the GORM schema for the API using Neon's direct connection.
package main

import (
	"context"
	"log"

	"nodistractions-online/backend/internal/config"
	"nodistractions-online/backend/internal/db"
)

func main() {
	databaseURL, err := config.MigrationDatabaseURL(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	database, err := db.Open(context.Background(), databaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(database)

	if err := db.Migrate(context.Background(), database); err != nil {
		log.Fatal("migrate database: ", err)
	}

	log.Println("database migration completed")
}
