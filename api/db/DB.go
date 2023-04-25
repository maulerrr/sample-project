package db

import (
	"github.com/maulerrr/sample-project/api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func ConnectDB() {
	var err error
	//dsn := os.Getenv("DATABASE_DSN")
	dsn := "postgres://postgres:1111@localhost:5432/sample?sslmode=disable"

	//log.Print(dsn)

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

	//log.Println("Connected to database!")
	//log.Println("Running migrations")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Like{},
		&models.Comment{},
	)

	if err != nil {
		log.Fatal("Failed to connect to migrate! \n", err)
	}

	//log.Println("Migrations done!")
}
