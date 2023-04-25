package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"github.com/maulerrr/sample-project/api/routes"
	"github.com/maulerrr/sample-project/api/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

func TestGetAllPosts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	token, err := getToken(router, user.Email, "password")

	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Authorized user gets all posts",
			headers:      map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			expectedCode: 200,
			expectedBody: `[]`,
		},
		{
			name:         "Unauthorized user can't get all posts",
			headers:      map[string]string{},
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/api/v1/post/", nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())
		})
	}

	db.DB.Delete(&user)
}

func TestAddPost(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	token, err := getToken(router, user.Email, "password")
	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	newPost := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		payload      interface{}
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Authorized user adds post",
			headers:      map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			payload:      newPost,
			expectedCode: 200,
			expectedBody: newPost,
		},
		{
			name:         "Unauthorized user can't add post",
			headers:      map[string]string{},
			payload:      newPost,
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/post/", bytes.NewBuffer(reqBody))
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response models.Post

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal("error unmarshalling response..")
				}

				response.PostID = newPost.PostID
				assert.Equal(t, tc.expectedBody, response)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}

	db.DB.Delete(&user)
}

func TestGetPostByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	post := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}
	db.DB.Create(&post)

	token, err := getToken(router, user.Email, "password")

	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		params       gin.Params
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:    "Authorized user gets post",
			headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			params: gin.Params{gin.Param{
				Key:   "id",
				Value: strconv.Itoa(post.PostID),
			}},
			expectedCode: 200,
			expectedBody: post,
		},
		{
			name:    "Unauthorized user can't get post",
			headers: map[string]string{},
			params: gin.Params{gin.Param{
				Key:   "id",
				Value: strconv.Itoa(post.PostID),
			}},
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/post/%s", tc.params.ByName("id"))
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response models.Post

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal("error unmarshalling response..")
				}

				assert.Equal(t, tc.expectedBody, response)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}
	db.DB.Delete(&post)
	db.DB.Delete(&user)
}

