# taskcli

**taskcli** es una herramienta de gestion de tareas escrita desde cero en Go.
Incluye una interfaz de linea de comandos y una interfaz web moderna, reutiliza
la misma logica de aplicacion y persiste los datos localmente en JSON.

> Proyecto educativo creado para aprender Go mediante un producto real,
> incremental y testeable.

## Que se hizo

La primera version ya permite:

- Crear tareas desde la terminal o desde el navegador.
- Listar todas las tareas.
- Marcar una tarea como completada.
- Editar el titulo de una tarea.
- Eliminar una tarea.
- Agregar comentarios a cada tarea.
- Consultar un hilo de comentarios cronologico.
- Persistir tareas, estados y comentarios en `.taskcli/tasks.json`.
- Ejecutar un servidor web local en `http://localhost:8080`.
- Ejecutar tests unitarios y pruebas HTTP.
- Ejecutar el detector de carreras de Go.

Actualmente existen dos interfaces sobre el mismo servicio:

```text
Terminal ---> application.Service ---> JSONRepository
                                        |
Navegador --> HTTP Server ---------> application.Service
```

No hace falta una base de datos ni Node.js: este proyecto usa paquetes
estandar de Go para HTTP, HTML, JSON, archivos y pruebas.

## Para que sirve

`taskcli` sirve para:

- Practicar modelos de dominio con `struct`.
- Entender interfaces y sustitucion de dependencias.
- Separar CLI, aplicacion, dominio y almacenamiento.
- Aprender errores explicitos con `if err != nil`.
- Usar punteros y metodos con seguridad.
- Construir un servidor HTTP sin framework externo.
- Testear reglas, persistencia y respuestas HTTP.
- Persistir datos sin una base de datos externa.
- Conservar comentarios asociados a una tarea.

Ademas puede usarse como plantilla para otros programas locales: gestores de
notas, agendas, inventarios, incidencias o utilidades de linea de comandos.

## Que se podria llegar a hacer

Ideas naturales para evolucionarlo:

1. **Etiquetas y prioridades:** organizar tareas por categoria o importancia.
2. **Fechas limite y recordatorios:** agregar vencimientos y alertas.
3. **Busqueda y filtros:** filtrar pendientes, completadas y comentarios.
4. **API REST JSON:** exponer los mismos casos de uso a clientes externos.
5. **Base de datos:** sustituir `JSONRepository` por SQLite o PostgreSQL.
6. **Autenticacion:** usuarios, cuentas y sesiones.
7. **Multiusuario:** permisos y tareas compartidas.
8. **Sincronizacion:** almacenar datos de forma remota.
9. **Concurrencia real:** procesar importaciones, exportaciones y tareas
   independientes con goroutines y canales.
10. **Aplicacion movil o escritorio:** consumir una API mediante React,
    Flutter, Fyne o Wails.
11. **Publicacion Docker:** empaquetar el servidor en una imagen.
12. **CI/CD:** ejecutar tests, race detection y vet en GitHub Actions.

La arquitectura esta preparada para estos cambios: el servicio no depende
directamente de HTTP o JSON, sino de una interfaz de repositorio.

## Por que Go es una buena opcion

Go puede funcionar especialmente bien en este proyecto por varios motivos:

- **Binario unico:** `go build` produce un ejecutable sin runtime adicional.
- **Compilacion rapida:** ayuda a mantener un ciclo corto de desarrollo.
- **Concurrencia integrada:** goroutines y canales permiten crecer hacia
  servidores y procesamiento paralelo sin una capa de concurrencia externa.
- **Servidor incluido:** `net/http`, plantillas y servidores se encuentran
  dentro de la biblioteca estandar.
- **Herramientas integradas:** `gofmt`, `go test`, `go vet` y el detector de
  carreras forman parte del flujo habitual.
- **Errores visibles:** los errores se retornan y se procesan explicitamente.
- **Despliegue sencillo:** un binario ligero simplifica un contenedor o una
  ejecucion local.
