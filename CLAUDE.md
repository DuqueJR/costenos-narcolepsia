# Proyecto: Hospital de los Costeños con Narcolepsia (Go)

Trabajo individual del curso de POO. Enunciado: `docs/enunciado.pdf`. Plan por fases y ramas: `docs/PLAN.md`.
Antes de cada tarea, lee en `docs/PLAN.md` la sección de la rama que indique el prompt y trabaja SOLO en ese alcance.

## Sobre mí (el estudiante)
- Estoy aprendiendo Go y mañana debo sustentar y modificar en vivo cualquier línea. Prefiere código simple, idiomático y explicable antes que código ingenioso.
- Al terminar cada tarea, explícame en español qué hiciste y por qué (receivers, encapsulación, errores, decisiones).
- Si algo del enunciado o del plan es ambiguo o contradictorio, pregúntame antes de inventar.

## Reglas técnicas (vienen del enunciado, no negociables)
- Go >= 1.21. El módulo principal (raíz) usa SOLO librería estándar: su go.mod nunca gana dependencias.
- La interfaz bonus es un módulo Go aparte en `gui/` (su propio go.mod con `replace costenos-narcolepsia => ../`) y es el ÚNICO lugar donde se permiten terceros: Fyne (`fyne.io/fyne/v2`) y `github.com/gopxl/beep/v2` para el audio. Nada de web (ni HTML/JS, ni APIs REST, ni Wails).
- Fotos y música viven en `gui/media` y en `gui/data/` (gitignored). Nunca en `hospital/`. Los paquetes `gui/media` y `gui/player` no importan Fyne.
- En la Fase 9, `hospital/` y el go.mod raíz NO se modifican. Comprobación: `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío.
- En la GUI, los widgets solo se actualizan en el hilo de la interfaz; desde otra goroutine usa `fyne.Do(...)`. La GUI no ejecuta Simulate.
- Paquete `hospital/` = toda la lógica. `main.go` = solo construir el escenario e imprimir. Nunca `fmt.Print*` dentro de `hospital/`.
- Campos de structs en minúscula (no exportados), acceso por métodos, constructores `NewX`.
- Métodos que pueden fallar retornan `error` como último valor. Nunca `panic` para casos previsibles.
- Errores centinela en `hospital/errors.go` (`var ErrNoRoomAvailable = errors.New(...)`), envueltos con `fmt.Errorf("...: %w", err)`.
- Estados como constantes tipadas con `iota` + método `String()`. Nunca strings sueltos para estados.
- `Person` embebido por valor en `Patient`, `Doctor` y `Orderly`. Ningún campo de identidad duplicado.
- El personal del hospital vive en un `[]Attender` y se despacha sin type switch ni type assertion.
- Firmas exactas de las consultas:
  - `func (h *Hospital) PatientsInHallway() []*Patient`
  - `func (d *Doctor) MyEpisodes() []EpisodeRecord`
  - `func (h *Hospital) AssignRoom(p *Patient) (*Room, error)`
  - `func (h *Hospital) SevereReport() map[*Patient]int`
- Identificadores en inglés. Doc comments en inglés en todo lo exportado, empezando con el nombre del identificador.

## Git: feature branching con rama de integración `dev`
- Estructura: ramas de trabajo (feature/, test/, docs/, chore/, fix/) → PR a `dev` → promoción por PR de `dev` a `main` (sección 3 de docs/PLAN.md).
- NUNCA hagas commits directos en `main` ni en `dev` (única excepción: el commit inicial del Prompt 0.1). Antes de tocar código corre `git branch --show-current`. Si no estás en la rama que indica el prompt, créala desde `dev` actualizado (`git switch dev && git pull && git switch -c <rama>`). Si `git pull` falla porque aún no hay remoto, continúa sin él.
- Los PR de ramas de trabajo siempre llevan base `dev` (`gh pr create --base dev`). Solo `dev` abre PR hacia `main`.
- Commits pequeños y frecuentes con Conventional Commits: `feat:`, `fix:`, `test:`, `docs:`, `chore:`, `refactor:`.
- No hagas merge (ni a `dev` ni a `main`), rebase, `push --force` ni borres ramas a menos que te lo pida explícitamente. Nunca borres `dev`.

## Definición de terminado (córrelo siempre antes de commitear)
(Si la tarea tocó `gui/`, agrega `cd gui && go build ./... && go vet ./...`.)
1. `gofmt -l .` no lista nada (si lista, aplica `gofmt -w .`)
2. `go vet ./...` sin advertencias
3. `go build ./...` compila
4. `go test ./...` pasa

## Comandos del proyecto
- `go run .` → escenario obligatorio
- `go run . -sim` → simulación concurrente (bonus)
- `cd gui && go run .` → interfaz gráfica Fyne (bonus; necesita gcc + CGO; la primera compilación tarda)
- `cd gui && go test ./media/` → tests del almacenamiento de fotos y música
- `go test ./...` y `go test -race ./...`

## Registro de uso de IA
Cuando generes algo incorrecto o no idiomático que luego se corrija (por mí o por ti), agrega una línea a `notas_ia.md` (está en .gitignore): qué estaba mal, por qué y cómo se corrigió. Sirve para escribir `AI_USAGE.md`.
