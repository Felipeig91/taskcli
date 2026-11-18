// Package storage contiene implementaciones de persistencia.
package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"taskcli/domain"
)

var ErrTaskNotFound = errors.New("tarea no encontrada")

// JSONRepository guarda todas las tareas en un único archivo JSON.
type JSONRepository struct {
	path string
	mu   sync.Mutex
}

func NewJSONRepository(path string) *JSONRepository {
	return &JSONRepository{path: path}
}

func (r *JSONRepository) List(ctx context.Context) ([]domain.Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := checkContext(ctx); err != nil {
		return nil, err
	}
	tasks, err := r.read()
	if errors.Is(err, os.ErrNotExist) {
		return []domain.Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("leer %s: %w", r.path, err)
	}
	return tasks, nil
}

func (r *JSONRepository) Save(ctx context.Context, task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := checkContext(ctx); err != nil {
		return err
	}
	tasks, err := r.read()
	if errors.Is(err, os.ErrNotExist) {
		tasks = []domain.Task{}
	} else if err != nil {
		return fmt.Errorf("leer %s: %w", r.path, err)
	}

	updated := false
	for i := range tasks {
		if tasks[i].ID == task.ID {
			tasks[i] = task
			updated = true
			break
		}
	}
	if !updated {
		tasks = append(tasks, task)
	}
	return r.write(tasks)
}

func (r *JSONRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if err := checkContext(ctx); err != nil {
		return err
	}
	tasks, err := r.read()
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("%w: %s", ErrTaskNotFound, id)
	}
	if err != nil {
		return fmt.Errorf("leer %s: %w", r.path, err)
	}

	filtered := tasks[:0]
	found := false
	for _, task := range tasks {
		if task.ID == id {
			found = true
			continue
		}
		filtered = append(filtered, task)
	}
	if !found {
		return fmt.Errorf("%w: %s", ErrTaskNotFound, id)
	}
	return r.write(filtered)
}

func (r *JSONRepository) read() ([]domain.Task, error) {
	data, err := os.ReadFile(r.path)
	if err != nil {
		return nil, err
	}
	var tasks []domain.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("decodificar JSON: %w", err)
	}
	return tasks, nil
}

func (r *JSONRepository) write(tasks []domain.Task) error {
	if err := os.MkdirAll(filepath.Dir(r.path), 0755); err != nil {
		return fmt.Errorf("crear directorio: %w", err)
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("codificar JSON: %w", err)
	}
	if err := os.WriteFile(r.path, data, 0600); err != nil {
		return fmt.Errorf("escribir %s: %w", r.path, err)
	}
	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
