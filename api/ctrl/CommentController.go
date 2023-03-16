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

type ResponseToCreateComment struct {
	Comment models.Comment `json:"comment"`
	User    models.User    `json:"user"`
}

type ResponseToGetAllComments struct {
	CommentID int    `json:"comment_id"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
}

func CreateComment(context *gin.Context) {
	json := new(dto.CreateCommentDTO)

	if err := context.BindJSON(json); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	newComment := models.Comment{
		UserID: json.UserID,
		PostID: json.PostID,
		Text:   json.Text,
	}

	user := models.User{}
	post := models.Post{}

	err := db.DB.First(&user, &json.UserID).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Cannot comment, user was not found!", 404)
		return
	}

	err = db.DB.First(&post, &json.PostID).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Cannot comment, post was not found!", 404)
		return
	}

	db.DB.Create(&newComment)

	response := ResponseToCreateComment{
		Comment: newComment,
		User:    user,
	}

	context.JSON(200, response)
}

func DeleteComment(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	comment, query := models.Comment{}, models.Comment{CommentID: id}

	err = db.DB.First(&comment, &query).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Couldn't delete! Comment is not found", 404)
		return
	}

	db.DB.Model(&comment).Association("comments").Delete()
	db.DB.Delete(&comment)

	context.JSON(200, "Successfully deleted the comment!")
}

func FindCommentByWords(context *gin.Context) {
	json := new(dto.FindByWordsDTO)

	if err := context.BindJSON(json); err != nil {
		utils.SendMessageWithStatus(context, "Invalid JSON", 400)
		return
	}

	found := models.Comment{}

	err := db.DB.Where("text LIKE ?", "%"+json.Text+"%").First(&found).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "Comment with such keywords was not found", 404)
		return
	}

	context.JSON(200, found)
}

func FindAllCommentsByUserID(context *gin.Context) {
	id, err := strconv.Atoi(context.Param("id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	comments := []models.Comment{}
	query := models.Comment{UserID: id}

	db.DB.Find(&comments, &query)

	context.JSON(200, comments)
}

func GetAllComments(context *gin.Context) {
	PostID, err := strconv.Atoi(context.Param("post_id"))

	if err != nil {
		utils.SendMessageWithStatus(context, "Invalid ID format", 400)
		return
	}

	response := []ResponseToGetAllComments{}

	comment := models.Comment{}
	post := models.Post{}
	SearchQuery := models.Post{PostID: PostID}
	DeleteQuery := models.Comment{PostID: PostID}

	err = db.DB.First(&post, &SearchQuery).Error

	if err == gorm.ErrRecordNotFound {
		utils.SendMessageWithStatus(context, "That post does not exist! Truncating all comments..", 404)
		db.DB.Model(&comment).Delete(&comment, &DeleteQuery)
		return
	}

	db.DB.Raw(`SELECT comments.comment_id, comments.user_id, users.username, comments.text
FROM comments INNER JOIN users ON comments.user_id = users.user_id
WHERE comments.post_id=?`, PostID).Scan(&response)

	context.JSON(200, response)
}