func TestDeletePostByID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	user2 := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser2@example.com",
		Username:  "testuser2",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user2)

	post := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}
	db.DB.Create(&post)

	token, err := getToken(router, user.Email, "password")
	token2, err := getToken(router, user2.Email, "password")

	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		params       gin.Params
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:    "Authorized user deletes post",
			headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			params: gin.Params{
				gin.Param{Key: "id", Value: strconv.Itoa(post.PostID)},
				gin.Param{Key: "user_id", Value: strconv.Itoa(user.UserID)},
			},
			expectedCode: 200,
			expectedBody: `"Successfully deleted the post!"`,
		},
		{
			name:    "Unauthorized user can't delete post",
			headers: map[string]string{},
			params: gin.Params{
				gin.Param{Key: "id", Value: strconv.Itoa(post.PostID)},
				gin.Param{Key: "user_id", Value: strconv.Itoa(user.UserID)},
			},
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
		{
			name:    "Authorized user is not allowed to delete someone else's posts",
			headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token2)},
			params: gin.Params{
				gin.Param{Key: "id", Value: strconv.Itoa(post.PostID)},
				gin.Param{Key: "user_id", Value: strconv.Itoa(user.UserID)},
			},
			expectedCode: 403,
			expectedBody: `{"code":403,"message":"U are not allowed to delete someone's post!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/post/%s/%s", tc.params.ByName("user_id"), tc.params.ByName("id"))
			req, _ := http.NewRequest(http.MethodDelete, url, nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)
			assert.Equal(t, tc.expectedBody, w.Body.String())

		})
	}

	db.DB.Delete(&post)
	db.DB.Delete(&user)
	db.DB.Delete(&user2)
}

func TestAddLike(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	token, err := getToken(router, user.Email, "password")
	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	post := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}
	db.DB.Create(&post)

	newLike := models.Like{
		LikeID: int(uuid.New().ID()),
		UserID: user.UserID,
		PostID: post.PostID,
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		payload      interface{}
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Authorized user likes post",
			headers:      map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			payload:      newLike,
			expectedCode: 200,
			expectedBody: LikeResponse{
				Data:    newLike,
				Message: "Liked",
				Liked:   true,
			},
		},
		{
			name:         "Unauthorized user can't like post",
			headers:      map[string]string{},
			payload:      newLike,
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/post/like", bytes.NewBuffer(reqBody))
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response LikeResponse

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal("error unmarshalling response..")
				}

				response.Data.LikeID = newLike.LikeID
				assert.Equal(t, tc.expectedBody, response)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}

	db.DB.Delete(&user)
	db.DB.Delete(&post)
	db.DB.Delete(&newLike)
}

func TestRemoveLike(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	token, err := getToken(router, user.Email, "password")
	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	post := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}
	db.DB.Create(&post)

	newLike := models.Like{
		LikeID: int(uuid.New().ID()),
		UserID: user.UserID,
		PostID: post.PostID,
	}
	db.DB.Create(&newLike)

	testCases := []struct {
		name         string
		headers      map[string]string
		payload      interface{}
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:         "Authorized user removes like on post",
			headers:      map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			payload:      newLike,
			expectedCode: 200,
			expectedBody: LikeResponse{
				Data:    newLike,
				Message: "Removed Like",
				Liked:   false,
			},
		},
		{
			name:         "Unauthorized user can't interact post",
			headers:      map[string]string{},
			payload:      newLike,
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			reqBody, _ := json.Marshal(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/api/v1/post/like", bytes.NewBuffer(reqBody))
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response LikeResponse

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal("error unmarshalling response..")
				}

				response.Data.LikeID = newLike.LikeID
				assert.Equal(t, tc.expectedBody, response)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}

	db.DB.Delete(&user)
	db.DB.Delete(&post)
}

func TestGetLike(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	routes.InitRoutes(router)

	mockDB := mockDB()
	db.DB = mockDB

	//os.Setenv("JWT_KEY", "testkey")

	user := models.User{
		UserID:    int(uuid.New().ID()),
		Email:     "testuser@example.com",
		Username:  "testuser",
		Password:  utils.HashPassword([]byte("password")),
		CreatedAt: time.Now(),
	}
	db.DB.Create(&user)

	post := models.Post{
		PostID: int(uuid.New().ID()),
		UserID: user.UserID,
		Header: "Test post HEADER",
		Body:   "Test post BODY",
	}
	db.DB.Create(&post)

	newLike := models.Like{
		LikeID: int(uuid.New().ID()),
		UserID: user.UserID,
		PostID: post.PostID,
	}
	db.DB.Create(&newLike)

	token, err := getToken(router, user.Email, "password")

	if err != nil {
		t.Errorf("Error getting token: %v", err)
	}

	testCases := []struct {
		name         string
		headers      map[string]string
		params       gin.Params
		expectedCode int
		expectedBody interface{}
	}{
		{
			name:    "Authorized user gets like on post",
			headers: map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
			params: gin.Params{
				gin.Param{Key: "id", Value: strconv.Itoa(post.PostID)},
				gin.Param{Key: "user_id", Value: strconv.Itoa(user.UserID)},
			},
			expectedCode: 200,
			expectedBody: GetLikeResponse{
				Status: 200,
				Liked:  true,
			},
		},
		{
			name:    "Unauthorized user can't get anything",
			headers: map[string]string{},
			params: gin.Params{
				gin.Param{Key: "id", Value: strconv.Itoa(post.PostID)},
				gin.Param{Key: "user_id", Value: strconv.Itoa(user.UserID)},
			},
			expectedCode: 401,
			expectedBody: `{"code":401,"message":"Authorize!"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			url := fmt.Sprintf("/api/v1/post/like/%s/%s", tc.params.ByName("user_id"), tc.params.ByName("id"))
			req, _ := http.NewRequest(http.MethodGet, url, nil)
			for key, value := range tc.headers {
				req.Header.Set(key, value)
			}
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedCode, w.Code)

			if w.Code == 200 {
				var response GetLikeResponse

				err := json.Unmarshal(w.Body.Bytes(), &response)
				if err != nil {
					log.Fatal("error unmarshalling response..")
				}

				assert.Equal(t, tc.expectedBody, response)
			} else {
				assert.Equal(t, tc.expectedBody, w.Body.String())
			}
		})
	}
	db.DB.Delete(&post)
	db.DB.Delete(&user)
	db.DB.Delete(&newLike)
}

// response structs
type getTokenResponse struct {
	Code    int                  `json:"code"`
	Message string               `json:"message"`
	Data    models.TokenResponse `json:"data"`
}

type LikeResponse struct {
	Data    models.Like `json:"data"`
	Message string      `json:"message"`
	Liked   bool        `json:"liked"`
}

type GetLikeResponse struct {
	Status int  `json:"status"`
	Liked  bool `json:"liked"`
}

// helper functions
func getToken(router *gin.Engine, email string, password string) (string, error) {
	loginCredentials := dto.Login{Email: email, Password: password}
	loginJSON, _ := json.Marshal(loginCredentials)
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginJSON))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != 200 {
		return "", fmt.Errorf("error getting token: %s", w.Body.String())
	}

	var tokenResponse getTokenResponse

	err := json.Unmarshal(w.Body.Bytes(), &tokenResponse)
	if err != nil {
		return "", fmt.Errorf("error getting token: %v", err)
	}

	return tokenResponse.Data.Token, nil
}

func mockDB() *gorm.DB {
	open, _ := gorm.Open(postgres.Open("postgres://postgres:1111@localhost:5432/mock_db?sslmode=disable"), &gorm.Config{})
	open.AutoMigrate(&models.User{})
	open.AutoMigrate(&models.Post{})
	open.AutoMigrate(&models.Like{})
	return open
}
