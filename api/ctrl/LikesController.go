package ctrl

import (
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/utils"
	"gorm.io/gorm"
	"log"
)

type Response struct {
	Data    models.Like `json:"data"`
	Message string      `json:"message"`
	Liked   bool        `json:"liked"`
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

	log.Print(response)

	context.JSON(200, response)
}
