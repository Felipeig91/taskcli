// Package application contiene los casos de uso de taskcli.
package application

import (
	"context"

	"taskcli/domain"
)

// Repository es el contrato que necesita la aplicación para guardar tareas.
//
// La aplicación depende de esta interfaz, no de JSON. Esto permite sustituir
// el almacenamiento por un fake en tests o por otra tecnología en el futuro.
type Repository interface {
	List(ctx context.Context) ([]domain.Task, error)
	Save(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id string) error
}
