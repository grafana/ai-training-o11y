package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/grafana/ai-training-o11y/ai-training-api/internal"
	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
	"github.com/grafana/ai-training-o11y/ai-training-api/testutil"
)

const (
	listenAddress = "localhost"
	listenPort    = 0

	sampleProcessNestedJSON = `{
		"metadata": {
			"key1": "value1"
		}
	}`
	sampleProcessWithGroupNameJSON = `{
		"group": "group1",
		"metadata": {
			"key1": "value1"
		}
	}`
	sampleUpdateMetadataJSON = `{
		"metadata": {
			"key1": "completely_different_value",
			"key2": "value2"
		}
	}`
)

type createProcessResponse struct {
	middleware.ResponseWrapper
	Data model.Process `json:"data"`
}

type getProcessResponse struct {
	middleware.ResponseWrapper
	Data model.Process `json:"data"`
}

type getGroupResponse struct {
	middleware.ResponseWrapper
	Data model.Group `json:"data"`
}

func read[T any](t *testing.T, resp *http.Response) T {
	t.Helper()
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode, string(body))
	var v T
	require.NoError(t, json.Unmarshal(body, &v))
	return v
}

func newHTTPClient(tenant string) *http.Client {
	return &http.Client{
		Transport: testutil.NewTenantRoundTripper(http.DefaultTransport, tenant),
		Timeout:   time.Second * 5,
	}
}

func NewTestApp(t *testing.T, logger log.Logger) *App {
	logLevel := &promlog.AllowedLevel{}
	logLevel.Set("debug")
	logFormat := &promlog.AllowedFormat{}
	logFormat.Set("logfmt")
	testApp, err := New(listenAddress, listenPort, filepath.Join(t.TempDir(), "test.db"), db.SQLite, "0", &promlog.Config{Level: logLevel, Format: logFormat})
	require.NoError(t, err)
	// Run the server in parallel
	go testApp.Run()
	return testApp
}

func TestAppCreatesNewProcess(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)
	defer testApp.Shutdown()

	httpC := newHTTPClient(t.Name())
	registerProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/new"
	resp, err := httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessNestedJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	cpr := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr.Data.ID)

	// Verify the process was created.
	getProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String()
	resp, err = httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	gpr := read[getProcessResponse](t, resp)
	assert.Equal(t, cpr.Data.ID, gpr.Data.ID)
	assert.Len(t, gpr.Data.Metadata, 1)
	assert.Equal(t, "key1", gpr.Data.Metadata[0].Key)
	assert.Equal(t, "value1", gpr.Data.Metadata[0].Value)
}

func TestAppCreatesNewProcessAndGroup(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)
	defer testApp.Shutdown()

	httpC := newHTTPClient(t.Name())
	registerProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/new"
	resp, err := httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessWithGroupNameJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	cpr := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr.Data.ID)

	// Verify the process was created.
	getProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String()
	resp, err = httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	gpr := read[getProcessResponse](t, resp)
	assert.Equal(t, cpr.Data.ID, gpr.Data.ID)

	// Verify the group was created.
	getGroupEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/group/" + gpr.Data.GroupID.String()
	resp, err = httpC.Get(getGroupEndpoint)
	require.NoError(t, err)
	ggr := read[getGroupResponse](t, resp)
	assert.Equal(t, gpr.Data.GroupID, &ggr.Data.ID)
}

func TestAppCreatesAndUpdatesMetadata(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)
	defer testApp.Shutdown()

	httpC := newHTTPClient(t.Name())
	registerProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/new"
	resp, err := httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessNestedJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	cpr := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr.Data.ID)

	// Verify the process was created.
	getProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String()
	resp, err = httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	gpr := read[getProcessResponse](t, resp)
	assert.Equal(t, cpr.Data.ID, gpr.Data.ID)

	// Update the metadata.
	updateMetadataEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String() + "/update-metadata"
	resp, err = httpC.Post(updateMetadataEndpoint, "application/json", bytes.NewBufferString(sampleUpdateMetadataJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Verify the metadata was updated.
	resp, err = httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	gpr = read[getProcessResponse](t, resp)
	assert.Equal(t, cpr.Data.ID, gpr.Data.ID)
	assert.Len(t, gpr.Data.Metadata, 2)
	assert.Equal(t, "key1", gpr.Data.Metadata[0].Key)
	assert.Equal(t, "completely_different_value", gpr.Data.Metadata[0].Value)
	assert.Equal(t, "key2", gpr.Data.Metadata[1].Key)
	assert.Equal(t, "value2", gpr.Data.Metadata[1].Value)
}
