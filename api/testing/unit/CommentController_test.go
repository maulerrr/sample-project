package unit

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	testing2 "github.com/maulerrr/sample-project/api/testing"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetAllComments(t *testing.T) {
	db.ConnectDB()

	type testcase struct {
		name     string
		param    gin.Param
		expected int
	}

	testcases := []testcase{
		{
			name:     "Test: Success",
			param:    gin.Param{Key: "post_id", Value: strconv.Itoa(1)},
			expected: 200,
		},
		{
			name:     "Test: Not Found",
			param:    gin.Param{Key: "post_id", Value: strconv.Itoa(-1)},
			expected: 404,
		},
		{
			name:     "Test: Invalid ID (empty)",
			param:    gin.Param{Key: "post_id", Value: ""},
			expected: 400,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Params = append(context.Params, tc.param)

			ctrl.GetAllComments(context)

			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}

func TestCreateComment(t *testing.T) {
	db.ConnectDB()

	type testcase struct {
		name     string
		payload  dto.CreateCommentDTO
		expected int
	}

	testcases := []testcase{
		{
			name: "Test: Successful comment creation",
			payload: dto.CreateCommentDTO{
				UserID: 1,
				PostID: 1,
				Text:   "Test Comment",
			},
			expected: 200,
		},
		{
			name: "Test: User not found",
			payload: dto.CreateCommentDTO{
				UserID: 0,
				PostID: 1,
				Text:   "Test Comment",
			},
			expected: 404,
		},
		{
			name: "Test: Post not found",
			payload: dto.CreateCommentDTO{
				UserID: 1,
				PostID: 0,
				Text:   "Test Comment",
			},
			expected: 404,
		},
		{
			name:     "Test: Invalid JSON (empty)",
			payload:  dto.CreateCommentDTO{},
			expected: 400,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			commentJSON, _ := json.Marshal(tc.payload)

			request := httptest.NewRequest(http.MethodPost, "/comment", bytes.NewBuffer(commentJSON))
			if tc.payload.UserID == 0 && tc.payload.PostID == 0 && tc.payload.Text == "" {
				request = httptest.NewRequest(http.MethodPost, "/comment", nil)
			}
			request.Header.Set("Content-Type", "application/json")
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Request = request

			ctrl.CreateComment(context)

			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}

func TestDeleteComment(t *testing.T) {
	db.ConnectDB()

	newComment := models.Comment{
		CommentID: 1,
		UserID:    1,
		PostID:    1,
		Text:      "Test comment text",
	}
	db.DB.Create(&newComment)

	type testcase struct {
		name     string
		params   gin.Params
		expected int
	}

	testcases := []testcase{
		{
			name: "Test: Success",
			params: gin.Params{
				gin.Param{Key: "user_id", Value: strconv.Itoa(1)},
				gin.Param{Key: "id", Value: strconv.Itoa(1)},
			},
			expected: 200,
		},
		{
			name: "Test: Invalid ID (empty)",
			params: gin.Params{
				gin.Param{Key: "user_id", Value: strconv.Itoa(1)},
				gin.Param{Key: "id", Value: ""},
			},
			expected: 400,
		},
		{
			name: "Test: Not Found",
			params: gin.Params{
				gin.Param{Key: "user_id", Value: strconv.Itoa(1)},
				gin.Param{Key: "id", Value: strconv.Itoa(-1)},
			},
			expected: 404,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Params = append(context.Params, tc.params[0], tc.params[1])

			ctrl.DeleteComment(context)

			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}

func TestFindAllCommentsByUserID(t *testing.T) {
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
			name:     "Test: Not Found",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(10)},
			expected: 404,
		},
		{
			name:     "Test: Invalid ID (empty)",
			param:    gin.Param{Key: "id", Value: ""},
			expected: 400,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			context, _ := gin.CreateTestContext(httptest.NewRecorder())
			context.Params = append(context.Params, tc.param)

			ctrl.FindAllCommentsByUserID(context)

			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}

func TestFindCommentByWords(t *testing.T) {
	db.ConnectDB()

	type testcase struct {
		name     string
		payload  dto.FindByWordsDTO
		expected int
	}

	testcases := []testcase{
		{
			name:     "Test: Success",
			payload:  dto.FindByWordsDTO{Text: "Test"},
			expected: 200,
		},
		{
			name:     "Test: Not Found",
			payload:  dto.FindByWordsDTO{Text: "there should be no result"},
			expected: 404,
		},
		{
			name:     "Test: Invalid JSON (empty)",
			expected: 400,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			requestJSON, _ := json.Marshal(tc.payload)
			request := httptest.NewRequest(http.MethodGet, "/byword", bytes.NewBuffer(requestJSON))
			if tc.expected == 400 {
				request = httptest.NewRequest(http.MethodGet, "/byword", nil)
			}
			context.Request = request

			ctrl.FindCommentByWords(context)

			log.Println(recorder.Body)
			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}

func TestFindComment1(t *testing.T) {
	db.ConnectDB()

	testcases := []testing2.Testcase{
		{
			Name:         "Test: Success",
			Payload:      dto.FindByWordsDTO{Text: "Test"},
			ExpectedCode: 200,
		},
		{
			Name:         "Test: Not Found",
			Payload:      dto.FindByWordsDTO{Text: "there should be no result"},
			ExpectedCode: 404,
		},
		{
			Name:         "Test: Invalid JSON (empty)",
			ExpectedCode: 400,
		},
	}

	testing2.TestRun(testcases, ctrl.FindCommentByWords, http.MethodGet, false, t)
}
