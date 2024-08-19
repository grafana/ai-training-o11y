package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCallResourceResponseSender implements backend.CallResourceResponseSender
// for use in tests.
type mockCallResourceResponseSender struct {
	response *backend.CallResourceResponse
}

// Send sets the received *backend.CallResourceResponse to s.response
func (s *mockCallResourceResponseSender) Send(response *backend.CallResourceResponse) error {
	s.response = response
	return nil
}

// TestCallResource tests CallResource calls, using backend.CallResourceRequest and backend.CallResourceResponse.
// This ensures the httpadapter for CallResource works correctly.
func TestCallResource(t *testing.T) {
	// Set up a mock server for metadata requests
	mockMetadataServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"metadata": "test"}`))
	}))
	defer mockMetadataServer.Close()

	// Initialize app with the mock server URL
	testJSON := []byte(fmt.Sprintf(`{"metadataUrl": "%s", "metadataToken": "test-token"}`, mockMetadataServer.URL))
	inst, err := NewApp(context.Background(), backend.AppInstanceSettings{
		JSONData: testJSON,
	})
	require.NoError(t, err, "Failed to create new app")
	require.NotNil(t, inst, "App instance should not be nil")

	app, ok := inst.(*App)
	require.True(t, ok, "Instance should be of type *App")

	// Create a test server to mock the metadata service
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"), "Bearer token should be set")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"metadata": "test"}`))
	}))
	defer testServer.Close()

	// Update the app's metadataUrl to point to our test server
	app.metadataUrl = testServer.URL

	// Set up and run test cases
	testCases := []struct {
		name      string
		method    string
		path      string
		body      []byte
		expStatus int
		expBody   []byte
	}{
		{
			name:      "GET ping 200",
			method:    http.MethodGet,
			path:      "ping",
			expStatus: http.StatusOK,
			expBody:   []byte(`{"message":"ok"}`),
		},
		{
			name:      "GET echo 405",
			method:    http.MethodGet,
			path:      "echo",
			expStatus: http.StatusMethodNotAllowed,
		},
		{
			name:      "POST echo 200",
			method:    http.MethodPost,
			path:      "echo",
			body:      []byte(`{"message":"ok"}`),
			expStatus: http.StatusOK,
			expBody:   []byte(`{"message":"ok"}`),
		},
		{
			name:      "GET non-existing handler 404",
			method:    http.MethodGet,
			path:      "not_found",
			expStatus: http.StatusNotFound,
		},
		{
			name:      "POST echo with invalid JSON 400",
			method:    http.MethodPost,
			path:      "echo",
			body:      []byte(`{"message":invalid}`),
			expStatus: http.StatusBadRequest,
		},
		{
			name:      "GET metadata 200",
			method:    http.MethodGet,
			path:      "metadata/test",
			expStatus: http.StatusOK,
			expBody:   []byte(`{"metadata": "test"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var r mockCallResourceResponseSender
			err = app.CallResource(context.Background(), &backend.CallResourceRequest{
				Method: tc.method,
				Path:   tc.path,
				Body:   tc.body,
			}, &r)

			require.NoError(t, err, "CallResource should not return an error")
			require.NotNil(t, r.response, "No response received from CallResource")

			assert.Equal(t, tc.expStatus, r.response.Status, "Unexpected response status")

			if len(tc.expBody) > 0 {
				assert.JSONEq(t, string(tc.expBody), string(r.response.Body), "Unexpected response body")
			}
		})
	}
}

func TestMetadataHandlerTokenInjection(t *testing.T) {
	app := &App{
		metadataUrl:   "http://example.com",
		metadataToken: "test-token",
		stackId:       "test-stack-id",
	}

	// Create a test server to mock the metadata service
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer test-stack-id:test-token", r.Header.Get("Authorization"), "Bearer token should be set correctly")
		assert.Equal(t, "/test", r.URL.Path, "Path should be correctly modified")
		assert.Equal(t, "test-host", r.Header.Get("X-Forwarded-Host"), "X-Forwarded-Host should be set")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"metadata": "test"}`))
	}))
	defer testServer.Close()

	// Update the app's metadataUrl to point to our test server
	app.metadataUrl = testServer.URL

	// Create a test request
	req := httptest.NewRequest("GET", "/metadata/test", nil)
	req.Host = "test-host"  // Set the Host header

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the metadataHandler
	handler := app.metadataHandler(app.metadataUrl)
	handler(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")

	body, err := io.ReadAll(rr.Body)
	require.NoError(t, err, "Failed to read response body")

	var responseBody map[string]string
	err = json.Unmarshal(body, &responseBody)
	require.NoError(t, err, "Failed to unmarshal response body")

	assert.Equal(t, "test", responseBody["metadata"], "Unexpected response body")
}

