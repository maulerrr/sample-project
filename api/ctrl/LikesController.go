package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"gorm.io/gorm"
	"strconv"
)

type Response struct {
	Data    models.Like `json:"data"`
	Message string      `json:"message"`
	Liked   bool        `json:"liked"`
}

type GetLikeResponse struct {
	Status int  `json:"status"`
	Liked  bool `json:"liked"`
}

func AddLike(context *gin.Context) {
	json := new(dto.AddLike)

	if err := context.BindJSON(json); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	newLike := models.Like{
		UserID: json.UserID,
		PostID: json.PostID,
	}

	found := models.Like{}

	err := db.DB.First(&found, &newLike).Error

	if err != gorm.ErrRecordNotFound {
		db.DB.Delete(&found, &newLike)

		response := Response{
			Data:    newLike,
			Message: "Removed Like",
			Liked:   false,
		}
		context.JSON(200, response)
		return
	}

	db.DB.Create(&newLike)

	response := Response{
		Data:    newLike,
		Message: "Liked",
		Liked:   true,
	}

	context.JSON(200, response)
}

func GetLike(context *gin.Context) {
	UserID, err := strconv.Atoi(context.Param("user_id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	PostID, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	found := models.Like{}
	query := models.Like{
		PostID: PostID,
		UserID: UserID,
	}
	dbErr := db.DB.First(&found, &query).Error

	if dbErr == gorm.ErrRecordNotFound {
		context.JSON(404, GetLikeResponse{Status: 404, Liked: false})
		return
	}

	context.JSON(200, GetLikeResponse{Status: 200, Liked: true})
}

func GetLikesCountOnPost(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	post := models.Post{}
	err = db.DB.First(&post, id).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Post not found", 404)
		return
	}

	var count int
	db.DB.Raw(`SELECT COUNT(post_id) FROM likes WHERE post_id=?`, id).Scan(&count)

	context.JSON(200, count)
}
