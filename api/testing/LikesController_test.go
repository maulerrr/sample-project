package testing

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestAddLike(t *testing.T) {
	db.ConnectDB()

	post := models.Post{PostID: 5, Header: "Test Post", Body: "Test Post Body", UserID: 1}
	db.DB.Create(&post)

	tests := []struct {
		name         string
		payload      dto.AddLike
		statusCode   int
		expectedResp ctrl.Response
	}{
		{
			name: "Test: Liked",
			payload: dto.AddLike{
				UserID: 2,
				PostID: post.PostID,
			},
			statusCode: 200,
			expectedResp: ctrl.Response{
				Message: "Removed Like",
				Liked:   false,
			},
		},
		{
			name: "Test: Unliked",
			payload: dto.AddLike{
				UserID: 2,
				PostID: post.PostID,
			},
			statusCode: 200,
			expectedResp: ctrl.Response{
				Message: "Removed Like",
				Liked:   false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			likeJSON, _ := json.Marshal(tt.payload)
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			request := httptest.NewRequest(http.MethodPost, "/post/like", bytes.NewBuffer(likeJSON))
			context.Request = request

			ctrl.AddLike(context)

			if recorder.Code != tt.statusCode {
				t.Errorf("Expected status code %d but got %d", tt.statusCode, recorder.Code)
			}

			var response ctrl.Response
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
			}
			if !reflect.DeepEqual(response.Liked, tt.expectedResp.Liked) {
				t.Errorf("Expected response %v but got %v", tt.expectedResp, response)
			}
		})
	}

}

func TestGetLike(t *testing.T) {
	db.ConnectDB()

	post := models.Post{PostID: 5, Header: "Test Post", Body: "Test Post Body", UserID: 1}
	like := models.Like{PostID: post.PostID, UserID: 1}
	db.DB.Create(&post)
	db.DB.Create(&like)

	tests := []struct {
		name           string
		userID, postID string
		statusCode     int
		expectedResp   ctrl.GetLikeResponse
	}{
		{
			name:       "Test: Liked",
			userID:     strconv.Itoa(post.UserID),
			postID:     strconv.Itoa(post.PostID),
			statusCode: 200,
			expectedResp: ctrl.GetLikeResponse{
				Status: 200,
				Liked:  true,
			},
		},
		{
			name:       "Test: Not liked",
			userID:     strconv.Itoa(post.UserID),
			postID:     strconv.Itoa(6),
			statusCode: 404,
			expectedResp: ctrl.GetLikeResponse{
				Status: 404,
				Liked:  false,
			},
		},
	}

	// run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			context.Params = append(context.Params,
				gin.Param{Key: "user_id", Value: tt.userID},
				gin.Param{Key: "id", Value: tt.postID})

			ctrl.GetLike(context)

			// check response status code
			if recorder.Code != tt.statusCode {
				t.Errorf("Expected status code %d but got %d", tt.statusCode, recorder.Code)
			}

			// check response body
			var response ctrl.GetLikeResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
			}
			if !reflect.DeepEqual(response, tt.expectedResp) {
				t.Errorf("Expected response %v but got %v", tt.expectedResp, response)
			}
		})
	}
}

func TestGetLikesCountOnPost(t *testing.T) {
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
			name:     "Test: Post not found",
			param:    gin.Param{Key: "id", Value: strconv.Itoa(-1)},
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

			ctrl.GetLikesCountOnPost(context)

			assert.Equal(t, tc.expected, context.Writer.Status())
		})
	}
}
