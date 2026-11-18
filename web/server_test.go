package web

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"taskcli/application"
	"taskcli/storage"
)

func TestServerCreatesAndListsTask(t *testing.T) {
	repository := storage.NewJSONRepository(t.TempDir() + "/tasks.json")
	service := application.NewService(repository)
	server, err := NewServer(service)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	createRequest := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("title=aprender+HTTP"))
	createRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(createResponse, createRequest)

	if createResponse.Code != http.StatusSeeOther {
		t.Fatalf("POST /tasks status = %d, want %d", createResponse.Code, http.StatusSeeOther)
	}

	listResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	if listResponse.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", listResponse.Code, http.StatusOK)
	}
	if !strings.Contains(listResponse.Body.String(), "aprender HTTP") {
		t.Error("la página no contiene la tarea creada")
	}
}

func TestServerEditsAndCommentsTask(t *testing.T) {
	repository := storage.NewJSONRepository(t.TempDir() + "/tasks.json")
	service := application.NewService(repository)
	server, err := NewServer(service)
	if err != nil {
		t.Fatalf("NewServer() error = %v", err)
	}

	createRequest := httptest.NewRequest(http.MethodPost, "/tasks", strings.NewReader("title=titulo+original"))
	createRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	createResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(createResponse, createRequest)
	if createResponse.Code != http.StatusSeeOther {
		t.Fatalf("POST /tasks status = %d, want %d", createResponse.Code, http.StatusSeeOther)
	}

	tasks, err := service.List(createRequest.Context())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	id := tasks[0].ID

	editRequest := httptest.NewRequest(http.MethodPost, "/tasks/"+id+"/edit", strings.NewReader("title=titulo+editado"))
	editRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	editResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(editResponse, editRequest)
	if editResponse.Code != http.StatusSeeOther {
		t.Fatalf("POST edit status = %d, want %d", editResponse.Code, http.StatusSeeOther)
	}

	commentRequest := httptest.NewRequest(http.MethodPost, "/tasks/"+id+"/comments", strings.NewReader("comment=mi+comentario"))
	commentRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	commentResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(commentResponse, commentRequest)
	if commentResponse.Code != http.StatusSeeOther {
		t.Fatalf("POST comment status = %d, want %d", commentResponse.Code, http.StatusSeeOther)
	}

	pageResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(pageResponse, httptest.NewRequest(http.MethodGet, "/", nil))
	body := pageResponse.Body.String()
	if !strings.Contains(body, "titulo editado") {
		t.Error("la página no contiene el título editado")
	}
	if !strings.Contains(body, "mi comentario") {
		t.Error("la página no contiene el comentario")
	}
}
