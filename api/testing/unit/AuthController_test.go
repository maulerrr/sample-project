package unit

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/maulerrr/sample-project/api/ctrl"
	"github.com/maulerrr/sample-project/api/db"
	"github.com/maulerrr/sample-project/api/dto"
	"github.com/maulerrr/sample-project/api/models"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestLogin(t *testing.T) {
	db.ConnectDB()

	type response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	type testcase struct {
		name     string
		payload  interface{}
		expected response
	}

	testcases := []testcase{
		{
			name:    "Test: Success",
			payload: dto.Login{Email: "test@mail.ru", Password: "test"},
			expected: response{
				Code:    200,
				Message: "success",
				Data: models.TokenResponse{
					UserID:   2,
					Username: "test",
					Email:    "test@mail.ru",
				},
			},
		},
		{
			name:    "Test: User does not exist",
			payload: dto.Login{Email: "random@random.ru", Password: "random"},
			expected: response{
				Code:    404,
				Message: "User not found",
			},
		},
		{
			name: "Test: Invalid JSON",
			expected: response{
				Code:    400,
				Message: "Invalid JSON",
			},
		},
		{
			name:    "Test: Incorrect Password",
			payload: dto.Login{Email: "test@mail.ru", Password: "123456"},
			expected: response{
				Code:    401,
				Message: "Password is not correct",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			loginJSON, _ := json.Marshal(tc.payload)
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			request := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(loginJSON))
			if tc.payload == nil {
				request = httptest.NewRequest(http.MethodPost, "/auth/signup", nil)
			}

			context.Request = request

			ctrl.Login(context)

			if recorder.Code != tc.expected.Code {
				t.Errorf("Expected status code %d but got %d", tc.expected.Code, recorder.Code)
			}

			var response response
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
			}

			if reflect.ValueOf(response.Data).Kind() == reflect.Map {
				dataMap := response.Data
				dataStruct := models.TokenResponse{}

				jsonData, _ := json.Marshal(dataMap)
				json.Unmarshal(jsonData, &dataStruct)

				dataStruct.Token = ""
				response.Data = dataStruct
			}

			if !reflect.DeepEqual(response, tc.expected) {
				t.Errorf("Expected response %v but got %v", tc.expected, response)
			}
		})
	}
}

func TestSignUp(t *testing.T) {
	db.ConnectDB()

	type response struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	type testcase struct {
		name     string
		payload  interface{}
		expected response
	}

	testcases := []testcase{
		{
			name:    "Test: Success",
			payload: dto.Registration{Username: "test22", Email: "test22@mail.ru", Password: "test"},
			expected: response{
				Code:    200,
				Message: "success",
				Data: models.TokenResponse{
					UserID:   22,
					Username: "test22",
					Email:    "test22@mail.ru",
				},
			},
		},
		{
			name:    "Test: User already exist",
			payload: dto.Registration{Username: "test", Email: "test@mail.ru", Password: "test"},
			expected: response{
				Code:    400,
				Message: "User already exists",
			},
		},
		{
			name: "Test: Invalid JSON",
			expected: response{
				Code:    400,
				Message: "Invalid JSON",
			},
		},
		{
			name:    "Test: Invalid Email Address",
			payload: dto.Registration{Username: "test2", Email: "test2", Password: "test"},
			expected: response{
				Code:    400,
				Message: "Invalid Email Address",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			loginJSON, _ := json.Marshal(tc.payload)
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			request := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(loginJSON))
			if tc.payload == nil {
				request = httptest.NewRequest(http.MethodPost, "/auth/signup", nil)
			}

			context.Request = request

			ctrl.SignUp(context)

			if recorder.Code != tc.expected.Code {
				t.Errorf("Expected status code %d but got %d", tc.expected.Code, recorder.Code)
			}

			var response response
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
			}

			if reflect.ValueOf(response.Data).Kind() == reflect.Map {
				dataMap := response.Data
				dataStruct := models.TokenResponse{}

				jsonData, _ := json.Marshal(dataMap)
				json.Unmarshal(jsonData, &dataStruct)

				dataStruct.Token = ""
				response.Data = dataStruct
			}

			if !reflect.DeepEqual(response, tc.expected) {
				t.Errorf("Expected response %v but got %v", tc.expected, response)
			}
		})
	}
}
