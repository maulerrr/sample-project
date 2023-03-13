package db

import (
	"github.com/maulerrr/sample-project/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	dsn := os.Getenv("DATABASE_DSN")

	DB, err = gorm.Open(
		postgres.Open(dsn),
		&gorm.Config{},
	)

	if err != nil {
		log.Fatal(
			"Failed to connect to the database! \n",
			err,
		)
	}

	log.Println("Connected to database!")
	log.Println("Running migrations")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
	)

	if err != nil {
		log.Fatal("Failed to connect to migrate! \n", err)
	}

	log.Println("Migrations done!")
}
