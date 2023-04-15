package testing

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testcase struct {
	name         string
	payload      interface{}
	params       gin.Params
	expectedCode int
	expectedData interface{}
}

func TestRun(testcases []testcase, function func(context *gin.Context), method string, withData bool, t *testing.T) {
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			context, _ := gin.CreateTestContext(recorder)

			requestJSON, _ := json.Marshal(tc.payload)
			request := httptest.NewRequest(method, "/does-not-matter", bytes.NewBuffer(requestJSON))
			if tc.payload == nil {
				request = httptest.NewRequest(method, "/does-not-matter", nil)
			}

			context.Params = tc.params
			context.Request = request

			function(context)

			if recorder.Code != tc.expectedCode {
				t.Errorf("Expected status code %d but got %d", tc.expectedCode, recorder.Code)
			}

			if !withData || tc.expectedData == nil {
				return
			}

			var response interface{}
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Errorf("Error unmarshalling response body: %v", err)
			}

			if reflect.ValueOf(response).Kind() == reflect.Map {
				dataMap := response
				var dataStruct interface{}

				jsonData, _ := json.Marshal(dataMap)
				json.Unmarshal(jsonData, &dataStruct)

				response = dataStruct
			}

			expectedJSON, _ := json.Marshal(tc.expectedData)
			actualJSON, _ := json.Marshal(response)

			assertJSONResponse(t, expectedJSON, actualJSON)

			//if !reflect.DeepEqual(response, tc.expectedData) {
			//	t.Errorf("Expected response %v but got %v", tc.expectedData, response)
			//}
		})
	}
}

func assertJSONResponse(t *testing.T, expectedJSON, actualJSON []byte) {
	var expected, actual interface{}
	if err := json.Unmarshal(expectedJSON, &expected); err != nil {
		t.Fatalf("Error unmarshalling expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actual); err != nil {
		t.Fatalf("Error unmarshalling actual JSON: %v", err)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected JSON response %s, but got %s", expectedJSON, actualJSON)
	}
}
