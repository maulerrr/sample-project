package middlewares

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"os"
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
			utils.SendMessageWithStatus(context, "Authorize!", 401)
			context.Abort()
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.SendMessageWithStatus(context, "Incorrect authorization header", 400)
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

			utils.SendMessageWithStatus(context, err.Error(), 401)
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
