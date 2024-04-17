package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/go-kit/log"
	"github.com/prometheus/common/promlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	db "github.com/grafana/ai-training-o11y/ai-training-api/internal"
	"github.com/grafana/ai-training-o11y/ai-training-api/middleware"
	"github.com/grafana/ai-training-o11y/ai-training-api/testutil"
)

const (
	listenAddress = "localhost"
	listenPort    = 8000

	sampleProcessJSON = `{
		"name": "Test Process"
	}`
)

var (
	apiEndpointPrefix = "http://" + listenAddress + ":" + strconv.Itoa(listenPort) + "/api/v1"
)

type createProcessResponse struct {
	middleware.ResponseWrapper
	Data CreateProcessResponse `json:"data"`
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
	return testApp
}

func TestAppCreatesNewProcess(t *testing.T) {
	logger := log.NewNopLogger()
	testApp := NewTestApp(t, logger)
	require.NotNil(t, testApp)

	httpC := newHTTPClient(t.Name())
	resp, err := httpC.Post(apiEndpointPrefix+"/process/new", "application/json", bytes.NewBufferString(sampleProcessJSON))
	require.NoError(t, err)
	defer resp.Body.Close()

	cpr := read[createProcessResponse](t, resp)
	assert.NotEmpty(t, cpr.Data.ID)
}
