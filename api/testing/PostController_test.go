package testing

import (
	"bytes"
	"encoding/json"
	"github.com/maulerrr/sample-project/api/db"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
)

func TestGetAllPosts(t *testing.T) {
	db.ConnectDB()

	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	ctrl.GetAllPosts(context)

	assert.Equal(t, 200, context.Writer.Status())
}

func TestAddPost1(t *testing.T) {
	db.ConnectDB()

	testcases := []testcase{
		{
			name: "Test: Success",
			payload: dto.CreatePost{
				UserID: 2,
				Header: "New Test",
				Body:   "Body",
			},
			expectedCode: 200,
			expectedData: nil,
		},
		{
			name:         "Test: Invalid JSON",
			expectedCode: 400,
			expectedData: gin.H{"code": 400, "message": "Invalid JSON"},
		},
	}

	TestRun(testcases, ctrl.AddPost, http.MethodPost, true, t)
}

func TestAddPost(t *testing.T) {
	db.ConnectDB()

	context, _ := gin.CreateTestContext(httptest.NewRecorder())

	newPost := dto.CreatePost{
		UserID: 1,
		Header: "Test Header",
		Body:   "Test Body",
	}

	postJSON, _ := json.Marshal(newPost)

	// Create a new HTTP request to add the post
	request, _ := http.NewRequest(http.MethodPost, "/post", bytes.NewBuffer(postJSON))
	request.Header.Set("Content-Type", "application/json")

	context.Request = request
	ctrl.AddPost(context)

	assert.Equal(t, 200, context.Writer.Status())

	var post models.Post
	db.DB.Last(&post)
	assert.Greater(t, post.PostID, 0)
}

func TestDeletePostByID(t *testing.T) {
	db.ConnectDB()

	context, _ := gin.CreateTestContext(httptest.NewRecorder())

	newPost := models.Post{
		UserID: 1,
		Header: "Test Header",
		Body:   "Test Body",
	}
	db.DB.Create(&newPost)

	context.Params = append(context.Params, gin.Param{Key: "user_id", Value: strconv.Itoa(newPost.UserID)})
	context.Params = append(context.Params, gin.Param{Key: "id", Value: strconv.Itoa(newPost.PostID)})
	ctrl.DeletePostByID(context)

	assert.Equal(t, 200, context.Writer.Status())

	var post models.Post
	db.DB.First(&post, newPost.PostID)
	assert.Equal(t, 0, post.PostID)
}

func TestGetByPostID(t *testing.T) {
	db.ConnectDB()

	type testcase struct {
		name     string
		param    gin.Param
		expected int
	}

	testcases := []testcase{
		{
			name:     "Test: Success",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(1)},
			expected: 200,
		},
		{
			name:     "Test: Invalid ID",
			param:    gin.Param{Key: "id", Value: ""},
			expected: 400,
		},
		{
			name:     "Test: Not found",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(-1)},
			expected: 404,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Params = append(context.Params, tc.param)

			ctrl.GetByPostID(context)

			assert.Equal(t, tc.expected, recorder.Code)
		})
	}
}

func TestGetAllByUserID(t *testing.T) {
	db.ConnectDB()

	type testcase struct {
		name     string
		param    gin.Param
		expected int
	}

	testcases := []testcase{
		{
			name:     "Test: Success",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(1)},
			expected: 200,
		},
		{
			name:     "Test: Invalid ID",
			param:    gin.Param{Key: "id", Value: ""},
			expected: 400,
		},
		{
			name:     "Test: No posts",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(-1)},
			expected: 404,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)
			context.Params = append(context.Params, tc.param)

			ctrl.GetAllByUserID(context)

			assert.Equal(t, tc.expected, recorder.Code)
		})
	}
}
