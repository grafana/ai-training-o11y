package plugin

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewApp(t *testing.T) {
	tests := []struct {
		name                    string
		jsonData                map[string]interface{}
		decryptedSecureJSONData map[string]string
		expected                map[string]string
	}{
		{
			name: "All fields set correctly",
			jsonData: map[string]interface{}{
				"metadataUrl":         "https://example.com",
				"stackId":             "stack123",
				"lokiDatasourceName":  "loki",
				"mimirDatasourceName": "mimir",
			},
			decryptedSecureJSONData: map[string]string{
				"metadataToken": "token123",
			},
			expected: map[string]string{
				"metadataUrl":         "https://example.com",
				"stackId":             "stack123",
				"lokiDatasourceName":  "loki",
				"mimirDatasourceName": "mimir",
				"metadataToken":       "token123",
			},
		},
		{
			name: "Missing some fields",
			jsonData: map[string]interface{}{
				"metadataUrl":        "https://example.com",
				"lokiDatasourceName": "loki",
			},
			decryptedSecureJSONData: map[string]string{
				"metadataToken": "token456",
			},
			expected: map[string]string{
				"metadataUrl":         "https://example.com",
				"stackId":             "",
				"lokiDatasourceName":  "loki",
				"mimirDatasourceName": "",
				"metadataToken":       "token456",
			},
		},
		{
			name: "No data provided",
			expected: map[string]string{
				"metadataUrl":         "",
				"stackId":             "",
				"lokiDatasourceName":  "",
				"mimirDatasourceName": "",
				"metadataToken":       "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonData, err := json.Marshal(tt.jsonData)
			require.NoError(t, err)

			appSettings := backend.AppInstanceSettings{
				JSONData:                jsonData,
				DecryptedSecureJSONData: tt.decryptedSecureJSONData,
			}

			instance, err := NewApp(context.Background(), appSettings)
			require.NoError(t, err)
			require.NotNil(t, instance)

			app, ok := instance.(*App)
			require.True(t, ok)

			assert.Equal(t, tt.expected["metadataUrl"], app.metadataUrl)
			assert.Equal(t, tt.expected["stackId"], app.stackId)
			assert.Equal(t, tt.expected["lokiDatasourceName"], app.lokiDatasourceName)
			assert.Equal(t, tt.expected["mimirDatasourceName"], app.mimirDatasourceName)
			assert.Equal(t, tt.expected["metadataToken"], app.metadataToken)
		})
	}
}

func TestNewAppWithInvalidJSONData(t *testing.T) {
	appSettings := backend.AppInstanceSettings{
		JSONData: []byte(`invalid json`),
	}

	_, err := NewApp(context.Background(), appSettings)
	assert.Error(t, err)
}

func TestNewAppWithNonStringValues(t *testing.T) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"metadataUrl":         123,
		"stackId":             true,
		"lokiDatasourceName":  []string{"not", "a", "string"},
		"mimirDatasourceName": map[string]string{"not": "a string"},
	})
	require.NoError(t, err)

	appSettings := backend.AppInstanceSettings{
		JSONData: jsonData,
	}

	instance, err := NewApp(context.Background(), appSettings)
	require.NoError(t, err)
	require.NotNil(t, instance)

	app, ok := instance.(*App)
	require.True(t, ok)

	assert.Empty(t, app.metadataUrl)
	assert.Empty(t, app.stackId)
	assert.Empty(t, app.lokiDatasourceName)
	assert.Empty(t, app.mimirDatasourceName)
}
