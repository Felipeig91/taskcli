// Package domain contiene las reglas de negocio de taskcli.
//
// Este paquete no conoce la terminal, archivos JSON ni bases de datos. Esa
// separación permite probar las reglas rápidamente y cambiar la infraestructura
// sin modificar el modelo de la tarea.
package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	// ErrInvalidID identifica una tarea sin un identificador usable.
	ErrInvalidID = errors.New("el ID de la tarea no puede estar vacío")
	// ErrEmptyTitle identifica una tarea sin título.
	ErrEmptyTitle = errors.New("el título de la tarea no puede estar vacío")
	// ErrEmptyComment identifica un comentario sin contenido.
	ErrEmptyComment = errors.New("el comentario no puede estar vacío")
)

// Comment representa un mensaje dentro del hilo de una tarea.
type Comment struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

// Task representa una tarea del dominio.
//
// Los campos exportados permiten que las capas exteriores lean la tarea y,
// más adelante, puedan serializarla. Las modificaciones con reglas de negocio
// pasan por métodos como Complete en vez de cambiar el estado sin control.
type Task struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
	Comments  []Comment
}

// NewTask crea una tarea válida.
//
// Recibir el ID desde fuera mantiene esta capa independiente de la estrategia
// de generación de IDs. En una fase posterior la aplicación podrá usar un
// contador, un UUID o un ID generado por el almacenamiento sin cambiar Task.
func NewTask(id, title string) (Task, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return Task{}, ErrInvalidID
	}

	title = strings.TrimSpace(title)
	if title == "" {
		return Task{}, ErrEmptyTitle
	}

	return Task{
		ID:        id,
		Title:     title,
		CreatedAt: time.Now(),
		Comments:  []Comment{},
	}, nil
}

// UpdateTitle cambia el título después de validar la nueva entrada.
func (t *Task) UpdateTitle(title string) error {
	if t == nil {
		return errors.New("no se puede editar una tarea nil")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return ErrEmptyTitle
	}
	t.Title = title
	return nil
}

// AddComment agrega un mensaje al hilo de la tarea.
func (t *Task) AddComment(id, body string) error {
	if t == nil {
		return errors.New("no se puede comentar una tarea nil")
	}
	id = strings.TrimSpace(id)
	if id == "" {
		return errors.New("el ID del comentario no puede estar vacío")
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return ErrEmptyComment
	}
	t.Comments = append(t.Comments, Comment{
		ID:        id,
		Body:      body,
		CreatedAt: time.Now(),
	})
	return nil
}

// Complete marca la tarea como completada.
//
// La operación es idempotente: completar una tarea ya completada sigue siendo
// una operación válida y no produce un error inesperado para la CLI.
func (t *Task) Complete() error {
	if t == nil {
		return errors.New("no se puede completar una tarea nil")
	}

	t.Completed = true
	return nil
}

// String ofrece una representación humana útil para logs y la futura CLI.
func (t Task) String() string {
	status := "pendiente"
	if t.Completed {
		status = "completada"
	}

	return fmt.Sprintf("[%s] %s: %s", status, t.ID, t.Title)
}
