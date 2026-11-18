// Command taskcli será la entrada de nuestra herramienta de línea de comandos.
//
// En esta primera fase mantenemos el programa pequeño a propósito: el objetivo
// es tener un punto de partida ejecutable antes de añadir persistencia,
// interfaces y concurrencia.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"taskcli/application"
	"taskcli/storage"
	"taskcli/web"
)

// main delega la lógica que puede fallar en run.
//
// En Go es habitual que una función de nivel superior devuelva un error y que
// main decida cómo presentarlo. Así no ocultamos errores dentro de helpers y
// mantenemos una única política para terminar el proceso.
func main() {
	if err := run(context.Background(), os.Args[1:], os.Stdout, os.Stderr); err != nil {
		// Los errores van a stderr y un código distinto de cero permite que
		// scripts y la shell detecten que el comando falló.
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// run será el punto de composición de la aplicación. En fases posteriores
// recibirá configuración, construirá servicios y ejecutará el comando elegido.
// Por ahora solo demuestra que el esqueleto compila y se puede ejecutar.
func run(ctx context.Context, args []string, out, errOut *os.File) error {
	if len(args) == 0 || args[0] == "help" {
		printUsage(out)
		return nil
	}

	path := filepath.Join(".taskcli", "tasks.json")
	service := application.NewService(storage.NewJSONRepository(path))

	switch args[0] {
	case "web":
		server, err := web.NewServer(service)
		if err != nil {
			return fmt.Errorf("crear servidor web: %w", err)
		}
		fmt.Fprintln(out, "taskcli web escuchando en http://localhost:8080")
		return http.ListenAndServe(":8080", server.Handler())
	case "add":
		if len(args) < 2 {
			return fmt.Errorf("uso: taskcli add <título>")
		}
		task, err := service.Add(ctx, strings.Join(args[1:], " "))
		if err != nil {
			return err
		}
		fmt.Fprintln(out, task)
	case "list":
		tasks, err := service.List(ctx)
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			fmt.Fprintln(out, "no hay tareas")
			return nil
		}
		for _, task := range tasks {
			fmt.Fprintln(out, task)
		}
	case "complete":
		if len(args) != 2 {
			return fmt.Errorf("uso: taskcli complete <id>")
		}
		return service.Complete(ctx, args[1])
	case "delete":
		if len(args) != 2 {
			return fmt.Errorf("uso: taskcli delete <id>")
		}
		return service.Delete(ctx, args[1])
	case "edit":
		if len(args) < 3 {
			return fmt.Errorf("uso: taskcli edit <id> <nuevo título>")
		}
		return service.UpdateTitle(ctx, args[1], strings.Join(args[2:], " "))
	case "comment":
		if len(args) < 3 {
			return fmt.Errorf("uso: taskcli comment <id> <comentario>")
		}
		return service.AddComment(ctx, args[1], strings.Join(args[2:], " "))
	default:
		return fmt.Errorf("comando desconocido %q; usa \"taskcli help\"", args[0])
	}
	return nil
}

func printUsage(out *os.File) {
	fmt.Fprintln(out, "taskcli - gestor de tareas")
	fmt.Fprintln(out, "")
	fmt.Fprintln(out, "Uso:")
	fmt.Fprintln(out, "  taskcli add <título>")
	fmt.Fprintln(out, "  taskcli list")
	fmt.Fprintln(out, "  taskcli complete <id>")
	fmt.Fprintln(out, "  taskcli delete <id>")
	fmt.Fprintln(out, "  taskcli edit <id> <nuevo título>")
	fmt.Fprintln(out, "  taskcli comment <id> <comentario>")
	fmt.Fprintln(out, "  taskcli web")
}
