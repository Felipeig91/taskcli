// Package web contiene la interfaz HTTP de taskcli.
package web

import (
	"embed"
	"html/template"
	"net/http"
	"strings"

	"taskcli/application"
	"taskcli/domain"
)

//go:embed templates/index.html static/styles.css
var assets embed.FS

type Server struct {
	service  *application.Service
	template *template.Template
}

func NewServer(service *application.Service) (*Server, error) {
	page, err := template.ParseFS(assets, "templates/index.html")
	if err != nil {
		return nil, err
	}
	return &Server{service: service, template: page}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/tasks", s.handleTasks)
	mux.HandleFunc("/tasks/", s.handleTasks)
	mux.Handle("/static/", http.FileServer(http.FS(assets)))
	return mux
}

type pageData struct {
	Tasks []domain.Task
	Error string
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	tasks, err := s.service.List(r.Context())
	if err != nil {
		s.render(w, http.StatusInternalServerError, pageData{Error: err.Error()})
		return
	}
	s.render(w, http.StatusOK, pageData{Tasks: tasks})
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/tasks")
	switch {
	case path == "":
		s.createTask(w, r)
	case strings.HasSuffix(path, "/edit"):
		s.editTask(w, r, strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/edit"))
	case strings.HasSuffix(path, "/comments"):
		s.addComment(w, r, strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/comments"))
	case strings.HasSuffix(path, "/complete"):
		s.updateTask(w, r, strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/complete"), true)
	case strings.HasSuffix(path, "/delete"):
		s.updateTask(w, r, strings.TrimSuffix(strings.TrimPrefix(path, "/"), "/delete"), false)
	default:
		http.NotFound(w, r)
	}
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	if _, err := s.service.Add(r.Context(), title); err != nil {
		s.renderWithError(w, r, http.StatusBadRequest, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) editTask(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.service.UpdateTitle(r.Context(), id, r.FormValue("title")); err != nil {
		s.renderWithError(w, r, http.StatusBadRequest, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) addComment(w http.ResponseWriter, r *http.Request, id string) {
	if err := s.service.AddComment(r.Context(), id, r.FormValue("comment")); err != nil {
		s.renderWithError(w, r, http.StatusBadRequest, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) updateTask(w http.ResponseWriter, r *http.Request, id string, complete bool) {
	var err error
	if complete {
		err = s.service.Complete(r.Context(), id)
	} else {
		err = s.service.Delete(r.Context(), id)
	}
	if err != nil {
		s.renderWithError(w, r, http.StatusBadRequest, err)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *Server) renderWithError(w http.ResponseWriter, r *http.Request, status int, err error) {
	tasks, listErr := s.service.List(r.Context())
	if listErr != nil {
		err = listErr
	}
	s.render(w, status, pageData{Tasks: tasks, Error: err.Error()})
}

func (s *Server) render(w http.ResponseWriter, status int, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := s.template.Execute(w, data); err != nil {
		// El encabezado ya se envió; no podemos cambiar el status, pero sí
		// dejamos que net/http registre el fallo de escritura.
		http.Error(w, "error al renderizar la página", http.StatusInternalServerError)
	}
}
