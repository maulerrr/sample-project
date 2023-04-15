package ctrl

import (
	"github.com/badoux/checkmail"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"gorm.io/gorm"
	"os"
	"time"
)

func Login(context *gin.Context) {
	credentials := new(dto.Login)

	if err := context.BindJSON(credentials); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	user := models.User{}
	query := models.User{Email: credentials.Email}
	err := db.DB.First(&user, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "User not found", 404)
		return
	}

	if !utils.ComparePasswords(user.Password, credentials.Password) {
		utils.SendMessageWithStatus(context, "Password is not correct", 401)
		return
	}

	tokenString, err := GenerateToken(user)

	if err != nil {
		utils.SendMessageWithStatus(context, "Auth error (token creation)", 500)
		return
	}

	response := &models.TokenResponse{
		UserID:   user.UserID,
		Username: user.Username,
		Email:    user.Email,
		Token:    tokenString,
	}

	utils.SendSuccessJSON(context, response)
}

func SignUp(context *gin.Context) {
	json := new(dto.Registration)

	if err := context.BindJSON(json); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	if len(json.Password) < 4 {
		utils.SendMessageWithStatus(context, "Minimum password length is 4", 400)
		return
	}

	password := utils.HashPassword([]byte(json.Password))
	err := checkmail.ValidateFormat(json.Email)
	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid Email Address", 400)
		return
	}

	newUser := models.User{
		Password:  password,
		Email:     json.Email,
		Username:  json.Username,
		CreatedAt: time.Now(),
	}

	found := models.User{}
	query := models.User{Email: json.Email}
	err = db.DB.First(&found, &query).Error
	if err != gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "User already exists", 400)
		return
	}

	err = db.DB.Create(&newUser).Error
	if err != nil {
		utils.SendMessageWithStatus(context, err.Error(), 400)
		return
	}

	tokenString, err := GenerateToken(newUser)

	if err != nil {
		utils.SendMessageWithStatus(context, "Auth error (token creation)", 500)
		return
	}

	response := &models.TokenResponse{
		UserID:   newUser.UserID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Token:    tokenString,
	}

	utils.SendSuccessJSON(context, response)
}

func GenerateToken(user models.User) (string, error) {
	expirationTime := time.Now().Add(time.Hour * 24)

	claims := &models.Claims{
		ID:       user.UserID,
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	jwtKey := os.Getenv("JWT_KEY")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))

	return tokenString, err
}
