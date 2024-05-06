package api

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/grafana/ai-training-o11y/ai-training-api/model"
)

//go:embed templates/processes.html
var processesPageContent string

//go:embed templates/process.html
var processPageContent string

//go:embed templates/style.css
var stylesheetContent string

type Admin struct {
	app                   *App
	processesPageTemplate *template.Template
	processPageTemplate   *template.Template
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

	tProcess := template.New("process")
	tProcess.Funcs(funcs)
	processPageTemplate := template.Must(tProcess.Parse(processPageContent))

	return &Admin{
		app:                   app,
		processesPageTemplate: processesPageTemplate,
		processPageTemplate:   processPageTemplate,
	}
}

func (admin *Admin) Register(router *mux.Router) {
	router.HandleFunc("/processes", admin.processes).Methods("GET")
	router.HandleFunc("/process/{id}", admin.process).Methods("GET")
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

func (a *Admin) process(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	processId := namedParam(req, "id")
	query := req.URL.Query()
	tenantID := query.Get("tenant")

	process := model.Process{}
	err := a.app.db(ctx).
		Unscoped().
		Where(&model.Process{ID: uuid.MustParse(processId), TenantID: tenantID}).
		Find(&process).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// find all associated metadata
	err = a.app.db(ctx).
		Where(model.MetadataKV{TenantID: tenantID, ProcessID: uuid.MustParse(processId)}).
		Find(&process.Metadata).Error
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	type processResponse struct {
		Process model.Process
	}

	data := processResponse{process}
	err = a.processPageTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
