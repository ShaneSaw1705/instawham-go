package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDb() {
	// Data Source Name (DSN) formatted for PostgreSQL connection
	dsn := os.Getenv("DB_URL")
	// Open a connection to the database using Gorm
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Assign to the package-level DB variable
	// Handle potential errors
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// If the connection is successful, you can use `DB` to interact with your database
	log.Println("Database connected successfully")
}