- **Legibilidad consistente:** el formateador oficial reduce discusiones
  inutiles sobre estilo.

No significa que Go sea siempre mejor que Node.js, Python, Rust o Java:

- Python puede ser mas rapido para prototipos centrados en datos o IA.
- Node.js puede simplificar un proyecto donde todo el equipo usa TypeScript.
- Rust ofrece un control de memoria mas estricta en software de alto nivel
  de exigencia.
- Java ofrece ecosistemas y frameworks empresariales muy maduros.

Go destaca aqui por combinar una sintaxis pequena, ejecucion rapida,
concurrencia accesible y despliegue sencillo.

## Arquitectura

```text
main.go
  |
  +-- CLI
  |     +-- application.Service
  |
  +-- HTTP Server (web/)
        +-- application.Service
              |
              +-- domain.Task
              +-- domain.Comment
              |
              +-- Repository interface
                    |
                    +-- JSONRepository
                          |
                          +-- .taskcli/tasks.json
```

### Paquetes

- `domain`: entidades y reglas de negocio.
- `application`: casos de uso como crear, editar, comentar y completar.
- `storage`: implementacion persistente basada en JSON.
- `web`: servidor HTTP, plantilla HTML y estilos.
- `main.go`: seleccion del modo CLI o servidor web.

El dominio no conoce la terminal, HTTP ni JSON. Cada capa puede evolucionar
sin arrastrar necesariamente a las otras.

## Requisitos

- macOS, Linux o Windows.
- Go 1.23 o posterior.
- Navegador moderno.

En macOS, Go puede instalarse con Homebrew:

```bash
brew install go
```

Comprueba la instalacion:

```bash
go version
```

## Ejecutar el proyecto

Clona el repositorio y entra en la carpeta:

```bash
git clone https://github.com/TU_USUARIO/taskcli.git
cd taskcli
```

Ejecuta directamente:

```bash
go run .
```

El proyecto actual usa subcomandos, por lo que no necesita argumentos para
mostrar ayuda:

```bash
go run . help
```

Compila un ejecutable:

```bash
go build -o taskcli .
```

Usa el binario con:

```bash
./taskcli help
```

## Uso desde la terminal

Crear tareas:

```bash
./taskcli add "Aprender structs"
./taskcli add "Practicar interfaces"
```

Listar tareas:

```bash
./taskcli list
```

Completar una tarea:

```bash
./taskcli complete task-ID
```

Editar una tarea:

```bash
./taskcli edit task-ID "Nuevo titulo"
```

Agregar comentarios:

```bash
./taskcli comment task-ID "Este es el primer comentario"
```

Eliminar una tarea:

```bash
./taskcli delete task-ID
```

## Interfaz web

Inicia el servidor:

```bash
go run . web
```

Abre:

```text
http://localhost:8080
```

La interfaz permite agregar, editar, completar, eliminar y comentar tareas.
Deten el servidor con `Ctrl+C`.

Las rutas principales son:

| Metodo | Ruta | Funcion |
|---|---|---|
| `GET` | `/` | Mostrar tareas y comentarios |
| `POST` | `/tasks` | Crear una tarea |
| `POST` | `/tasks/{id}/edit` | Editar el titulo |
| `POST` | `/tasks/{id}/complete` | Completar una tarea |
| `POST` | `/tasks/{id}/comments` | Agregar un comentario |
| `POST` | `/tasks/{id}/delete` | Eliminar una tarea |

Los datos se guardan de forma local. El archivo generado no contiene
credenciales, pero normalmente no conviene publicarlo como parte del codigo
fuente.

## Ejemplo de datos

El almacenamiento actual guarda una estructura similar a:

```json
[
  {
    "ID": "task-123456",
    "Title": "Aprender interfaces",
    "Completed": false,
    "CreatedAt": "2026-11-18T12:00:00Z",
    "Comments": [
      {
        "ID": "comment-123457",
        "Body": "Revisar el repositorio y el servicio",
        "CreatedAt": "2026-11-18T12:05:00Z"
      }
    ]
  }
]
```

