// File copied over from machine-learning/internal/middleware/middleware.go

package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	"github.com/grafana/dskit/middleware"
	"github.com/grafana/dskit/user"
)

var AuthenticateUser = middleware.AuthenticateUser

func AuthnMiddleware(constTenant string) mux.MiddlewareFunc {
	// This will be true in prod.
	if constTenant == "" {
		return mux.MiddlewareFunc(middleware.AuthenticateUser)
	}

	// When we are using a constant tenant we need to inject the constant
	// tenant into the context for the request.
	// NOTE: THIS ONLY HAPPENS IN DEV!
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := user.InjectOrgID(r.Context(), constTenant)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ContentTypeJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Request func(tenant string, req *http.Request) (interface{}, error)

func RequestResponseMiddleware(logger log.Logger) func(Request) http.HandlerFunc {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	return func(f Request) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			// The OrgID needs to be added to the router before we reach this middleware.
			tenantID, err := user.ExtractOrgID(req.Context())
			if err != nil {
				level.Error(logger).Log("msg", "Error in api request", "err", err)
				resp, _ := json.Marshal(ResponseWrapper{
					Status: "error",
					Error:  err.Error(),
				})
				http.Error(w, string(resp), http.StatusInternalServerError)
				return
			}

			data, err := f(tenantID, req)
			if err != nil {
				statusCode := errorStatusCode(err)
				level.Error(logger).Log("msg", "Error in api request", "err", err, "code", statusCode)
				resp, _ := json.Marshal(ResponseWrapper{
					Status: "error",
					Error:  err.Error(),
				})
				http.Error(w, string(resp), statusCode)
				return
			}

			// If there is no data return a 204 instead of a 200.
			if data == nil {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			res, err := json.Marshal(ResponseWrapper{
				Status: "success",
				Data:   data,
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			//nolint:errcheck // Just do our best to write.
			w.Write(res)
		}
	}
}

type ResponseWrapper struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

type errNotFound struct{ error }

func ErrNotFound(err error) error {
	return errNotFound{err}
}

type errBadRequest struct{ error }

func ErrBadRequest(err error) error {
	return errBadRequest{err}
}

func errorStatusCode(err error) int {
	switch err {
	case context.Canceled:
		return http.StatusBadRequest
	case context.DeadlineExceeded:
		return http.StatusRequestTimeout
	}
	switch err.(type) {
	case errNotFound:
		return http.StatusNotFound
	case errBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
