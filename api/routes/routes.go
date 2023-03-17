package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/middleware"
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
	postRouter.DELETE("/:user_id/:id", middlewares.PostDeletionMiddleware(), ctrl.DeletePostByID)
	//postRouter.PUT("/update/:id", middlewares.AuthMiddleware(), ctrl.UpdatePostByID)

	postRouter.POST("/like", middlewares.AuthMiddleware(), ctrl.AddLike)
	postRouter.GET("/like/:user_id/:id", middlewares.AuthMiddleware(), ctrl.GetLike)

	commentRouter := router.Group("/comment")
	commentRouter.GET("/:post_id", middlewares.AuthMiddleware(), ctrl.GetAllComments)
	commentRouter.POST("/", middlewares.AuthMiddleware(), ctrl.CreateComment)
	commentRouter.DELETE("/:user_id/:id", middlewares.CommentDeletionMiddleware(), ctrl.DeleteComment)
}
