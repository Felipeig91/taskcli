package application

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"taskcli/domain"
)

var ErrTaskNotFound = errors.New("tarea no encontrada")

// Service coordina los casos de uso de la aplicación.
type Service struct {
	repository Repository
	now        func() time.Time
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
		now:        time.Now,
	}
}

// Add crea y persiste una tarea. El ID se genera aquí, fuera del dominio.
func (s *Service) Add(ctx context.Context, title string) (domain.Task, error) {
	id := "task-" + strconv.FormatInt(s.now().UnixNano(), 10)
	task, err := domain.NewTask(id, title)
	if err != nil {
		return domain.Task{}, err
	}
	if err := s.repository.Save(ctx, task); err != nil {
		return domain.Task{}, fmt.Errorf("guardar tarea: %w", err)
	}
	return task, nil
}

func (s *Service) List(ctx context.Context) ([]domain.Task, error) {
	tasks, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listar tareas: %w", err)
	}
	return tasks, nil
}

func (s *Service) Complete(ctx context.Context, id string) error {
	task, err := s.find(ctx, id)
	if err != nil {
		return err
	}
	if err := task.Complete(); err != nil {
		return err
	}
	if err := s.repository.Save(ctx, task); err != nil {
		return fmt.Errorf("guardar tarea completada: %w", err)
	}
	return nil
}

// UpdateTitle edita el título de una tarea existente.
func (s *Service) UpdateTitle(ctx context.Context, id, title string) error {
	task, err := s.find(ctx, id)
	if err != nil {
		return err
	}
	if err := task.UpdateTitle(title); err != nil {
		return err
	}
	if err := s.repository.Save(ctx, task); err != nil {
		return fmt.Errorf("guardar título editado: %w", err)
	}
	return nil
}

// AddComment agrega un comentario al hilo de una tarea existente.
func (s *Service) AddComment(ctx context.Context, id, body string) error {
	task, err := s.find(ctx, id)
	if err != nil {
		return err
	}
	commentID := "comment-" + strconv.FormatInt(s.now().UnixNano(), 10)
	if err := task.AddComment(commentID, body); err != nil {
		return err
	}
	if err := s.repository.Save(ctx, task); err != nil {
		return fmt.Errorf("guardar comentario: %w", err)
	}
	return nil
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("eliminar tarea: %w", err)
	}
	return nil
}

func (s *Service) find(ctx context.Context, id string) (domain.Task, error) {
	tasks, err := s.repository.List(ctx)
	if err != nil {
		return domain.Task{}, fmt.Errorf("buscar tarea: %w", err)
	}
	for _, task := range tasks {
		if task.ID == id {
			return task, nil
		}
	}
	return domain.Task{}, fmt.Errorf("%w: %s", ErrTaskNotFound, id)
}