## Conceptos de Go utilizados

### Struct

Un `struct` agrupa datos relacionados. `Task` representa una tarea y
`Comment` representa un mensaje dentro de su hilo.

### Punteros

Los metodos con receptor `*Task` pueden modificar el objeto original. Por
ejemplo, `Complete` cambia `Completed` sin construir una nueva tarea manual.

### Interfaces

`application.Repository` define las operaciones necesarias para persistencia.
El servicio no conoce si los datos vienen de JSON, memoria o PostgreSQL.

### Errores idiomaticos

Go representa fallos con `error` y evita excepciones para los caminos normales:

```go
task, err := service.Add(ctx, title)
if err != nil {
    return err
}
```

Cuando una capa agrega contexto, usa `%w` para conservar la causa original.

### Concurrencia

El servidor HTTP maneja varias peticiones simultaneas. El repositorio JSON
actual usa `sync.Mutex` para proteger las operaciones dentro del mismo proceso.

### Pruebas

Los tests verifican reglas del dominio y respuestas HTTP con `httptest`.
Los tests de tabla permiten cubrir varios casos con una estructura clara.

## Calidad y verificacion

Desde la raiz del proyecto:

```bash
gofmt -w main.go application/*.go storage/*.go domain/*.go web/*.go
go test ./...
go test -race ./...
go vet ./...
```

Comandos rapidos:

```bash
# Solo tests
go test ./...

# Tests detectando posibles data races
go test -race ./...

# Analisis estatico
go vet ./...
```

## Dependencias

El proyecto principal usa exclusivamente la biblioteca estandar de Go:

- `net/http`
- `html/template`
- `embed`
- `encoding/json`
- `context`
- `sync`
- `testing`
- `httptest`

No hay dependencias externas ni framework HTTP.

## Flujo diario con VS Code

1. Ejecuta `go test ./...` despues de modificar una regla.
2. Pide a Copilot una propuesta pequena y revisa cada linea.
3. Usa sus sugerencias para tests adicionales, no como reemplazo del test.
4. Pide una revision de casos limite o errores de concurrencia.
5. Ejecuta `go test -race ./...` antes de subir cambios.
6. Usa el depurador de VS Code para reproducir una entrada concreta.
7. Cuando encuentres un problema, conviertelo en un test de regresion.

Ejemplo de prompt para Copilot:

```text
Genera casos de prueba para Editar y AgregarComentario.
Incluye titulos vacios, comentarios vacios, tareas inexistentes
y persistencia correcta. Explica cada caso sin modificar el dominio.
```

## Publicar en GitHub

Inicializa Git si aun no existe:

```bash
git init
git add .
git commit -m "feat: crear taskcli con CLI y interfaz web"
git branch -M main
git remote add origin https://github.com/TU_USUARIO/taskcli.git
git push -u origin main
```

Antes de publicar:

- Sustituye `TU_USUARIO` por tu usuario real.
- Revisa que no haya archivos personales dentro de `.taskcli`.
- Decide si quieres publicar el ejecutable generado.
- Añade una licencia si deseas indicar como puede reutilizarse el proyecto.
- Considera agregar un workflow de GitHub Actions.

## Estado del proyecto

Implementado:

- Dominio `Task` y `Comment`.
- CLI.
- Persistencia JSON.
- Interfaz web responsiva.
- Edicion y comentarios.
- Tests HTTP y del dominio.
- Validacion con `go vet` y `go test -race`.

Pendiente como siguiente paso:

- SQLite o PostgreSQL.
- API REST formal.
- Seguridad y autenticacion.
- Concurrencia con goroutines.
- Publicacion automatizada.

## Licencia

Agrega una licencia antes de convertir el repositorio en publico si quieres
indicar claramente los permisos de uso, distribucion y modificacion.
