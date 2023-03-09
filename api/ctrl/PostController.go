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

	if err := context.BindJSON(json); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	newPost := models.Post{
		Header: json.Header,
		Body:   json.Body,
	}

	db.DB.Create(&newPost)

	context.JSON(200, newPost)
}

func DeletePostByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	post, query := models.Post{}, models.Post{PostID: id}

	err = db.DB.First(&post, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Couldn't delete! Post is not found", 404)
		return
	}

	db.DB.Model(&post).Association("posts").Delete()
	db.DB.Delete(&post)

	context.JSON(200, "Successfully deleted the post!")
}

func GetPostByID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	post := models.Post{}
	query := models.Post{PostID: id}

	err = db.DB.First(&post, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Post is not found", 404)
		return
	}

	context.JSON(200, post)
}

//func UpdatePostByID(context *gin.Context) {
//	id, err := strconv.Atoi(context.Param("id"))
//
//	if err != nil {
//		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
//		return
//	}
//
//	updatedPost := new(dto.UpdatePost)
//
//	if updatedPost.Body == "" {
//
//	}
//}
