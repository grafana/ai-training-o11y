// File copied over from machine-learning/internal/middleware/middleware_test.go

package middleware

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grafana/dskit/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestResponseMiddleware(t *testing.T) {
	requestMiddlware := RequestResponseMiddleware(nil)

	t.Run("NoTenantHeader", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			return nil, errors.New("internal server error")
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, string(data), `{"status":"error","error":"no org id"}`)
	})

	t.Run("BadRequest", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req = req.WithContext(user.InjectOrgID(req.Context(), "mytenant"))

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			assert.Equal(t, "mytenant", tenant)
			return nil, ErrBadRequest(errors.New("bad request"))
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Contains(t, string(data), `{"status":"error","error":"bad request"}`)
	})

	t.Run("NotFound", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req = req.WithContext(user.InjectOrgID(req.Context(), "mytenant"))

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			assert.Equal(t, "mytenant", tenant)
			return nil, ErrNotFound(errors.New("not found"))
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, res.StatusCode)
		assert.Contains(t, string(data), `{"status":"error","error":"not found"}`)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req = req.WithContext(user.InjectOrgID(req.Context(), "mytenant"))

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			assert.Equal(t, "mytenant", tenant)
			return nil, errors.New("internal server error")
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Contains(t, string(data), `{"status":"error","error":"internal server error"}`)
	})

	t.Run("SuccessWithBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req = req.WithContext(user.InjectOrgID(req.Context(), "mytenant"))

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			assert.Equal(t, "mytenant", tenant)
			return map[string]interface{}{"hello": "world"}, nil
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Contains(t, string(data), `{"status":"success","data":{"hello":"world"}}`)
	})

	t.Run("SuccessWithoutBody", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/hello", nil)
		req = req.WithContext(user.InjectOrgID(req.Context(), "mytenant"))

		requestMiddlware(func(tenant string, req *http.Request) (interface{}, error) {
			assert.Equal(t, "mytenant", tenant)
			return nil, nil
		})(w, req)

		res := w.Result()
		defer res.Body.Close()
		data, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		assert.Empty(t, data)
	})
}
