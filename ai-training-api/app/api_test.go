package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/google/uuid"
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
		"user_metadata": {
			"key1": "value1",
			"key2": 2
		}
	}`
	sampleProcessWithGroupNameJSON = `{
		"group": "group1",
		"user_metadata": {
			"key1": "value1"
		}
	}`
	sampleUpdateMetadataJSON = `{
		"user_metadata": {
			"key1": "completely_different_value",
			"key3": "value3"
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

type getGroupsResponse struct {
	middleware.ResponseWrapper
	Data []model.Group `json:"data"`
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
	testApp, err := New(listenAddress, listenPort, filepath.Join(t.TempDir(), "test.db"), db.SQLite, "0", "", &promlog.Config{Level: logLevel, Format: logFormat})
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
	assert.Len(t, gpr.Data.Metadata, 2)
	assert.Equal(t, "key1", gpr.Data.Metadata[0].Key)
	assert.Equal(t, "string", gpr.Data.Metadata[0].Type)
	assert.Equal(t, "value1", string(gpr.Data.Metadata[0].Value))
	assert.Equal(t, "key2", gpr.Data.Metadata[1].Key)
	assert.Equal(t, "int", gpr.Data.Metadata[1].Type)
	value, err := model.UnmarshalMetadataValue(gpr.Data.Metadata[1].Value, gpr.Data.Metadata[1].Type)
	require.NoError(t, err)
	assert.Equal(t, 2, value)
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
	assert.Len(t, ggr.Data.Processes, 1)
	assert.Equal(t, gpr.Data.ID, ggr.Data.Processes[0].ID)
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
	assert.Len(t, gpr.Data.Metadata, 3)
	assert.Contains(t, gpr.Data.Metadata, model.MetadataKV{TenantID: "0", Key: "key1", Type: "string", Value: []byte("completely_different_value"), ProcessID: cpr.Data.ID})
	assert.Contains(t, gpr.Data.Metadata, model.MetadataKV{TenantID: "0", Key: "key2", Type: "int", Value: []byte("2"), ProcessID: cpr.Data.ID})
	assert.Contains(t, gpr.Data.Metadata, model.MetadataKV{TenantID: "0", Key: "key3", Type: "string", Value: []byte("value3"), ProcessID: cpr.Data.ID})
}

func TestAppCreatesAndDeletesProcessAndGroup(t *testing.T) {
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

	// Delete the process.
	deleteProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String() + "/delete"
	req, err := http.NewRequest(http.MethodPost, deleteProcessEndpoint, nil)
	require.NoError(t, err)
	resp, err = httpC.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify the process was deleted.
	getProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + cpr.Data.ID.String()
	resp, err = httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)

	// Delete the group.
	deleteGroupEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/group/" + cpr.Data.GroupID.String() + "/delete"
	req, err = http.NewRequest(http.MethodPost, deleteGroupEndpoint, nil)
	require.NoError(t, err)
	resp, err = httpC.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNoContent, resp.StatusCode)

	// Verify the group was deleted.
	getGroupEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/group/" + cpr.Data.GroupID.String()
	resp, err = httpC.Get(getGroupEndpoint)
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestAppSetsCorrectEndTime(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)
	defer testApp.Shutdown()

	// Create a process with an old start time directly in the DB.
	startTime := time.Now().Add(-2 * time.Hour)
	process := model.Process{
		ID:        uuid.New(),
		TenantID:  "0",
		StartTime: startTime,
	}
	db := testApp.db(context.Background())
	require.NoError(t, db.Create(&process).Error)

	// Query the process and verify that the end time was set correctly.
	httpC := newHTTPClient(t.Name())
	getProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/" + process.ID.String()
	resp, err := httpC.Get(getProcessEndpoint)
	require.NoError(t, err)
	gpr := read[getProcessResponse](t, resp)
	assert.Equal(t, process.ID, gpr.Data.ID)
	assert.Equal(t, gpr.Data.EndTime, gpr.Data.StartTime.Add(time.Hour))
}

// This tests for the case where the same group name is added to two process
// metadatas. The group should be created only once and both processes should
// be added to the group.
func TestAppAddsProcessesToAGroup(t *testing.T) {
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

	// Create a new process and add it to the same group.
	resp, err = httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessWithGroupNameJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	cpr2 := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr2.Data.ID)

	// Get all groups and verify that the processes were added to the group.
	getGroupsEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/groups"
	resp, err = httpC.Get(getGroupsEndpoint)
	require.NoError(t, err)
	ggsr := read[getGroupsResponse](t, resp)
	assert.Len(t, ggsr.Data, 1)
	assert.Len(t, ggsr.Data[0].Processes, 2)
	assert.Contains(t, ggsr.Data[0].Processes, cpr.Data)
	assert.Contains(t, ggsr.Data[0].Processes, cpr2.Data)
}

// This tests for the case where two processes were created independently and
// then added to a group. This path is likely to be triggered via the UI.
func TestAppCreatesGroupWithMultipleProcesses(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)
	defer testApp.Shutdown()

	// Create two processes without group names.
	httpC := newHTTPClient(t.Name())
	registerProcessEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/process/new"
	resp, err := httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessNestedJSON))
	require.NoError(t, err)
	defer resp.Body.Close()
	cpr := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr.Data.ID)

	resp, err = httpC.Post(registerProcessEndpoint, "application/json", bytes.NewBufferString(sampleProcessNestedJSON))
	require.NoError(t, err)
	defer resp.Body.Close()
	cpr2 := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr2.Data.ID)

	// Create a group with the two processes.
	createGroupEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/group/new"
	resp, err = httpC.Post(createGroupEndpoint, "application/json", bytes.NewBufferString(`{"process_ids": ["`+cpr.Data.ID.String()+`", "`+cpr2.Data.ID.String()+`"]}`))
	require.NoError(t, err)
	defer resp.Body.Close()

	// Get all groups and verify that the processes were added to the group.
	getGroupsEndpoint := "http://" + testApp.server.HTTPListenAddr().String() + "/api/v1/groups"
	resp, err = httpC.Get(getGroupsEndpoint)
	require.NoError(t, err)
	ggsr := read[getGroupsResponse](t, resp)
	assert.Len(t, ggsr.Data, 1)
	assert.Len(t, ggsr.Data[0].Processes, 2)
	assert.Equal(t, ggsr.Data[0].Processes[0].ID, cpr.Data.ID)
	assert.Equal(t, ggsr.Data[0].Processes[1].ID, cpr2.Data.ID)
}
