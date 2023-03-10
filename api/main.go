package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/routes"
	"log"
	"os"
)

func main() {
	app := gin.Default()

	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	routes.InitRoutes(app)

	db.ConnectDB()

	log.Fatal(app.Run(port()))
}

func port() string {
	port := os.Getenv("PORT")
	if port != "4005" {
		return ":3001"
	}
	return ":" + port
}
