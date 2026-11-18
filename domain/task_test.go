package domain

import (
	"errors"
	"strings"
	"testing"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name      string
		id        string
		title     string
		wantErr   error
		wantID    string
		wantTitle string
	}{
		{
			name:      "crea y normaliza una tarea",
			id:        " task-1 ",
			title:     "  estudiar structs  ",
			wantID:    "task-1",
			wantTitle: "estudiar structs",
		},
		{
			name:    "rechaza ID vacío",
			title:   "estudiar structs",
			wantErr: ErrInvalidID,
		},
		{
			name:    "rechaza título vacío",
			id:      "task-1",
			title:   "   ",
			wantErr: ErrEmptyTitle,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task, err := NewTask(tt.id, tt.title)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewTask() error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr != nil {
				return
			}

			if task.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", task.ID, tt.wantID)
			}
			if task.Title != tt.wantTitle {
				t.Errorf("Title = %q, want %q", task.Title, tt.wantTitle)
			}
			if task.CreatedAt.IsZero() {
				t.Error("CreatedAt no debe ser cero")
			}
			if task.Completed {
				t.Error("una tarea nueva no debe estar completada")
			}
		})
	}
}

func TestTaskComplete(t *testing.T) {
	task, err := NewTask("task-1", "estudiar punteros")
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}

	if err := task.Complete(); err != nil {
		t.Fatalf("Complete() error = %v", err)
	}
	if !task.Completed {
		t.Error("Complete() no marcó la tarea como completada")
	}

	// La segunda llamada verifica la propiedad idempotente.
	if err := task.Complete(); err != nil {
		t.Fatalf("Complete() segunda llamada error = %v", err)
	}
}

func TestTaskCompleteNil(t *testing.T) {
	var task *Task

	err := task.Complete()
	if err == nil {
		t.Fatal("Complete() sobre nil debería devolver un error")
	}
	if !strings.Contains(err.Error(), "nil") {
		t.Fatalf("error = %q, debería explicar el puntero nil", err)
	}
}

func TestTaskString(t *testing.T) {
	task, err := NewTask("task-1", "leer interfaces")
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}

	if got, want := task.String(), "[pendiente] task-1: leer interfaces"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}

	task.Completed = true
	if got, want := task.String(), "[completada] task-1: leer interfaces"; got != want {
		t.Errorf("String() completada = %q, want %q", got, want)
	}
}

func TestTaskUpdateTitle(t *testing.T) {
	task, err := NewTask("task-1", "título original")
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	if err := task.UpdateTitle("  título editado  "); err != nil {
		t.Fatalf("UpdateTitle() error = %v", err)
	}
	if task.Title != "título editado" {
		t.Errorf("Title = %q, want título editado", task.Title)
	}
	if err := task.UpdateTitle(" "); !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("UpdateTitle() error = %v, want ErrEmptyTitle", err)
	}
}

func TestTaskAddComment(t *testing.T) {
	task, err := NewTask("task-1", "título")
	if err != nil {
		t.Fatalf("NewTask() error = %v", err)
	}
	if err := task.AddComment("comment-1", "  primer comentario  "); err != nil {
		t.Fatalf("AddComment() error = %v", err)
	}
	if len(task.Comments) != 1 {
		t.Fatalf("len(Comments) = %d, want 1", len(task.Comments))
	}
	if task.Comments[0].Body != "primer comentario" {
		t.Errorf("Body = %q, want primer comentario", task.Comments[0].Body)
	}
	if err := task.AddComment("comment-2", " "); !errors.Is(err, ErrEmptyComment) {
		t.Fatalf("AddComment() error = %v, want ErrEmptyComment", err)
	}
}
