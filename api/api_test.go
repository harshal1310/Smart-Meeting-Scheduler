package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessJson(t *testing.T) {
	tests := []struct {
		name        string
		data        interface{}
		expectError bool
	}{
		{
			name:        "Valid data",
			data:        map[string]string{"message": "success"},
			expectError: false,
		},
		{
			name:        "Array data",
			data:        []string{"item1", "item2"},
			expectError: false,
		},
		{
			name:        "Nil data",
			data:        nil,
			expectError: false,
		},
		{
			name: "Complex struct",
			data: struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			}{ID: 1, Name: "Test"},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Call SuccessJson
			SuccessJson(w, req, tt.data)

			// Check status code
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check CORS header
			corsHeader := w.Header().Get("Access-Control-Allow-Origin")
			if corsHeader != "*" {
				t.Errorf("Expected CORS header *, got %s", corsHeader)
			}

			// Verify response is valid JSON
			var response interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Errorf("Response is not valid JSON: %v", err)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		code           int
		expectedStatus int
	}{
		{
			name:           "Custom error with code",
			err:            errors.New("test error"),
			code:           http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Error with zero code (uses default)",
			err:            errors.New("server error"),
			code:           0,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Nil error",
			err:            nil,
			code:           http.StatusBadRequest,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Call Error
			Error(w, req, tt.err, tt.code)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check content type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("Expected Content-Type application/json, got %s", contentType)
			}

			// Check CORS header
			corsHeader := w.Header().Get("Access-Control-Allow-Origin")
			if corsHeader != "*" {
				t.Errorf("Expected CORS header *, got %s", corsHeader)
			}

			// Verify response is valid JSON with error structure
			var errorResponse ErrorResponse
			if err := json.Unmarshal(w.Body.Bytes(), &errorResponse); err != nil {
				t.Errorf("Error response is not valid JSON: %v", err)
			}

			// Check error message
			if tt.err != nil && errorResponse.Message != tt.err.Error() {
				t.Errorf("Expected error message '%s', got '%s'", tt.err.Error(), errorResponse.Message)
			}
			if tt.err == nil && errorResponse.Message != "nil err" {
				t.Errorf("Expected 'nil err' message for nil error, got '%s'", errorResponse.Message)
			}
		})
	}
}

func TestSuccess(t *testing.T) {
	tests := []struct {
		name    string
		jsonMsg []byte
	}{
		{
			name:    "Valid JSON message",
			jsonMsg: []byte(`{"message": "success"}`),
		},
		{
			name:    "Array JSON message",
			jsonMsg: []byte(`[1, 2, 3]`),
		},
		{
			name:    "Empty JSON object",
			jsonMsg: []byte(`{}`),
		},
		{
			name:    "Empty array",
			jsonMsg: []byte(`[]`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request and response recorder
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Call Success
			Success(w, req, tt.jsonMsg)

			// Check status code (Success doesn't set status, defaults to 200)
			if w.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d", w.Code)
			}

			// Check CORS header
			corsHeader := w.Header().Get("Access-Control-Allow-Origin")
			if corsHeader != "*" {
				t.Errorf("Expected CORS header *, got %s", corsHeader)
			}

			// Check response body
			if string(w.Body.Bytes()) != string(tt.jsonMsg) {
				t.Errorf("Expected body '%s', got '%s'", string(tt.jsonMsg), w.Body.String())
			}
		})
	}
}

func TestToHTTPStatusCode(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
	}{
		{
			name:           "Generic error",
			err:            errors.New("generic error"),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "Nil error",
			err:            nil,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := toHTTPStatusCode(tt.err)
			if status != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, status)
			}
		})
	}
}

func TestErrorResponseStructure(t *testing.T) {
	// Test ErrorResponse struct
	err := ErrorResponse{Message: "test error"}

	jsonData, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		t.Errorf("Failed to marshal ErrorResponse: %v", marshalErr)
	}

	var unmarshaled ErrorResponse
	if unmarshalErr := json.Unmarshal(jsonData, &unmarshaled); unmarshalErr != nil {
		t.Errorf("Failed to unmarshal ErrorResponse: %v", unmarshalErr)
	}

	if unmarshaled.Message != "test error" {
		t.Errorf("Expected message 'test error', got '%s'", unmarshaled.Message)
	}
}

func TestAPIResponseHeaders(t *testing.T) {
	// Test that all API functions set proper headers
	req := httptest.NewRequest("GET", "/test", nil)

	t.Run("SuccessJson headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		SuccessJson(w, req, map[string]string{"test": "data"})

		if w.Header().Get("Content-Type") != "application/json" {
			t.Error("SuccessJson should set Content-Type to application/json")
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("SuccessJson should set CORS header")
		}
	})

	t.Run("Error headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		Error(w, req, errors.New("test"), http.StatusBadRequest)

		if w.Header().Get("Content-Type") != "application/json" {
			t.Error("Error should set Content-Type to application/json")
		}
		if w.Header().Get("Access-Control-Allow-Origin") != "*" {
			t.Error("Error should set CORS header")
		}
	})
}

func TestAPILogging(t *testing.T) {
	// Test that API functions log appropriately (this is hard to test without capturing logs)
	req := httptest.NewRequest("GET", "/test", nil)

	t.Run("Success logging", func(t *testing.T) {
		w := httptest.NewRecorder()
		// This should log but we can't easily capture it in tests
		Success(w, req, []byte(`{"test": "data"}`))
		// Just verify it doesn't panic
	})

	t.Run("Error logging", func(t *testing.T) {
		w := httptest.NewRecorder()
		// This should log but we can't easily capture it in tests
		Error(w, req, errors.New("test error"), http.StatusBadRequest)
		// Just verify it doesn't panic
	})
}
