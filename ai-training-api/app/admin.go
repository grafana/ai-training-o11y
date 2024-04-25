package api

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

//go:embed templates/processes.html
var processesPageContent string

//go:embed templates/style.css
var stylesheetContent string

type Admin struct {
	app                   *App
	processesPageTemplate *template.Template
}

func NewAdmin(app *App) *Admin {
	funcs := template.FuncMap{
		"prettyjson": func(x any) any {
			b, err := json.MarshalIndent(x, "", "  ")
			if err != nil {
				return x
			}
			return string(b)
		},
		"stylesheet": func() string { return stylesheetContent },
	}

	tProcesses := template.New("processes")
	tProcesses.Funcs(funcs)
	processesPageTemplate := template.Must(tProcesses.Parse(processesPageContent))

	return &Admin{
		app:                   app,
		processesPageTemplate: processesPageTemplate,
	}
}

func (admin *Admin) Register(router *mux.Router) {
	router.HandleFunc("/processes", admin.processes).Methods("GET")
}

type listResponse[T any] struct {
	Items    []T
	Limit    int
	TenantID string
}

func (a *Admin) processes(w http.ResponseWriter, req *http.Request) {
	var err error
	db := a.app.db(req.Context())

	type processResponse struct {
		listResponse[model.Process]
	}
	resp := processResponse{}

	query := req.URL.Query()
	l := query.Get("limit")
	resp.Limit, err = strconv.Atoi(l)
	if err != nil {
		resp.Limit = 20
	}
	q := db.Unscoped().
		Order("end_time desc, start_time desc").
		Limit(resp.Limit)

	resp.TenantID = query.Get("tenant")
	if resp.TenantID != "" {
		q = q.Where(map[string]string{"tenant_id": resp.TenantID})
	}

	resp.Items = make([]model.Process, 0, resp.Limit)
	err = q.Find(&resp.Items).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = a.processesPageTemplate.Execute(w, resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
