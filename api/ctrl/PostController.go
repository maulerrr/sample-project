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
		UserID: json.UserID,
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

	comments := models.Comment{}
	likes := models.Like{}
	DeleteCommentsQuery := models.Comment{PostID: id}
	DeleteLikesQuery := models.Like{PostID: id}

	err = db.DB.First(&post, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Couldn't delete! Post is not found", 404)
		db.DB.Model(&comments).Delete(&comments, &DeleteCommentsQuery)
		db.DB.Model(&likes).Delete(&likes, &DeleteLikesQuery)
		return
	}

	db.DB.Model(&post).Association("posts").Delete()
	db.DB.Model(&comments).Delete(&comments, &DeleteCommentsQuery)
	db.DB.Model(&likes).Delete(&likes, &DeleteLikesQuery)
	db.DB.Delete(&post)

	context.JSON(200, "Successfully deleted the post!")
}

func GetByPostID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	post := models.Post{}
	SearchPostQuery := models.Post{PostID: id}

	//to refresh DB, truncate all unnecessary things after unchecked deletion
	comments := models.Comment{}
	likes := models.Like{}
	DeleteCommentsQuery := models.Comment{PostID: id}
	DeleteLikesQuery := models.Like{PostID: id}

	err = db.DB.First(&post, &SearchPostQuery).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Post is not found", 404)
		db.DB.Model(&comments).Delete(&comments, &DeleteCommentsQuery)
		db.DB.Model(&likes).Delete(&likes, &DeleteLikesQuery)
		return
	}

	context.JSON(200, post)
}

func GetByUserID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	posts := []models.Post{}
	query := models.Post{UserID: id}

	err = db.DB.Find(&posts, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Post is not found", 404)
		return
	}

	context.JSON(200, posts)
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
