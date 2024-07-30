package plugin

import (
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

// handlePing is an HTTP GET resource that returns a {"message": "ok"} JSON response.
func (a *App) handlePing(w http.ResponseWriter, req *http.Request) {
	log.DefaultLogger.Debug("Handling ping request", "method", req.Method, "url", req.URL.String())
	w.Header().Add("Content-Type", "application/json")
	if _, err := w.Write([]byte(`{"message": "ok"}`)); err != nil {
		log.DefaultLogger.Error("Failed to write ping response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.DefaultLogger.Debug("Ping request handled successfully")
}

// handleEcho is an HTTP POST resource that accepts a JSON with a "message" key and
// returns to the client whatever it is sent.
func (a *App) handleEcho(w http.ResponseWriter, req *http.Request) {
	log.DefaultLogger.Debug("Handling echo request", "method", req.Method, "url", req.URL.String())
	if req.Method != http.MethodPost {
		log.DefaultLogger.Warn("Method not allowed for echo", "method", req.Method)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		log.DefaultLogger.Error("Failed to decode echo request body", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.DefaultLogger.Error("Failed to encode echo response", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	log.DefaultLogger.Debug("Echo request handled successfully", "message", body.Message)
}

func (a *App) metadataHandler(target string) func(http.ResponseWriter, *http.Request) {
	log.DefaultLogger.Info("Creating metadata handler", "target", target)
	remote, err := url.Parse(target)
	if err != nil {
		log.DefaultLogger.Error("Failed to parse metadata URL", "error", err, "url", target)
	}
	p := httputil.NewSingleHostReverseProxy(remote)

	// Modify the director to add the bearer token and X-Forwarded-Host
	originalDirector := p.Director
	p.Director = func(req *http.Request) {
		originalDirector(req)
		
		// Set X-Forwarded-Host
		if req.Host != "" {
			req.Header.Set("X-Forwarded-Host", req.Host)
		}
		
		// Set Authorization header if token is available
		if a.metadataToken != "" && a.stackId != "" { // Change belongs here
			req.Header.Set("Authorization", "Bearer "+ a.stackId + ":" + a.metadataToken)
			log.DefaultLogger.Debug("Added bearer token to request")
		} else {
			log.DefaultLogger.Debug("No metadata token set, not adding to request")
		}

		// Remove the "/metadata" prefix from the path
		req.URL.Path = req.URL.Path[len("/metadata"):]
		
		log.DefaultLogger.Debug("Proxying metadata request", 
			"method", req.Method, 
			"url", req.URL.String(), 
			"host", req.Host,
			"x-forwarded-host", req.Header.Get("X-Forwarded-Host"))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		log.DefaultLogger.Debug("Handling metadata request", "method", r.Method, "url", r.URL.String())
		p.ServeHTTP(w, r)
	}
}

// registerRoutes takes a *http.ServeMux and registers some HTTP handlers.
func (a *App) registerRoutes(mux *http.ServeMux) {
	log.DefaultLogger.Info("Registering routes")
	mux.HandleFunc("/ping", a.handlePing)
	mux.HandleFunc("/echo", a.handleEcho)
	mux.HandleFunc("/metadata/", a.metadataHandler(a.metadataUrl)) // Use config for this
	log.DefaultLogger.Info("Routes registered successfully")
}
