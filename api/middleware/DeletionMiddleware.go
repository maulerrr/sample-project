package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"gorm.io/gorm"
	"log"
	"strconv"
	"strings"
)

func PostDeletionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		post := GetByPostIDAndUserID(context)

		if authHeader == "" {
			utils.SendMessageWithStatus(context, "Authorize!", 404)
			context.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.SendMessageWithStatus(context, "Incorrect authorization header", 403)
			context.Abort()
			return
		}

		token := headerParts[1]

		user, err := parseToken(token)

		if user.ID != post.UserID {
			utils.SendMessageWithStatus(context, "U are not allowed to delete someone's post!", 403)
			context.Abort()
			return
		}

		if err != nil {
			log.Print("Error occurred!")
			utils.SendMessageWithStatus(context, err.Error(), 400)
			context.Abort()
			return
		}

		context.Next()
	}
}

func CommentDeletionMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")
		comment := GetByCommentIDAndUserID(context)

		if authHeader == "" {
			utils.SendMessageWithStatus(context, "Authorize!", 404)
			context.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.SendMessageWithStatus(context, "Incorrect authorization header", 403)
			context.Abort()
			return
		}

		token := headerParts[1]

		user, err := parseToken(token)

		if user.ID != comment.UserID {
			utils.SendMessageWithStatus(context, "U are not allowed to delete someone's comment!", 403)
			context.Abort()
			return
		}

		if err != nil {
			log.Print("Error occurred!")
			utils.SendMessageWithStatus(context, err.Error(), 400)
			context.Abort()
			return
		}

		context.Next()
	}
}

func GetByCommentIDAndUserID(context *gin.Context) models.Comment {
	ID, err := strconv.Atoi(context.Param("id"))
	UserID, err := strconv.Atoi(context.Param("user_id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		context.Abort()
	}

	comment := models.Comment{}
	query := models.Comment{
		UserID:    UserID,
		CommentID: ID,
	}

	err = db.DB.First(&comment, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "That user does not have such comment", 404)
		context.Abort()
	}

	return comment
}

func GetByPostIDAndUserID(context *gin.Context) models.Post {
	ID, err := strconv.Atoi(context.Param("id"))
	UserID, err := strconv.Atoi(context.Param("user_id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		context.Abort()
	}

	post := models.Post{}
	query := models.Post{
		UserID: UserID,
		PostID: ID,
	}

	err = db.DB.First(&post, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "That user does not have such post", 404)
		context.Abort()
	}

	return post
}