func TestMetadataHandlerPathTransformations(t *testing.T) {
	app := &App{
		metadataUrl:   "http://example.com",
		metadataToken: "test-token",
		stackId:       "test-stack-id",
	}

	testCases := []struct {
		name           string
		inputPath      string
		expectedPath   string
		expectedStatus int
	}{
		{"Metadata prefix", "/metadata/api/v1/processes", "/api/v1/processes", http.StatusOK},
		{"No metadata prefix", "/api/v1/processes", "/api/v1/processes", http.StatusOK},
		{"Metadata in middle", "/some/path/metadata/api/v1/processes", "/api/v1/processes", http.StatusOK},
		{"Root path with metadata", "/metadata", "/", http.StatusOK},
		{"Root path", "/", "/", http.StatusOK},
		{"Metadata root path", "/metadata/", "/", http.StatusOK},
		{"Multiple metadata in path", "/metadata/something/metadata/api", "/something/metadata/api", http.StatusOK},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a test server to mock the metadata service
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, tc.expectedPath, r.URL.Path, "Path should be correctly modified")
				assert.Equal(t, "Bearer test-stack-id:test-token", r.Header.Get("Authorization"), "Bearer token should be set correctly")
				assert.Equal(t, "test-host", r.Header.Get("X-Forwarded-Host"), "X-Forwarded-Host should be set")
				w.WriteHeader(tc.expectedStatus)
				w.Write([]byte(`{"metadata": "test"}`))
			}))
			defer testServer.Close()

			// Update the app's metadataUrl to point to our test server
			app.metadataUrl = testServer.URL

			// Create a test request
			req := httptest.NewRequest("GET", tc.inputPath, nil)
			req.Host = "test-host"  // Set the Host header

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Call the metadataHandler
			handler := app.metadataHandler(app.metadataUrl)
			handler(rr, req)

			// Check the response
			assert.Equal(t, tc.expectedStatus, rr.Code, "Handler returned wrong status code")

			body, err := io.ReadAll(rr.Body)
			require.NoError(t, err, "Failed to read response body")

			var responseBody map[string]string
			err = json.Unmarshal(body, &responseBody)
			require.NoError(t, err, "Failed to unmarshal response body")

			assert.Equal(t, "test", responseBody["metadata"], "Unexpected response body")
		})
	}
}

func TestMetadataHandlerNoToken(t *testing.T) {
	app := &App{
		metadataUrl: "http://example.com",
		// No token set
	}

	// Create a test server to mock the metadata service
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get("Authorization"), "Authorization header should not be set")
		assert.Equal(t, "/test", r.URL.Path, "Path should be correctly modified")
		assert.Equal(t, "example.com", r.Header.Get("X-Forwarded-Host"), "X-Forwarded-Host should be set")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"metadata": "test"}`))
	}))
	defer testServer.Close()

	// Update the app's metadataUrl to point to our test server
	app.metadataUrl = testServer.URL

	// Create a test request
	req := httptest.NewRequest("GET", "/metadata/test", nil)
	req.Host = "example.com"

	// Create a response recorder
	rr := httptest.NewRecorder()

	// Call the metadataHandler
	handler := app.metadataHandler(app.metadataUrl)
	handler(rr, req)

	// Check the response
	assert.Equal(t, http.StatusOK, rr.Code, "Handler returned wrong status code")
}
