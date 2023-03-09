package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	middlewares "github.com/maulerrr/sample-project/api/middleware"
)

func InitRoutes(app *gin.Engine) {

	router := app.Group("api/v1")

	router.GET("/healthcheck")

	authRouter := router.Group("/auth")
	authRouter.POST("/signup", ctrl.SignUp)
	authRouter.POST("/login", ctrl.Login)

	postRouter := router.Group("/post")
	postRouter.GET("/", middlewares.AuthMiddleware(), ctrl.GetAllPosts)
	postRouter.GET("/:id", middlewares.AuthMiddleware(), ctrl.GetPostByID)
	postRouter.POST("/add", middlewares.AuthMiddleware(), ctrl.AddPost)
	postRouter.DELETE("/delete/:id", middlewares.AuthMiddleware(), ctrl.DeletePostByID)
	postRouter.PUT("/update/:id", middlewares.AuthMiddleware(), ctrl.UpdatePostByID)
}
