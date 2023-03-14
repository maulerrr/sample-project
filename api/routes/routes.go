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
	postRouter.POST("/", middlewares.AuthMiddleware(), ctrl.AddPost)
	postRouter.GET("/:id", middlewares.AuthMiddleware(), ctrl.GetByPostID)
	postRouter.DELETE("/:user_id/:id", middlewares.AccessMiddleware(), ctrl.DeletePostByID)
	//postRouter.PUT("/update/:id", middlewares.AuthMiddleware(), ctrl.UpdatePostByID)

	postRouter.POST("/like", middlewares.AuthMiddleware(), ctrl.AddLike)
}
