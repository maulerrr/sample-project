package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"strings"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

func parseToken(tokenString string) (*models.Claims, error) {
	// Parse the JWT token and extract the claims
	// just to check and handle errors.
	//log.Println(jwtKey)

	claims := &models.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func AuthMiddleware() gin.HandlerFunc {

	return func(context *gin.Context) {
		authHeader := context.GetHeader("Authorization")

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

		claims := &models.Claims{}

		tkn, err := jwt.ParseWithClaims(token, claims,
			func(t *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				utils.SendMessageWithStatus(context, "Unauthorized", 401)
				context.Abort()
				return
			}

			utils.SendMessageWithStatus(context, err.Error(), 404)
			context.Abort()
			return
		}

		if !tkn.Valid {
			utils.SendMessageWithStatus(context, "Unauthorized", 401)
			context.Abort()
			return
		}

		context.Next()
	}
}

func AccessMiddleware() gin.HandlerFunc {
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
			utils.SendMessageWithStatus(context, "U are not allowed to delete someone's posts!", 403)
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
