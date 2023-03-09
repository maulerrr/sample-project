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

func GetAllPosts(context *gin.Context) {
	posts := []models.Post{}
	db.DB.Preload("posts").Find(&posts)

	context.JSON(200, posts)
}

func AddPost(context *gin.Context) {
	json := new(dto.CreatePost)

}

func DeletePostByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	err = db.DB.Delete(id).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Couldn't delete! Post is not found", 404)
	}

	context.JSON(200, nil)
}

func UpdatePostByID(context *gin.Context) {
	updatedPost := new(dto.UpdatePost)

	if updatedPost.Header == "" {

	}

	if updatedPost.Body == "" {

	}
}

func GetPostByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	post := models.Post{}
	query := models.Post{PostID: id}

	err = db.DB.First(post, query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Post not found", 404)
	}

	context.JSON(200, post)
}
