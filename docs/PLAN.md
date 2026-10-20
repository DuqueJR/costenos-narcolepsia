# Plan de trabajo — Hospital de los Costeños con Narcolepsia (Go)

> **Entrega:** hoy 19/10/2026, 23:59 (link del repositorio) · **Sustentación:** 20/10/2026, individual · **Peso:** 15 %
> **Regla de oro:** la nota se multiplica por tu factor de sustentación (1.0 / 0.8 / 0.5 / 0.0). Todo lo que entregues lo debes poder explicar y modificar en vivo.

---

## 1. Cómo usar este plan con Claude Code

**Estructura de carpetas:**

```
TrabajoGO/
  TrabajoGo.pdf                 ← original, no se toca
  costenos-narcolepsia/         ← EL REPOSITORIO. Claude Code se abre AQUÍ
    CLAUDE.md                   ← reglas permanentes; Claude Code lo lee solo al arrancar
    docs/
      PLAN.md                   ← este archivo; los prompts le piden a Claude leer una sección
      enunciado.pdf             ← copia del enunciado (en .gitignore)
    go.mod, main.go, hospital/  ← lógica + escenario (los crea Claude Code)
    gui/                        ← bonus: interfaz Fyne (módulo Go aparte, con su propio go.mod)
      data/                     ← fotos y música elegidas (en .gitignore)
```

**Cómo se reparte la información:**
- `CLAUDE.md` tiene las reglas que aplican **siempre**: técnicas, de Git y la definición de terminado. Se carga automáticamente en cada sesión, así que no tienes que repetirlas en los prompts.
- `docs/PLAN.md` tiene el detalle de **cada rama**. No se importa entero en CLAUDE.md porque es largo y gastaría contexto; cada prompt le dice a Claude qué sección leer.
- Los prompts de abajo los copias y pegas tal cual en Claude Code.

**Ritual por rama:**
1. Abre la terminal en `costenos-narcolepsia/` y ejecuta `claude`. Si vienes de otra rama, escribe `/clear` para empezar con el contexto limpio.
2. Pega el prompt (o los prompts, en orden) de la rama. En las ramas marcadas con 🧭, activa antes el **modo plan** (Shift+Tab hasta ver "plan mode") para que Claude proponga el plan y tú lo apruebes antes de escribir código.
3. Lee lo que hizo. Si algo no te cuadra, pregúntale ahí mismo.
4. Pega el **prompt de cierre de rama** (sección 3.4): te explica el diff y te hace preguntas de sustentación.
5. Abre el Pull Request **hacia `dev`**, revisa el diff en GitHub y haz merge (sección 3.3). En los checkpoints, promueve `dev` → `main` (sección 3.5).
6. Anota en `notas_ia.md` cualquier cosa que la IA haya hecho mal.

---

## 2. Calendario sugerido (hoy, desde las 11:00)

| Hora objetivo | Fase | Ramas | Rúbrica |
|---|---|---|---|
| 11:00–11:30 | 0. Setup | `main` (commit inicial) + `dev`, `chore/project-setup` | Ingeniería |
| 11:30–12:00 | 1. Base del modelo | `feature/person-and-states` | Entidades |
| 12:00–13:15 | 2. Entidades | `feature/patient-and-room`, `feature/attender-and-episode`, `feature/doctor` | Entidades + Composición |
| 13:15–13:45 | almuerzo | — | — |
| 13:45–15:00 | 3. Orquestación | `feature/orderly`, `feature/hospital-core`, `feature/sleep-attack-flow` | Interfaces |
| 15:00–15:40 | 4. Consultas | `feature/queries` | **Consultas 30 %** |
| 15:40–16:10 | 5. Demo | `feature/demo-scenario` | Demo 15 % |
| 16:10–16:40 | 6. Tests | `test/hospital-tests` | Ingeniería |
| 16:40–17:15 | 7. Calidad | `chore/code-quality`, `docs/readme` | Ingeniería 20 % |
| — | ✅ **Checkpoint: promoción `dev` → `main` (nota base completa)** | | |
| 17:15–18:30 | 8. Bonus concurrencia | `feature/thread-safe-hospital`, `feature/concurrent-simulation` | +5 % |
| 18:30–21:00 | 9. Bonus interfaz Fyne | `feature/gui-shell`, `feature/gui-hospital-views`, `feature/gui-photos`, `feature/gui-music`, `docs/gui-readme` | +10 % |
| 21:00–21:45 | 10. Cierre | `docs/ai-usage`, tag `v1.0.0` | Condición de la nota |
| 21:45–23:59 | colchón | — | — |
| mañana temprano | 11. Sustentación | `practice/*` (nunca se mergean) | Factor ×1.0 |

**Corte duro:** si a las **22:00** una rama de bonus no está lista, **no la mergees a `dev`**. `dev` y `main` se quedan limpios y entregas sin ese bonus. Para eso sirve el feature branching.
En la Fase 9 el orden de las ramas también es la prioridad: con `feature/gui-shell` + `feature/gui-hospital-views` ya cumples el bonus; fotos y música son el plus visual.

---

## 3. Estrategia de versionado: feature branching

### 3.1 Estructura de ramas

```
main  ◄── PR solo desde dev (promociones en los checkpoints)
  │
dev   ◄── PR de cada rama de trabajo
  │
  ├── chore/project-setup
  ├── feature/person-and-states
  ├── feature/...
  └── docs/ai-usage
```

- **`main`** es la versión estable y entregable. **Solo recibe merges desde `dev`**, que son las promociones de la sección 3.5. Nunca recibe commits directos (la única excepción es el commit inicial) ni PR de ramas de trabajo.
- **`dev`** es la rama de integración. Se crea desde `main` justo después del commit inicial y no se borra nunca. **Tampoco recibe commits directos.**
- **Ramas de trabajo** (`feature/`, `test/`, `docs/`, `chore/`, `fix/`): una por unidad de trabajo. Nacen de `dev` actualizado y vuelven a `dev` por **Pull Request**.
- Todos los merges, tanto los PR a `dev` como las promociones a `main`, se hacen con **"Create a merge commit"** (no "Squash"). Así el historial conserva todos los commits y se ven las ramas, que es justo la evidencia que pide el enunciado.
- `dev` siempre compila y pasa tests. `main`, además, corre el escenario completo.
- Una rama de trabajo ya mergeada en `dev` se borra. Nunca se reutiliza una rama vieja.
- Para no equivocarte de base al abrir un PR, `dev` queda como **rama por defecto** del repositorio mientras trabajas. En la entrega final se vuelve a poner `main` (Fase 10).

### 3.2 Prefijos de rama

| Prefijo | Uso | Ejemplo |
|---|---|---|
| `main` | versión entregable; solo recibe PR desde `dev` | — |
| `dev` | integración; recibe los PR de todas las ramas de trabajo | — |
| `feature/` | funcionalidad nueva | `feature/doctor` |
| `test/` | solo tests | `test/hospital-tests` |
| `docs/` | solo documentación | `docs/readme` |
| `chore/` | setup, formato, mantenimiento | `chore/code-quality` |
| `fix/` | corrección de algo ya mergeado | `fix/hospital-getters` |
| `practice/` | práctica de sustentación, **nunca se mergea** | `practice/nurse` |

Los commits siguen Conventional Commits: `feat: add Room.Occupy`, `test: cover no-room case`, `docs: write README`.

### 3.3 Flujo de cada rama (comandos)

```bash
git switch dev
git pull
git switch -c feature/<nombre>
# ... prompts en Claude Code (él hace los commits en la rama) ...
gofmt -l . ; go vet ./... ; go test ./...
git push -u origin feature/<nombre>
# En GitHub: Compare & pull request → verifica que la base sea "dev" → revisar "Files changed" → "Create a merge commit" → Delete branch
#   (alternativa CLI: gh pr create --base dev --fill   y luego   gh pr merge --merge --delete-branch)
git switch dev
git pull
git branch -d feature/<nombre>
```

Si una rama necesita un cambio de una rama anterior que aún no está en `dev`, **termina y mergea primero la anterior**. No ramifiques desde ramas sin mergear: hoy no hay tiempo para resolver conflictos.

### 3.4 Prompts que se repiten

**Prompt de cierre de rama** (pégalo al terminar cada rama, antes del PR):
```
Terminamos esta rama. Haz lo siguiente:
1. Corre la definición de terminado de CLAUDE.md y confirma que todo pasa.
2. Muéstrame `git log --oneline dev..HEAD` y `git diff --stat dev...HEAD`.
3. Explícame en español, archivo por archivo, lo que cambió en esta rama, como si yo tuviera que sustentarlo: qué hace cada tipo/método, por qué ese receiver (puntero o valor), por qué ese tipo de retorno y qué errores puede devolver.
4. Hazme 3 preguntas tipo sustentación sobre esta rama, una a la vez, y espera mi respuesta antes de corregirme.
No hagas push ni merge.
```

**Prompt para abrir el PR** (opcional, si tienes `gh` instalado y autenticado):
```
Haz push de la rama actual a origin y crea un Pull Request con base dev (`gh pr create --base dev`). Nunca contra main. Título en formato Conventional Commits; en el cuerpo, un resumen en viñetas de lo que agrega la rama y cómo se verificó (gofmt, vet, tests). NO hagas merge: yo reviso el diff y lo mergeo.
```

### 3.5 Promoción `dev` → `main`

**Cuándo** (estos son los checkpoints):
1. Después de la Fase 5, cuando el escenario ya corre. Así el link del repositorio ya muestra algo entregable en `main`.
2. Después de la Fase 7: nota base completa.
3. Después de `docs/ai-usage`.
4. Después de cada bonus mergeado en `dev`.
5. Promoción final, antes de entregar (Fase 10).

**Flujo:**
```bash
git switch dev
git pull
gofmt -l . ; go vet ./... ; go test ./... ; go run .
gh pr create --base main --head dev --title "release: <qué incluye>" --fill
gh pr merge --merge          # SIN --delete-branch: dev no se borra
git switch main
git pull
git switch dev               # se sigue trabajando desde dev
```
En la web es lo mismo: **New pull request** → base `main`, compare `dev` → **Create a merge commit**. Si GitHub ofrece borrar la rama, **no borres `dev`**.

Después de una promoción, GitHub puede mostrar `dev` como "1 commit behind main" por el merge commit. Es normal y no causa conflictos, porque `main` solo recibe cambios desde `dev`.

**Prompt de promoción:**
```
Vamos a promover dev a main. NO modifiques código.
1. `git switch dev && git pull`. Corre la definición de terminado de CLAUDE.md y `go run .`; si algo falla, detente y dímelo.
2. Muéstrame `git log --oneline main..dev` (lo que entra en esta promoción) y resume qué ramas incluye.
3. Crea el PR con `gh pr create --base main --head dev`, con título "release: <resumen>" y en el cuerpo la lista de ramas incluidas y cómo se verificó.
4. NO hagas merge: yo lo reviso y lo mergeo con "Create a merge commit", sin borrar dev.
```

---

## 4. Decisiones de diseño (las sigue todo el plan)

Apréndetelas: son tus justificaciones en la sustentación.

1. **Nombres en inglés**, consistentes con el enunciado.
2. **El `Hospital` orquesta.** `HandleSleepAttack(p, location)` ejecuta: `p.SufferSleepAttack(location)` → `h.AssignRoom(p)` (si falla, se guarda el error) → `h.dispatch()` elige un `Attender` → `attender.Attend(p, location)` → `h.RegisterEpisode(rec)` → retorna `(rec, roomErr)`. `main` solo llama métodos.
3. **El paciente guarda `room *Room`.** El `EpisodeRecord` toma ese cuarto al crearse (nil = se quedó en el pasillo), y al despertar se sabe qué cuarto liberar.
4. **`Room.Occupy(p)` cambia al paciente mediante métodos no exportados** (`p.moveToBed(r)`), no tocando sus campos directamente, aunque estén en el mismo paquete.
5. **En `EpisodeRecord`, "quién atendió" es de tipo `Attender`**, no `*Doctor`, porque el camillero también atiende.
6. **`Doctor.Attend` agrega el registro a `d.episodes`.** Por eso `MyEpisodes()` funciona sin recibir el hospital.
7. **El ID del episodio lo genera quien atiende**, con su propio contador (`D-01-E1`), sin variables globales (que además romperían la concurrencia).
8. **Despacho round-robin** sobre el `[]Attender`: simple, determinista y fácil de defender.
9. **Errores centinela** en `errors.go`; los tests los verifican con `errors.Is`.
10. **`SufferSleepAttack` retorna `error`** (si el paciente ya está dormido). El enunciado da las firmas mínimas y pide que lo que pueda fallar retorne `error`.
11. **`SevereReport` incluye a los Severe con 0 episodios.** El map no tiene orden, así que `main` ordena por ID antes de imprimir.
12. **Getters de lectura en `Hospital` desde la Fase 3** (`Patients()`, `Rooms()`, `FindPatient(id)`...), para que el bonus de interfaz no obligue a tocar `hospital/`.
13. **Cuando el paciente del pasillo recibe cuarto más tarde, no se crea un episodio nuevo.** El episodio original queda como ocurrió.
14. **Se respeta el escenario exacto.** El reporte de severidad dará conteos bajos y está bien; el ejemplo del PDF es ilustrativo.
15. **Fotos y música son responsabilidad de la vista, no del dominio.** Viven en `gui/media` y en disco (`gui/data/`), asociadas al paciente por su ID. El paquete `hospital/` no sabe que existen.
16. **La interfaz Fyne es un módulo Go aparte** (`gui/go.mod` con `replace costenos-narcolepsia => ../`). El módulo principal queda 100 % librería estándar, y `go run .` / `go test ./...` en la raíz no necesitan gcc ni descargan Fyne. Es la prueba más fuerte de la separación lógica/vista: la GUI solo puede usar lo que `hospital` exporta.

---

## FASE 0 — Setup del proyecto

### Rama `main` (commit inicial) — Prompt 0.1

**Para qué:** no se puede ramificar desde un repositorio vacío. Este es el único commit directo en `main`, y deja en el historial el plan y las reglas desde el minuto cero. Aquí también nace `dev`, la rama de integración.

**Qué va:** `git init`, `.gitignore`, un README mínimo, `CLAUDE.md` y `docs/PLAN.md`. Repositorio en GitHub.

**Alcance:** nada de código Go.

**Terminado cuando:** el repositorio existe en GitHub con un commit en `main` y la rama `dev` creada desde él, y `dev` es la rama por defecto.

```
Esta es la única tarea que se hace directamente en main.
1. Inicializa git con rama principal main (`git init -b main`).
2. Crea un .gitignore para Go que incluya: binarios (*.exe, *.out, /costenos-narcolepsia), *.test, coverage.out, .vscode/, .idea/, notas_ia.md, docs/enunciado.pdf y, para la interfaz bonus, gui/data/.
3. Crea un README.md mínimo con el título del proyecto, mi nombre (Jerónimo) como autor y la línea "En construcción".
4. Crea notas_ia.md vacío con un encabezado (queda ignorado por git).
5. Haz un único commit: "chore: initial commit with plan and project rules", que incluya CLAUDE.md, docs/PLAN.md, README.md y .gitignore. Verifica con `git status` que enunciado.pdf y notas_ia.md NO queden incluidos.
6. Si tengo `gh` instalado y autenticado, crea el repositorio público "costenos-narcolepsia" en mi GitHub y haz push de main (`gh repo create costenos-narcolepsia --public --source=. --push`). Si no, dime los comandos exactos para hacerlo desde la web.
7. Crea la rama de integración dev desde main (`git switch -c dev`) y súbela (`git push -u origin dev`). Si tengo gh, déjala como rama por defecto del repositorio (`gh repo edit --default-branch dev`) para que los PR apunten a dev por defecto; si no, dime cómo hacerlo en Settings → General → Default branch.
Explícame qué hace cada línea del .gitignore y el rol de main, dev y las ramas de trabajo.
```

### Rama `chore/project-setup` — Prompt 0.2

**Para qué:** dejar el módulo Go y la estructura listos, y detectar ya el problema de Windows con `-race`, que puede costarte el bonus de concurrencia.

**Qué va:** `go.mod`, la estructura de carpetas y un `main.go` mínimo. Verificación de la versión de Go y de gcc/CGO, que necesitan tanto `-race` como Fyne (la interfaz bonus).

**Alcance:** solo estructura y herramientas, sin lógica.

**Terminado cuando:** `go build ./...` compila y gcc funciona (lo necesitan `go run -race` y Fyne).

```
Rama: chore/project-setup (créala desde dev según CLAUDE.md). Lee docs/PLAN.md, sección "FASE 0".
1. Corre `go version` y dime la versión: debe ser >= 1.21. Avísame si es muy antigua, porque la versión actual de Fyne (bonus de interfaz) puede pedir una más reciente.
2. `go mod init costenos-narcolepsia`.
3. Crea la estructura con archivos que compilen (solo `package hospital` y un comentario de paquete en cada uno): hospital/doc.go (comentario de paquete), errors.go, person.go, patient.go, room.go, episode.go, attender.go, doctor.go, orderly.go, hospital.go, hospital_test.go. Y un main.go con `package main` y un `func main()` que imprima solo el nombre del hospital.
4. Verifica el race detector: corre `go env CGO_ENABLED` y `gcc --version`, y después intenta `go run -race .`. Si falla, explícame por qué (en Windows el race detector necesita CGO y un compilador C) y dame las opciones concretas: instalar gcc con WinLibs/MSYS2 y agregarlo al PATH, o usar WSL. No instales nada sin preguntarme.
5. gcc también es requisito de Fyne (la interfaz bonus): sin gcc no hay ninguno de los dos bonus. Si no está, dame instrucciones paso a paso para instalar MSYS2 con el toolchain mingw-w64 (o TDM-GCC), agregarlo al PATH y verificar con `gcc --version` y `go env CGO_ENABLED` (debe dar 1). No lo instales sin preguntarme.
6. Corre la definición de terminado y haz commit: "chore: go module and package skeleton".
Explícame qué es go.mod, qué es un paquete y por qué hospital y main van separados.
```

---

## FASE 1 — Base del modelo

### Rama `feature/person-and-states` — Prompts 1.1 y 1.2

**Para qué:** los cimientos que usan todas las entidades. Cubre "Entidades" (constantes tipadas, encapsulación) y prepara la composición.

**Qué va:** `Person` con campos privados y getters, y los tres tipos de estado con `iota` y `String()`.

**Alcance:** solo tipos base. Nada de Patient, Doctor ni Hospital.

**Terminado cuando:** ningún estado es un string suelto y `fmt.Println(AsleepInHallway)` imprime texto legible.

**Debes entender:** exportado vs no exportado (el **compilador** impide `p.name` desde otro paquete), `iota`, la interfaz `fmt.Stringer` y por qué las constantes tipadas evitan errores de tipeo.

**Prompt 1.1 — Person**
```
Rama: feature/person-and-states (créala desde dev). Lee docs/PLAN.md, sección "FASE 1".
Implementa en hospital/person.go:
- `type Person struct { id string; name string; age int }` con campos NO exportados.
- Un constructor no exportado `newPerson(id, name string, age int) (Person, error)` que valide id y name no vacíos (strings.TrimSpace) y age > 0. Si los datos son inválidos, retorna un error que envuelva `ErrInvalidData`.
- Getters con receiver de VALOR: `ID() string`, `Name() string`, `Age() int`.
- En hospital/errors.go, declara `ErrInvalidData` con errors.New.
Restricciones: no exportes campos ni el constructor (Person solo se crea a través de NewPatient/NewDoctor/NewOrderly).
Explícame por qué aquí el receiver de valor está bien y qué pasaría si Person tuviera un setter con receiver de valor.
```

**Prompt 1.2 — Estados tipados**
```
Seguimos en feature/person-and-states. En hospital/person.go (o en un hospital/states.go, como te parezca más claro), implementa:
- `type PatientState int` con const Awake, AsleepInHallway, AsleepInBed usando iota.
- `type NarcolepsyLevel int` con Mild, Moderate, Severe.
- `type RoomState int` con Available, Occupied.
- Un método `String() string` para cada tipo, con switch y un default que devuelva algo como "PatientState(7)" para valores inválidos.
- En hospital_test.go, un test table-driven pequeño que verifique los String() de todos los valores.
Corre la definición de terminado y haz commits separados para Person y para los estados.
Explícame cómo funciona iota y por qué fmt.Println llama a String() automáticamente.
```

---

## FASE 2 — Entidades

### Rama `feature/patient-and-room` — Prompts 2.1a y 2.1b

**Para qué:** `Patient` y `Room` dependen uno del otro (el cuarto tiene ocupantes y el paciente sabe su cuarto), por eso van en la misma rama. Cubre "Entidades" y la primera parte de "Composición".

**Qué va:** `Patient` con `Person` embebido, `Room`, constructores y métodos con errores.

**Alcance:** **no** incluye el campo `assignedDoctor` (se agrega en `feature/doctor`), ni el Hospital.

**Terminado cuando:** puedes crear pacientes y cuartos, ocupar, liberar y despertar con errores correctos, y hay tests básicos.

**Debes entender:** métodos promovidos por embedding (no es herencia), receiver de puntero vs de valor, por qué existen los constructores y por qué `Occupy` usa un método no exportado del paciente.

**Prompt 2.1a — Patient**
```
Rama: feature/patient-and-room (créala desde dev). Lee docs/PLAN.md, secciones "Decisiones de diseño" y "FASE 2".
Implementa hospital/patient.go:
- `type Patient struct` con `Person` embebido POR VALOR y los campos no exportados: level NarcolepsyLevel, state PatientState, currentLocation string, room *Room. NO agregues assignedDoctor todavía.
- `NewPatient(id, name string, age int, level NarcolepsyLevel) (*Patient, error)`: usa newPerson, valida que level sea válido y empieza en estado Awake.
- Métodos con receiver PUNTERO:
  - `SufferSleepAttack(location string) error`: si ya está dormido, retorna `ErrAlreadyAsleep`; si no, state = AsleepInHallway, currentLocation = location, room = nil.
  - `WakeUp() error`: si ya está despierto, retorna `ErrAlreadyAwake`; si no, state = Awake y room = nil.
  - Getters: `State() PatientState`, `Level() NarcolepsyLevel`, `Location() string`, `Room() *Room`.
  - No exportados: `moveToBed(r *Room)` (state = AsleepInBed, room = r).
- Agrega los errores nuevos a errors.go.
- Para que compile, crea en room.go solo `type Room struct` con sus campos (number int, capacity int, state RoomState, occupants []*Patient). Los métodos de Room van en el prompt siguiente.
Explícame qué métodos tiene Patient gracias al embedding y por qué todos sus métodos usan receiver puntero.
```

**Prompt 2.1b — Room + tests**
```
Seguimos en feature/patient-and-room. Completa hospital/room.go:
- `NewRoom(number, capacity int) (*Room, error)`: number > 0 y capacity > 0, si no `ErrInvalidData`. Estado inicial Available.
- `IsAvailable() bool`: true si len(occupants) < capacity.
- `Occupy(p *Patient) error`: nil → ErrInvalidData; lleno → `ErrRoomFull`; paciente ya adentro → `ErrAlreadyInRoom`. Si todo va bien: agrega a occupants, llama p.moveToBed(r) y, si el cuarto quedó lleno, state = Occupied.
- `Release(p *Patient) error`: si el paciente no está → `ErrNotInRoom`; si está, lo quita de occupants y state = Available.
- Getters: `Number() int`, `Capacity() int`, `State() RoomState`, `Occupants() []*Patient` (retorna una COPIA del slice).
- Tests en hospital_test.go: Occupy exitoso (cuarto Occupied, paciente AsleepInBed), Occupy sobre cuarto lleno (errors.Is ErrRoomFull), Release de un paciente que no está (ErrNotInRoom), y WakeUp de un paciente ya despierto (ErrAlreadyAwake).
Corre la definición de terminado y haz commits separados (patient, room, tests).
Explícame por qué Occupants() retorna una copia y qué pasaría si retornara el slice interno.
```

### Rama `feature/attender-and-episode` — Prompt 2.2

**Para qué:** el `EpisodeRecord` necesita saber "quién atendió", que es un `Attender`. Por eso se declara la interfaz antes de implementarla.

**Qué va:** la interfaz `Attender` y `EpisodeRecord` con `Summary()`.

**Alcance:** solo la declaración de la interfaz. Doctor y Orderly la implementan en ramas posteriores.

**Terminado cuando:** compila y `Summary()` tiene un test.

**Debes entender:** las interfaces en Go se cumplen **implícitamente**, y por qué el campo es `Attender` y no `*Doctor`.

```
Rama: feature/attender-and-episode (créala desde dev). Lee docs/PLAN.md, sección "FASE 2".
1. En hospital/attender.go declara:
   type Attender interface {
       ID() string
       Name() string
       Attend(p *Patient, location string) (EpisodeRecord, error)
   }
   con un doc comment que explique que lo implementan Doctor y Orderly y que el hospital despacha a través de él.
2. En hospital/episode.go implementa `EpisodeRecord` con campos no exportados: id string, at time.Time, patient *Patient, attendedBy Attender, location string, room *Room (nil = pasillo).
   - Constructor NO exportado `newEpisodeRecord(id string, at time.Time, p *Patient, by Attender, location string, room *Room) EpisodeRecord` (solo los Attender crean registros).
   - Getters: ID(), At(), Patient(), AttendedBy(), Location(), Room().
   - `Summary() string` con un formato como: "2026-10-19 14:15 P-001 cafeteria -> room 101 (Dr. Karen Ospina)" o "-> hallway" si room es nil. Usa el layout "2006-01-02 15:04".
3. Test de Summary() para los dos casos (con cuarto y sin cuarto). Para el Attender del test puedes crear un fake mínimo dentro de hospital_test.go.
Definición de terminado y commit.
Explícame qué significa que la interfaz se cumpla implícitamente y por qué el fake del test ya es un Attender sin declararlo.
```

### Rama `feature/doctor` — Prompt 2.3

**Para qué:** el primer `Attender` real, que modela la relación doctor–episodios dentro de la entidad (consulta 5.2).

**Qué va:** `Doctor` con `Person` embebido, `DiagnosePatient`, `AttendSleepEmergency`, `Attend` y `MyEpisodes`. Se agrega `assignedDoctor` a `Patient`.

**Alcance:** nada de Hospital ni de Orderly.

**Terminado cuando:** `var _ Attender = (*Doctor)(nil)` compila y `MyEpisodes` tiene un test.

**Debes entender:** cómo `Doctor` cumple `ID()`/`Name()` gracias a `Person`, por qué solo `*Doctor` (y no `Doctor`) cumple `Attender`, y por qué `MyEpisodes` no necesita el hospital.

```
Rama: feature/doctor (créala desde dev). Lee docs/PLAN.md, "Decisiones de diseño" (puntos 5, 6 y 7) y "FASE 2".
1. En patient.go agrega el campo `assignedDoctor *Doctor`, el getter `AssignedDoctor() *Doctor` y el método no exportado `assignDoctor(d *Doctor)`.
2. En hospital/doctor.go implementa `Doctor` con Person embebido por valor y los campos no exportados: specialty string, patients []*Patient, episodes []EpisodeRecord.
   - `NewDoctor(id, name string, age int, specialty string) (*Doctor, error)`.
   - `DiagnosePatient(p *Patient) error`: nil → ErrInvalidData; si ya está a su cargo → `ErrAlreadyUnderCare`; si no, lo agrega a patients y llama p.assignDoctor(d).
   - `AttendSleepEmergency(p *Patient) error`: si el paciente está Awake → `ErrPatientAwake`; si no tiene doctor asignado, lo diagnostica.
   - `Attend(p *Patient, location string) (EpisodeRecord, error)`: llama AttendSleepEmergency; crea el registro con ID `fmt.Sprintf("%s-E%d", d.ID(), len(d.episodes)+1)`, time.Now() y p.Room(); lo agrega a d.episodes y lo retorna.
   - `MyEpisodes() []EpisodeRecord` y `Patients() []*Patient` retornan COPIAS. `Specialty() string`.
   - Debajo del tipo agrega la verificación en compilación `var _ Attender = (*Doctor)(nil)`.
3. Tests: Attend sobre un paciente dormido agrega exactamente un registro a MyEpisodes; Attend sobre uno despierto retorna ErrPatientAwake; DiagnosePatient dos veces retorna ErrAlreadyUnderCare.
Definición de terminado y commits.
Explícame: (a) qué hace `var _ Attender = (*Doctor)(nil)`, (b) qué pasaría con MyEpisodes si Attend tuviera receiver de valor, (c) qué cambiaría si Person se embebiera como *Person.
```

---

## FASE 3 — Orquestación

### Rama `feature/orderly` — Prompt 3.1

**Para qué:** sin un segundo tipo que cumpla `Attender` no hay polimorfismo, y ese ítem vale 0.

**Qué va:** `Orderly` (camillero) con `Person` embebido. Implementa `Attend` pero **no** puede diagnosticar.

**Alcance:** solo `Orderly`.

**Terminado cuando:** `var _ Attender = (*Orderly)(nil)` compila y hay un test que usa un `[]Attender` con Doctor y Orderly.

**Debes entender:** esta rama es casi idéntica a lo que te pedirán en vivo con el `Nurse`. **Estúdiala.**

```
Rama: feature/orderly (créala desde dev). Lee docs/PLAN.md, "FASE 3".
En hospital/orderly.go implementa `Orderly` (camillero): Person embebido por valor y un contador no exportado `attended int`.
- `NewOrderly(id, name string, age int) (*Orderly, error)`.
- `Attend(p *Patient, location string) (EpisodeRecord, error)`: si el paciente está Awake → ErrPatientAwake; si no, crea el registro con ID `fmt.Sprintf("%s-E%d", o.ID(), o.attended+1)`, incrementa el contador y lo retorna. NO diagnostica ni asigna doctor (doc comment explicando que un camillero traslada pero no diagnostica).
- `Attended() int`.
- `var _ Attender = (*Orderly)(nil)`.
Test: crea `staff := []Attender{doctor, orderly}`, duerme a dos pacientes, recorre staff llamando Attend sin type switch y verifica que ambos registros tengan el AttendedBy correcto.
Definición de terminado y commit.
Después, explícame paso a paso cómo agregaría yo un tipo Nurse que cumpla Attender, sin escribirlo.
```

### Rama `feature/hospital-core` — Prompt 3.2

**Para qué:** el contenedor del sistema, con registro y búsqueda de todo. También deja listos los getters que necesita la interfaz Fyne.

**Qué va:** `Hospital`, `NewHospital`, `AdmitPatient`, `HireDoctor`, `HireStaff`, `AddRoom`, getters y búsquedas.

**Alcance:** **no** incluye todavía el flujo del ataque de sueño ni las consultas.

**Terminado cuando:** puedes armar un hospital completo y los duplicados devuelven error.

```
Rama: feature/hospital-core (créala desde dev). Lee docs/PLAN.md, "Decisiones de diseño" (punto 12) y "FASE 3".
En hospital/hospital.go implementa `Hospital` con campos no exportados: name string, patients []*Patient, doctors []*Doctor, staff []Attender, rooms []*Room, history []EpisodeRecord, next int (índice del round-robin).
- `NewHospital(name string) *Hospital`.
- `AdmitPatient(p *Patient) error`: nil → ErrInvalidData; ID repetido → `ErrDuplicateID` (envuelto con el ID en el mensaje).
- `HireStaff(a Attender) error`: nil → ErrInvalidData; ID repetido entre todo el staff → ErrDuplicateID. Agrega al slice staff. Sirve para camilleros y para cualquier Attender futuro (un Nurse).
- `HireDoctor(d *Doctor) error`: usa HireStaff y además agrega a doctors.
- `AddRoom(r *Room) error`: nil o número repetido → error.
- Getters que retornan COPIAS: Name(), Patients(), Doctors(), Staff(), Rooms(), History().
- Búsquedas: `FindPatient(id string) (*Patient, error)`, `FindDoctor(id string) (*Doctor, error)`, `FindRoom(number int) (*Room, error)`, con `ErrNotFound` envuelto.
Tests: duplicados en AdmitPatient, HireStaff y AddRoom; FindPatient inexistente retorna ErrNotFound.
Definición de terminado y commits.
Explícame por qué HireDoctor guarda al doctor en dos slices y qué gano teniendo HireStaff(a Attender).
```

### Rama `feature/sleep-attack-flow` 🧭 — Prompts 3.3a y 3.3b

**Para qué:** es el corazón del sistema: dónde el hospital usa el `[]Attender` de forma polimórfica y dónde se maneja el caso "no hay cuarto".

**Qué va:** `AssignRoom`, `dispatch`, `RegisterEpisode`, `HandleSleepAttack` y `WakeUpPatient`.

**Alcance:** la consulta `AssignRoom` se implementa aquí; en la Fase 4 se revisa contra el enunciado.

**Terminado cuando:** el flujo completo funciona y, sin cuarto, el paciente queda en el pasillo con un error legible.

**Prompt 3.3a — Asignación, despacho y registro** (en modo plan)
```
Rama: feature/sleep-attack-flow (créala desde dev). Lee docs/PLAN.md, "Decisiones de diseño" (puntos 2, 3, 8 y 13) y "FASE 3". Primero propón el plan y espera mi aprobación.
En hospital.go implementa:
- `AssignRoom(p *Patient) (*Room, error)`: nil → ErrInvalidData; si el paciente no está AsleepInHallway → `ErrNotInHallway`; recorre h.rooms en orden y ocupa el PRIMERO con IsAvailable() usando r.Occupy(p). Si ninguno está libre, retorna `nil, fmt.Errorf("assign room to %s: %w", p.ID(), ErrNoRoomAvailable)` SIN cambiar el estado del paciente (sigue AsleepInHallway). Nunca crea cuartos ni hace panic.
- `dispatch() (Attender, error)` (no exportado): si staff está vacío → `ErrNoStaff`; si no, round-robin con h.next. Sin type switch.
- `RegisterEpisode(rec EpisodeRecord) error`: valida que el registro tenga paciente y lo agrega a history.
Tests: AssignRoom con cuarto libre (retorna el cuarto, cuarto Occupied, paciente AsleepInBed); AssignRoom sin cuartos libres (errors.Is ErrNoRoomAvailable, paciente sigue AsleepInHallway, room nil); dispatch alterna entre los Attender.
Definición de terminado y commits.
```

**Prompt 3.3b — Flujo completo**
```
Seguimos en feature/sleep-attack-flow.
- `HandleSleepAttack(p *Patient, location string) (EpisodeRecord, error)`: p.SufferSleepAttack(location) (si falla, retorna ese error); roomErr := h.AssignRoom(p) (NO aborta si es ErrNoRoomAvailable, pero sí ante cualquier otro error); attender := h.dispatch(); rec := attender.Attend(p, location); h.RegisterEpisode(rec); retorna (rec, roomErr). Doc comment que explique claramente que un error ErrNoRoomAvailable significa "episodio registrado, pero el paciente quedó en el pasillo".
- `WakeUpPatient(p *Patient) error`: si p.Room() != nil, primero r.Release(p); luego p.WakeUp().
Tests: 3 cuartos y 4 ataques → los 3 primeros sin error y el cuarto con ErrNoRoomAvailable, 4 registros en History(); WakeUpPatient libera el cuarto y luego AssignRoom del paciente del pasillo obtiene ese cuarto.
Definición de terminado y commits.
Explícame en qué línea exacta ocurre el polimorfismo y por qué HandleSleepAttack retorna el registro Y el error al mismo tiempo.
```

---

## FASE 4 — Las 4 consultas (30 %)

### Rama `feature/queries` — Prompts 4.1 y 4.2

**Para qué:** el ítem que más pesa. Cada consulta **retorna datos, no imprime**, con la firma exacta del enunciado.

**Qué va:** `PatientsInHallway` y `SevereReport` nuevas, más una revisión de `AssignRoom` y `MyEpisodes` contra el enunciado.

**Alcance:** nada de `fmt.Print` en `hospital/`.

**Terminado cuando:** las 4 firmas son exactas, están testeadas y hay cero prints en el paquete.

**Debes entender:** por qué `map[*Patient]int` (relación directa paciente→conteo, la clave identifica al paciente único), por qué el map no tiene orden, y cómo funcionan `%w` y `errors.Is`.

**Prompt 4.1 — Consultas nuevas**
```
Rama: feature/queries (créala desde dev). Lee docs/PLAN.md, "Decisiones de diseño" (punto 11) y "FASE 4". Revisa también la sección 5 de docs/enunciado.pdf.
Implementa en hospital.go:
- `func (h *Hospital) PatientsInHallway() []*Patient`: recorre h.patients y retorna los que están en AsleepInHallway, en orden de admisión. Si no hay, retorna un slice vacío (no nil).
- `func (h *Hospital) SevereReport() map[*Patient]int`: incluye a TODOS los pacientes Severe (con 0 si no tuvieron episodios) y cuenta sus episodios del día actual en h.history. Usa una función auxiliar no exportada `sameDay(a, b time.Time) bool` que compare año, mes y día.
Tests: PatientsInHallway después de un fallo de cuarto contiene exactamente al paciente esperado; SevereReport con dos pacientes Severe (uno con 2 episodios y otro con 0) y uno Mild que no debe aparecer.
Definición de terminado y commits.
```

**Prompt 4.2 — Auditoría contra el enunciado**
```
Seguimos en feature/queries. Audita las 4 consultas contra la sección 5 de docs/enunciado.pdf y contra CLAUDE.md:
1. Las firmas deben ser EXACTAMENTE las del enunciado. Muéstrame cada firma actual al lado de la exigida.
2. Ninguna consulta imprime. Busca cualquier fmt.Print en hospital/ y repórtalo.
3. MyEpisodes no recibe el hospital.
4. AssignRoom sin cuarto: retorna error, no crashea, no inventa cuartos y el paciente sigue AsleepInHallway.
Si algo no cumple, corrígelo con un commit `fix:`. Después explícame por qué SevereReport retorna un map y no un slice, y qué implica eso al imprimir.
```

---

## FASE 5 — Escenario obligatorio

### Rama `feature/demo-scenario` — Prompts 5.1 y 5.2

**Para qué:** el evaluador compara las salidas de todos con este escenario. Cubre "Demo" (15 %).

**Qué va:** `main.go` con el escenario exacto y la impresión de las 4 consultas.

**Alcance:** solo construcción e impresión. Si aparece lógica de negocio en `main`, va a `hospital/`.

**Terminado cuando:** `go run .` muestra el escenario completo, con el error legible y las 4 consultas identificables.

**Prompt 5.1 — Escenario**
```
Rama: feature/demo-scenario (créala desde dev). Lee docs/PLAN.md "FASE 5" y la sección 6 de docs/enunciado.pdf.
Escribe main.go con este escenario, en este orden exacto:
1. Crear el hospital "Hospital de los Costeños con Narcolepsia". Contratar 2 doctores (D-01, D-02, con especialidad) y 1 camillero (O-01) con HireDoctor/HireStaff. Crear 3 cuartos (101, 102, 103) de capacidad 1. Admitir 5 pacientes P-001…P-005, con P-002 y P-004 en Severe. Usa nombres como los del ejemplo del enunciado.
2. Los 4 primeros pacientes sufren un ataque en lugares distintos (cafetería, sala de espera de radiología, parqueadero, pasillo 2 segundo piso) usando h.HandleSleepAttack. Imprime cada evento con rec.Summary(). Al 4.º le corresponde ErrNoRoomAvailable: detéctalo con errors.Is e imprime un mensaje legible, sin terminar el programa.
3. Imprime "== Patients asleep in hallway ==" justo después del fallo (debe aparecer P-004).
4. P-002 despierta (h.WakeUpPatient), se imprime qué cuarto liberó y luego h.AssignRoom(P-004), imprimiendo el cuarto asignado.
5. Al final imprime las 4 consultas con encabezados claros: hallway (ahora vacío: imprime "(none)"), episodios de CADA doctor con MyEpisodes, un bloque "== Room assignment ==" con el resultado del fallo y de la reasignación, y "== Severity report (Severe patients) ==".
Cualquier error inesperado (no ErrNoRoomAvailable) se imprime y termina con os.Exit(1). main.go NO tiene lógica de negocio: solo llamadas a métodos e impresión. Si necesitas lógica, dime qué método agregar a hospital/.
Definición de terminado, corre `go run .`, muéstrame la salida completa y haz commit.
```

**Prompt 5.2 — Pulido de la salida**
```
Seguimos en feature/demo-scenario.
1. Las claves de SevereReport() vienen en orden aleatorio: ordénalas por ID (sort.Slice o slices.SortFunc) antes de imprimir y explícame por qué es necesario. Corre `go run .` 3 veces y confirma que la salida es estable (salvo la hora).
2. Alinea las columnas con fmt.Printf (ancho fijo para ID y nombre), como en el ejemplo del enunciado.
3. Separa la impresión en funciones auxiliares de main.go (printHallway, printDoctorEpisodes, printSevereReport) para que main() se lea como el guion del escenario.
4. Revisa que main.go no tenga lógica de negocio y explícame la diferencia entre "imprimir" y "decidir".
Definición de terminado y commit.
```

> ✅ **Checkpoint:** después de mergear `feature/demo-scenario` en `dev`, haz la **promoción `dev` → `main`** (sección 3.5). Desde aquí el link del repositorio ya es entregable.

---

## FASE 6 — Tests

### Rama `test/hospital-tests` — Prompts 6.1 y 6.2

**Para qué:** el enunciado exige mínimo 3 tests (asignación con cuarto, asignación sin cuarto y una consulta). Ya tienes varios de las ramas anteriores; aquí se completa y ordena la suite.

**Qué va:** revisión de cobertura, tests faltantes en estilo table-driven y la cobertura del paquete.

**Alcance:** solo `hospital_test.go` (y, si se vuelve largo, otros `*_test.go` del paquete).

**Terminado cuando:** `go test ./... -cover` pasa y la cobertura es razonable (>70 %).

**Prompt 6.1 — Auditoría y obligatorios**
```
Rama: test/hospital-tests (créala desde dev). Lee docs/PLAN.md "FASE 6".
1. Lista los tests que ya existen y dime cuáles cubren los 3 obligatorios del enunciado (asignación con cuarto libre, fallo sin cuarto, una consulta). Si alguno falta o es débil, refuérzalo con nombres claros: TestAssignRoom_FreeRoom, TestAssignRoom_NoRoomAvailable, TestPatientsInHallway.
2. Agrega un helper no exportado en el test (por ejemplo newTestHospital(t, rooms int)) que arme un hospital con staff, cuartos y pacientes, usando t.Helper().
3. Corre `go test ./... -v -cover` y muéstrame el resultado.
Commit: "test: required room assignment and query tests".
```

**Prompt 6.2 — Tests adicionales**
```
Seguimos en test/hospital-tests. Agrega tests table-driven para:
- SevereReport: conteos correctos, incluye Severe con 0 episodios y excluye Mild/Moderate.
- MyEpisodes: cada doctor ve solo sus propios episodios después de varios HandleSleepAttack.
- Polimorfismo: después de N ataques, History() tiene registros atendidos tanto por *Doctor como por *Orderly (verifícalo con AttendedBy().ID(), sin type switch).
- Errores: WakeUp sobre paciente despierto, SufferSleepAttack sobre paciente dormido, AdmitPatient duplicado (todos con errors.Is).
Corre `go test ./... -cover` y dime la cobertura. Commit: "test: severity report, episodes and error cases".
Explícame qué es un table-driven test, la diferencia entre t.Errorf y t.Fatalf, y por qué el test está en `package hospital` y no en `package hospital_test`.
```

---

## FASE 7 — Calidad y README

### Rama `chore/code-quality` — Prompt 7.1

**Para qué:** el 20 % de "Idiomático Go e ingeniería".

**Qué va:** gofmt, vet, doc comments y una revisión de idiomas de Go.

**Alcance:** no se agregan features.

```
Rama: chore/code-quality (créala desde dev). Lee docs/PLAN.md "FASE 7".
1. `gofmt -l .` y `go vet ./...`: todo limpio.
2. Doc comments en inglés en TODOS los identificadores exportados, empezando con el nombre (`// AssignRoom finds ...`), más un comentario de paquete en hospital/doc.go que explique el modelo.
3. Revisión de Go idiomático y reporte antes de cambiar nada: nombres (sin GetX, sin stuttering tipo hospital.HospitalX), errores en minúscula y sin punto final, cero panic, cero fmt.Print en hospital/, cero lógica de negocio en main.go, receivers consistentes por tipo.
4. Aplica las correcciones que te apruebe, con commits `refactor:` o `docs:`.
Explícame cada hallazgo que encuentres y por qué es o no idiomático (esto alimenta AI_USAGE.md: si algo lo generaste tú mal antes, anótalo en notas_ia.md).
```

### Rama `docs/readme` — Prompt 7.2

**Qué va:** un `README.md` completo.

```
Rama: docs/readme (créala desde dev). Escribe README.md en español con:
- Título, descripción breve, autor (Jerónimo, trabajo individual).
- Requisitos (versión de Go) y cómo correr: `go run .`; cómo testear: `go test ./...`. Deja una sección "Bonus" con `go run . -sim`, `go test -race ./...` y la interfaz Fyne (`cd gui && go run .`) marcados como "(se documenta en las fases 8 y 9)".
- Estructura del proyecto (árbol de carpetas) y qué contiene cada archivo.
- Modelo: tabla con las 5 entidades + Orderly, composición (Person) e interfaz (Attender), con un diagrama mermaid classDiagram.
- "Decisiones de diseño": toma la sección 4 de docs/PLAN.md y escribe cada punto con su justificación en 1–2 líneas.
- Flujo de una crisis de sueño paso a paso (HandleSleepAttack).
- Estrategia de Git: feature branching, prefijos de rama, PR con merge commit.
No inventes funcionalidades que no existan en el código: verifica cada afirmación contra el código. Commit: "docs: complete README".
```

> ✅ **Checkpoint:** después de mergear `docs/readme` en `dev`, haz la **promoción `dev` → `main`** (sección 3.5). Con eso `main` tiene la nota base completa (falta solo `AI_USAGE.md`).

---

## FASE 8 — Bonus: concurrencia (+5 %)

**Por qué va antes de la interfaz:** el bonus de interfaz exige que `hospital/` **no cambie** al agregarla, así que toda modificación al paquete (como el mutex) tiene que hacerse aquí, antes de la Fase 9.

> Nota: la GUI **no** ejecuta `Simulate`. Las goroutines de la simulación modifican pacientes mientras la GUI leería sus getters, y eso sería una data race. Simulación (`go run . -sim`) e interfaz (`gui/`) son programas separados.

> ⚠️ **Requisito en Windows:** `go run -race` necesita CGO y gcc (lo verificaste en la Fase 0). Sin eso, no puedes demostrar que no hay data races y el bonus no cuenta.

### Rama `feature/thread-safe-hospital` 🧭 — Prompt 8.1

**Para qué:** que el hospital sea seguro bajo acceso concurrente. Esta es la pregunta *"¿qué pasa si dos goroutines asignan el mismo cuarto a la vez?"*.

**Qué va:** un `sync.Mutex` en `Hospital`, versiones internas sin lock y un test concurrente con `-race`.

**Alcance:** no cambia el comportamiento de `go run .`.

**Terminado cuando:** `go test -race ./...` pasa y el test concurrente asigna exactamente tantos cuartos como existen.

**Debes entender:** qué es una data race; por qué el mutex de Go **no es reentrante** (de ahí las versiones `...Locked`); por qué `Hospital` nunca se copia (`go vet` detecta copias de un mutex); y qué contrato tienen los métodos de las entidades.

```
Rama: feature/thread-safe-hospital (créala desde dev). Lee docs/PLAN.md "FASE 8". Primero propón el plan y espera mi aprobación.
1. Agrega `mu sync.Mutex` a Hospital. TODOS los métodos exportados de Hospital (registro, getters, búsquedas, AssignRoom, HandleSleepAttack, WakeUpPatient y las consultas) hacen `h.mu.Lock(); defer h.mu.Unlock()`.
2. El sync.Mutex de Go NO es reentrante: si HandleSleepAttack tiene el lock y llama a AssignRoom, que también lo pide, hay deadlock. Crea versiones internas sin lock (assignRoomLocked, dispatchLocked, registerEpisodeLocked, findPatientLocked...) que solo se llaman con el lock tomado, y que los métodos públicos sean wrappers. Documenta esa convención en un comentario.
3. Documenta en el comentario del paquete el contrato de concurrencia: durante ejecución concurrente, Patient/Room/Doctor/Orderly solo se modifican a través de métodos de Hospital; Doctor.MyEpisodes y los getters de entidades se leen cuando no hay goroutines escribiendo (por ejemplo, después de wg.Wait()).
4. Test `TestConcurrentSleepAttacks`: hospital con 3 cuartos, 2 doctores, 1 camillero y 50 pacientes; 50 goroutines llaman HandleSleepAttack al mismo tiempo (sync.WaitGroup). Verifica: exactamente 3 sin error, 47 con ErrNoRoomAvailable, 50 registros en History() y ningún cuarto con más ocupantes que su capacidad.
5. Corre `go test -race ./...` y `go vet ./...`. Si -race no corre en mi máquina, avísame en lugar de saltártelo.
Commits separados (mutex, internos, test).
Explícame qué pasaría en el test SIN el mutex, línea por línea, en AssignRoom.
```

### Rama `feature/concurrent-simulation` — Prompt 8.2

**Para qué:** lo que el bonus pide literalmente: cada paciente es una goroutine que se duerme a intervalos aleatorios.

**Qué va:** `hospital/simulation.go` con `Simulate` y el flag `-sim` en `main`.

**Alcance:** `go run .` sigue mostrando exactamente el escenario obligatorio.

**Terminado cuando:** `go run -race . -sim` termina sin reportes de race.

```
Rama: feature/concurrent-simulation (créala desde dev). Lee docs/PLAN.md "FASE 8".
1. En hospital/simulation.go:
   - `type SimConfig struct { Rounds int; MaxDelay time.Duration; Locations []string }` (campos exportados: es un struct de configuración, no una entidad; explícame por qué esto no rompe la encapsulación).
   - `type SimResult struct { Attacks, RoomFailures, Errors int }`.
   - `func (h *Hospital) Simulate(cfg SimConfig) (SimResult, error)`: valida cfg; lanza UNA goroutine por paciente admitido; cada una repite Rounds veces: dormir un tiempo aleatorio (math/rand, hasta MaxDelay) → h.HandleSleepAttack(p, ubicación aleatoria) → dormir otro tiempo aleatorio → h.WakeUpPatient(p). ErrNoRoomAvailable cuenta como RoomFailures, no como error. Los contadores se protegen con sync/atomic o con un mutex propio. Espera con sync.WaitGroup y retorna el resultado. NO imprime nada.
2. En main.go agrega el flag `-sim` (paquete flag): si está activo, construye el mismo hospital del escenario (extrae la construcción a una función de main, sin lógica de negocio), corre Simulate con Rounds=3 y MaxDelay=50ms e imprime el SimResult y las 4 consultas. Sin el flag, go run . queda idéntico.
3. Test de Simulate con -race: después de Simulate, el total de registros en History() es igual a pacientes × Rounds.
4. Corre `go run -race . -sim` 3 veces y `go test -race ./...`, y muéstrame las salidas.
Commits y explicación: qué hace WaitGroup, por qué los contadores necesitan atomic y por qué las goroutines solo tocan el estado a través de métodos de Hospital.
```

---

## FASE 9 — Bonus: interfaz gráfica con Fyne (+10 %)

**Decisión:** la interfaz es una **app de escritorio nativa con [Fyne](https://fyne.io) (`fyne.io/fyne/v2`)**, escrita 100 % en Go y sin HTML ni JavaScript. Permite **elegir una foto para cada paciente** y tiene un **reproductor de música**. Fyne no trae audio, así que el reproductor usa **`github.com/gopxl/beep/v2`** (mp3, wav y ogg). El enunciado permite librerías de terceros solo para la interfaz, y estas dos se usan solo ahí.

**Arquitectura:**
```
costenos-narcolepsia/            ← módulo principal: SOLO librería estándar
  hospital/   (NO cambia en esta fase)
  main.go     (escenario por consola)
  gui/                           ← módulo Go APARTE (go.mod propio), con Fyne + beep
    main.go, app.go, patients.go, rooms.go, queries.go, history.go, music.go, seed.go
    media/    (paquete sin Fyne: guarda y valida fotos y música → testeable sin gcc)
    player/   (paquete sin Fyne: reproductor con beep)
    data/     (fotos y canciones copiadas; en .gitignore)
```
`gui/go.mod` usa `replace costenos-narcolepsia => ../` para importar `costenos-narcolepsia/hospital` desde la carpeta local.

**Cómo se corre:** `cd gui` y luego `go run .`. La **primera compilación tarda varios minutos** porque compila la parte gráfica con gcc; las siguientes son rápidas.

**Requisito:** gcc + `CGO_ENABLED=1`, el mismo que `-race` (verificado en la Fase 0).

**Reglas del bonus:**
- `hospital/` **no cambia en absoluto**: `git diff dev...HEAD --stat -- hospital/` vacío. Si falta un getter, **para**, crea `fix/hospital-getters`, mergéala y vuelve.
- El **módulo principal no gana dependencias**: `git diff dev...HEAD -- go.mod go.sum` vacío. Fyne y beep solo aparecen en `gui/go.mod`.
- Fotos y música **no tocan el dominio** (decisión 15).
- La GUI **no ejecuta `Simulate`**: mientras las goroutines modifican pacientes, la GUI leería sus getters y habría una data race.
- Los widgets solo se actualizan en el hilo de la interfaz. Los callbacks de botones ya corren ahí; desde otra goroutine (el reproductor) se usa **`fyne.Do(...)`**.

**Orden = prioridad:** `feature/gui-shell` → `feature/gui-hospital-views` (**con estas dos ya tienes el bonus**) → `feature/gui-photos` → `feature/gui-music` → `docs/gui-readme`. Si se acaba el tiempo, corta después de cualquier rama ya mergeada.

**Debes entender (toda la fase):** por qué la GUI es otro módulo y qué hace `replace`; el flujo botón → callback → método de `Hospital` → `refreshAll()`; cómo funciona `widget.List` (length / create / update); por qué existe `fyne.Do`; y por qué las fotos no son un campo de `Patient`.

### Rama `feature/gui-shell` 🧭 — Prompt 9.1

**Para qué:** crear el módulo de la interfaz, la ventana y la estructura de pestañas, sin tocar nada del módulo principal.

**Qué va:** `gui/go.mod` con `replace`, Fyne instalado, un hospital inicial (`seed.go`), un struct `App` y la ventana con pestañas vacías.

**Terminado cuando:** `cd gui && go run .` abre la ventana, y en la raíz `go test ./...` sigue igual con `go.mod` sin cambios.

```
Rama: feature/gui-shell (créala desde dev). Lee docs/PLAN.md "FASE 9" completa y "Decisiones de diseño" (puntos 15 y 16). Primero propón el plan y espera mi aprobación. NO toques hospital/, main.go ni el go.mod de la raíz.
1. Crea el módulo gui/: `go mod init costenos-narcolepsia/gui` dentro de gui/, con `require costenos-narcolepsia v0.0.0` y `replace costenos-narcolepsia => ../`. Instala Fyne con `go get fyne.io/fyne/v2@latest`. Explícame qué hace replace y por qué la GUI es un módulo aparte.
2. gui/seed.go: `func newSeedHospital() (*hospital.Hospital, error)` que arma el MISMO estado inicial del paso 1 del escenario (2 doctores, 1 camillero, 3 cuartos, 5 pacientes con los mismos IDs y nombres que main.go), sin ataques. Comentario: es una duplicación intencional de datos (no de lógica) porque la GUI es otro programa y otro módulo.
3. gui/app.go: `type App struct` con el *hospital.Hospital, la fyne.Window, el último resultado de asignación de cuarto y las funciones de refresco de cada vista; un método `refreshAll()` y helpers `showError(err)` y `showWarning(msg)` con dialog.ShowError / dialog.ShowInformation.
4. gui/main.go: app.NewWithID("co.edu.eia.costenos"), ventana "Hospital de los Costeños con Narcolepsia" de 1100x720, un encabezado con contadores (pacientes, dormidos en pasillo, dormidos en cama, cuartos libres) y container.NewAppTabs con las pestañas Pacientes, Cuartos, Consultas, Historial y Música (por ahora con un label "Próximamente"), cada una con su icono de theme.
5. Verifica: `cd gui && go build .` compila (avísame si falla por gcc/CGO); en la raíz `go test ./...` y `go vet ./...` pasan; `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` sale vacío. Agrega gui/data/ al .gitignore.
Commits (module, seed, app shell). Explícame la estructura y el ciclo de vida de una app Fyne (app → window → content → ShowAndRun).
```

### Rama `feature/gui-hospital-views` — Prompts 9.2a y 9.2b

**Para qué:** la interfaz funcional con pacientes, cuartos, las 4 consultas e historial. **Con esta rama mergeada ya cumples el bonus.**

**Qué va:** un archivo por vista: `patients.go`, `rooms.go`, `queries.go` y `history.go`.

**Terminado cuando:** desde la GUI puedes reproducir el escenario (4 ataques con 3 cuartos → aviso en el 4.º → despertar → asignar cuarto) y ver las 4 consultas actualizadas.

**Prompt 9.2a — Pestaña Pacientes**
```
Rama: feature/gui-hospital-views (créala desde dev). Lee docs/PLAN.md "FASE 9". NO toques hospital/ ni el go.mod raíz.
En gui/patients.go construye la pestaña Pacientes con un container.NewHSplit:
1. Izquierda: widget.List de pacientes. Cada fila tiene un avatar circular con iniciales (canvas.Circle + canvas.Text; la foto llega en otra rama), nombre, ID y una "pastilla" de estado con color: Awake = verde, AsleepInHallway = ámbar, AsleepInBed = índigo. Usa un helper `stateColor(s hospital.PatientState) color.Color`.
2. Derecha: panel de detalle del paciente seleccionado (widget.Card) con ID, edad, nivel, estado, ubicación, cuarto y doctor asignado. Acciones:
   - Entry "Ubicación" + botón "Sleep attack" (habilitado solo si está despierto) → h.HandleSleepAttack. Si el error es ErrNoRoomAvailable (errors.Is), usa showWarning ("Episodio registrado, pero no hay cuarto: el paciente queda en el pasillo"); cualquier otro error, showError.
   - "Wake up" (solo si está dormido) → h.WakeUpPatient.
   - "Assign room" (solo si está en el pasillo) → h.AssignRoom.
   - Después de cada acción: guarda el resultado de asignación en App y llama a refreshAll().
3. Habilita y deshabilita los botones según el estado en cada refresco.
Prueba el flujo del escenario a mano. Verifica `cd gui && go build .` y `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío. Commits.
Explícame cómo funciona widget.List (los callbacks length, createItem y updateItem y por qué reutiliza filas) y por qué aquí no hace falta polling.
```

**Prompt 9.2b — Cuartos, Consultas e Historial**
```
Seguimos en feature/gui-hospital-views. Un archivo por vista:
1. gui/rooms.go: grid de 3 columnas (container.NewGridWithColumns) con una widget.Card por cuarto: número, capacidad, estado (con color: libre verde, ocupado índigo) y nombre del ocupante o "Libre".
2. gui/queries.go: pestaña Consultas con un widget.Accordion (o AppTabs internas) con las 4 consultas rotuladas como en el enunciado:
   - "Pacientes en el pasillo" → h.PatientsInHallway() (lista; "(ninguno)" si está vacía).
   - "Episodios por doctor" → widget.Select con los doctores + lista de d.MyEpisodes() con Summary().
   - "Asignación de cuarto" → el último resultado guardado en App (cuarto asignado o el error/aviso).
   - "Reporte de severidad" → widget.Table con paciente y episodios de h.SevereReport(), ORDENADO por ID (explícame por qué hay que ordenar).
3. gui/history.go: lista de h.History() del más reciente al más antiguo, con Summary().
4. refreshAll() actualiza todas las vistas y los contadores del encabezado.
5. Reproduce el escenario completo desde la GUI y compáralo con la salida de `go run .` en la raíz.
Verifica build, `go test ./...` en la raíz y `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío. Commits por vista.
Explícame cómo la GUI usa exactamente los mismos métodos que main.go y por qué eso demuestra la separación lógica/vista.
```

### Rama `feature/gui-photos` — Prompt 9.3

**Para qué:** el primer plus visual: una foto por paciente, elegida con el diálogo de archivos.

**Qué va:** el paquete `gui/media` (sin Fyne y testeable) con `SavePhoto`, más el uso de la foto en la lista y en el detalle.

**Terminado cuando:** eliges una foto, se ve en la lista y en el detalle, y sigue ahí al reabrir la app.

```
Rama: feature/gui-photos (créala desde dev). Lee docs/PLAN.md "FASE 9" y "Decisiones de diseño" (punto 15). NO toques hospital/ ni el go.mod raíz.
1. Paquete gui/media (SIN importar Fyne, para poder testearlo sin gcc): `type Store struct { dir string }`, `NewStore(dir string) (*Store, error)` que crea dir/photos y dir/music.
   - `SavePhoto(patientID, srcPath string) (string, error)`: límite de 5 MB (ErrTooLarge); valida el CONTENIDO con http.DetectContentType sobre los primeros 512 bytes: solo image/png o image/jpeg, que son los que Fyne decodifica (si no, ErrUnsupportedType). Copia a photos/<ID>.<ext> con la extensión derivada del tipo detectado, nunca del nombre original, y borra la foto anterior con otra extensión. Valida que patientID no tenga separadores de ruta.
   - `PhotoPath(patientID string) (string, bool)`.
   - Tests con t.TempDir(): PNG válido (bytes de un PNG mínimo), un .txt renombrado a .png → ErrUnsupportedType, archivo grande → ErrTooLarge, ID con "../" → error. Corre `cd gui && go test ./media/`.
2. En la GUI usa gui/data como directorio del Store.
   - Detalle del paciente: foto grande (canvas.NewImageFromFile, FillMode ImageFillContain, SetMinSize 160x160) o el avatar de iniciales si no tiene. Botón "Cambiar foto" → dialog.ShowFileOpen con storage.NewExtensionFileFilter([".png", ".jpg", ".jpeg"]) → Store.SavePhoto → refresco. Los errores van a showError con un mensaje claro.
   - Lista: avatar pequeño con la foto si existe.
3. Verifica build, tests y `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío. Commits (media store, tests, UI).
Explícame por qué se valida el contenido y no la extensión, qué es path traversal y por qué la foto no es un campo de Patient.
```

### Rama `feature/gui-music` 🧭 — Prompts 9.4a y 9.4b

**Para qué:** el segundo plus visual: un reproductor de música con biblioteca propia.

**Qué va:** `SaveTrack`/`Tracks` en `gui/media`, el paquete `gui/player` con beep y la pestaña Música.

**Terminado cuando:** puedes agregar canciones, reproducir, pausar, ir a la siguiente o anterior, adelantar, cambiar el volumen y borrar.

**Prompt 9.4a — Biblioteca y reproductor (sin interfaz)**
```
Rama: feature/gui-music (créala desde dev). Lee docs/PLAN.md "FASE 9". Primero propón el plan y espera mi aprobación. NO toques hospital/ ni el go.mod raíz.
1. gui/media: `SaveTrack(srcPath string) (Track, error)` con una lista blanca de extensiones (.mp3, .wav, .ogg), un límite de 20 MB y rechazo si DetectContentType detecta text/*, image/*, application/pdf o application/zip. Nombre en disco: <unix>-<nombre saneado> (a-z, 0-9, guion, guion bajo; máx. 60). `Tracks() ([]Track, error)` ordenadas por fecha y `DeleteTrack(id string) error` (el id debe existir en el listado). `type Track struct { ID, Title, Path string }`. Tests con t.TempDir().
2. Paquete gui/player (SIN Fyne) con `go get github.com/gopxl/beep/v2`:
   - Inicializa speaker una sola vez con una frecuencia fija (44100 Hz) y usa beep.Resample cuando la canción tenga otra.
   - `type Player struct` con métodos: Load(path) (decodifica según la extensión con beep/v2/mp3, wav o vorbis), Play(), Pause(), Toggle(), Stop(), Seek(fraction float64), Position() y Length() (time.Duration), SetVolume(0..1) con effects.Volume, IsPlaying() y un callback OnFinished(func()) usando beep.Seq + beep.Callback.
   - Todo acceso al streamer va con speaker.Lock()/Unlock(). Cierra los streamers al cambiar de canción.
   - Doc comments explicando que speaker reproduce en su propia goroutine.
3. Verifica `cd gui && go build ./... && go test ./media/` y `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío. Commits.
Explícame por qué speaker.Lock() es necesario y en qué se parece al mutex del Hospital.
```

**Prompt 9.4b — Pestaña Música**
```
Seguimos en feature/gui-music. En gui/music.go construye la pestaña Música:
1. Lista de canciones (widget.List) con la canción actual resaltada; doble clic o selección + "Play" la reproduce.
2. Controles: anterior, play/pausa (el icono cambia), siguiente, slider de progreso con tiempo actual / total, slider de volumen, "Agregar canción" (dialog.ShowFileOpen filtrando .mp3 .wav .ogg → Store.SaveTrack) y "Eliminar".
3. Un ticker en una goroutine actualiza el progreso cada 500 ms con fyne.Do(...) (explícame por qué no se puede tocar un widget directo desde esa goroutine). Evita el bucle de retroalimentación cuando el usuario arrastra el slider (bandera mientras se arrastra, y Seek en OnChangeEnded).
4. Al terminar una canción pasa a la siguiente (OnFinished → fyne.Do).
5. Opcional: window.SetOnDropped para agregar canciones arrastrándolas a la ventana.
6. Detén el ticker y el reproductor al cerrar la ventana (SetCloseIntercept o SetOnClosed).
Prueba con 2 canciones: reproducir, pausar, adelantar, siguiente y borrar. Verifica build y `git diff dev...HEAD --stat -- hospital/ go.mod go.sum` vacío. Commits.
Explícame el recorrido completo: clic en Play → Player → speaker (otra goroutine) → ticker → fyne.Do → slider.
```

### Rama `docs/gui-readme` — Prompt 9.5

**Para qué:** que el README documente la interfaz que **realmente** quedó en `dev`.

```
Rama: docs/gui-readme (créala desde dev). Actualiza README.md con una sección "Interfaz gráfica (bonus)" que describa SOLO lo que existe en dev: stack (Fyne v2 + gopxl/beep para audio), por qué es un módulo aparte (gui/go.mod con replace), requisitos (gcc + CGO en Windows, con un enlace a la guía de instalación de Fyne), cómo correrla (`cd gui && go run .`, avisando que la primera compilación tarda), qué permite hacer (pestañas, fotos y música), dónde se guardan los archivos (gui/data, ignorado por git) y cómo se comprueba que hospital/ y el go.mod raíz no cambiaron (con los comandos git diff). Deja un espacio para una captura de pantalla (docs/screenshot.png, que tomaré yo). Commit: "docs: GUI section in README".
```

---

## FASE 10 — Cierre y entrega

### Rama `docs/ai-usage` — Prompt 10.1

**Para qué:** `AI_USAGE.md` es obligatorio. Si contradice lo que digas en la sustentación, se trata como deshonestidad académica (factor 0).

```
Rama: docs/ai-usage (créala desde dev). Ayúdame a redactar AI_USAGE.md en español, honesto, a partir de notas_ia.md y del historial (git log). Estructura:
1. Herramientas usadas y para qué: Claude (chat, planificación con docs/PLAN.md) y Claude Code (código Go, tests, docs e interfaz Fyne con fotos y reproductor de música), indicando por fase.
2. 2–3 prompts decisivos, copiados de docs/PLAN.md o de los que realmente usé, explicando por qué fueron decisivos.
3. Al menos 2 casos reales donde el código generado estaba mal o no era idiomático, sacados de notas_ia.md: qué estaba mal, cómo lo detecté y cómo se corrigió (con el commit si existe). NO inventes casos: si notas_ia.md tiene menos de 2, dímelo y los buscamos juntos en el historial.
4. Qué aprendí (déjame una lista de temas con espacios para escribirlo yo con mis palabras).
Commit: "docs: AI usage report".
```

### Promoción final `dev` → `main` y Prompt 10.2 (verificación final, sin cambios)

Primero haz la **promoción final `dev` → `main`** (sección 3.5) y vuelve a poner `main` como rama por defecto (`gh repo edit --default-branch main`, o desde Settings), para que quien abra o clone el link vea `main`. Después, el Prompt 10.2:

```
Estamos en main actualizado. NO modifiques archivos. Haz la verificación final:
1. Clona el repositorio remoto en una carpeta temporal fuera del proyecto y, ahí, corre `go run .`, `go test ./...`, `go vet ./...` y `gofmt -l .` en la raíz (deben funcionar sin descargar Fyne: el módulo principal es solo librería estándar). Si están disponibles, también `go test -race ./...` y `go run -race . -sim`. Después, en ese clon, `cd gui && go build .` y abre la app para comprobar que arranca (ciérrala al terminar).
2. Recorre el checklist de la sección 10 de docs/enunciado.pdf punto por punto y dime ✅/❌ con evidencia (archivo:línea o salida del comando). Agrega: `git log --oneline -- hospital/ go.mod go.sum` no muestra commits de las ramas de la Fase 9.
3. Muéstrame `git log --oneline --graph --all | head -60` para confirmar que el historial muestra las ramas.
4. Si todo está ✅, crea el tag anotado `v1.0.0` ("Entrega final") y haz push del tag. Si hay ❌, dime qué rama fix/ crear (desde dev, PR a dev y luego una nueva promoción a main).
5. Borra la carpeta temporal.
```

Luego: **entrega el link del repositorio antes de las 23:59.**

---

## FASE 11 — Sustentación (mañana temprano)

Practica en ramas `practice/*` creadas desde `main`. **Nunca se mergean** y al terminar se borran: `git switch main && git branch -D practice/<nombre>`.

**La meta es hacerlo tú.** Usa a Claude solo como evaluador: que te pregunte, que te corrija y que no escriba el código.

**Prompt 11.1 — Simulacro**
```
Eres el profesor evaluando mi sustentación de este proyecto (lee docs/enunciado.pdf sección 8.2 y el código). NO escribas código. Hazme UNA pregunta a la vez de este tipo, espera mi respuesta y califícala (correcta / parcial / incorrecta, con la explicación correcta):
- Value vs pointer embedding de Person.
- Cambiar un receiver de puntero a valor y predecir qué se rompe.
- Por qué un campo está en minúscula y qué exactamente impide tocarlo desde otro paquete.
- Qué pasa si dos goroutines asignan el mismo cuarto a la vez.
- Por qué map y no slice, por qué ese tipo de retorno, por qué esa separación de paquetes.
- Preguntas sobre mi AI_USAGE.md.
Después de 8 preguntas, dame un factor estimado (1.0/0.8/0.5) y en qué temas debo repasar.
```

**Modificaciones en vivo para practicar** (cada una en su rama `practice/`, cronometrándote):

| Rama | Ejercicio | Meta |
|---|---|---|
| `practice/nurse` | Crear `Nurse` (Person embebido + `Attend`), contratarlo con `HireStaff` y verlo despachado en el escenario | < 5 min sin ayuda |
| `practice/value-receiver` | Cambiar `Attend` de `*Doctor` a `Doctor`, compilar y correr: ¿qué pasa con `MyEpisodes`? ¿Sigue cumpliendo `Attender`? | Predecirlo **antes** de correr |
| `practice/pointer-person` | Cambiar `Person` a `*Person` en `Patient`: ¿qué cambia en los constructores? ¿Qué pasa si queda nil? | Explicarlo |
| `practice/no-mutex` | Quitar el lock de `AssignRoom` y correr `go test -race ./...` | Leer e interpretar el reporte de race |
| `practice/export-field` | Intentar `p.state = hospital.Awake` desde `main.go` | Leer el error del compilador y explicarlo |
| `practice/gui-button` | Agregar en la GUI un botón "Cuartos libres" que muestre en un diálogo los cuartos disponibles, sin tocar `hospital/` | Explicar el flujo callback → Hospital → widget |

**Tabla de respuestas rápidas:**

| Pregunta | Respuesta corta |
|---|---|
| `Person` por valor vs `*Person` | Valor: se copia con el struct y nunca es nil. Puntero: se comparte entre structs, puede ser nil (panic al usarlo) y un cambio se ve en todos los que lo comparten. |
| Receiver puntero → valor | El método modifica una copia y el cambio se pierde (`d.episodes` no crece y `MyEpisodes` sale vacío). Con receiver de valor, `T` y `*T` cumplen la interfaz; al revés, si un método pasa a puntero, un `T` por valor deja de cumplirla. |
| ¿Por qué minúscula? | No exportado: el **compilador** rechaza el acceso desde otro paquete. La encapsulación en Go es por paquete, no por tipo. |
| Dos goroutines y el mismo cuarto | Sin mutex: data race, ambas ven `IsAvailable()==true` y el cuarto queda con 2 ocupantes. Con mutex: la segunda espera, ve el cuarto ocupado y recibe `ErrNoRoomAvailable`. |
| ¿Por qué un map en `SevereReport`? | Relación directa paciente→conteo con búsqueda O(1), y la clave puntero identifica al paciente único. Costo: no tiene orden, por eso `main` ordena. |
| ¿Por qué `hospital` y `main` separados? | Lógica reutilizable y testeable, sin depender de cómo se muestra. Lo prueba el bonus: la consola, la simulación y la GUI Fyne usan el mismo paquete sin cambiarlo. |
| ¿Por qué `(*Room, error)`? | Puntero: es el cuarto real que quedó ocupado, no una copia. `error`: "no hay cuarto" es un caso previsible, no un `panic`. |
| ¿Cómo se cumple que `hospital/` no cambió con la interfaz? | `git log -- hospital/ go.mod` y `git diff` por rama: la Fase 9 solo agregó `gui/`, que es otro módulo. La GUI usa los métodos y getters públicos que ya existían. |
| ¿Por qué la foto no es un campo de `Patient`? | Es un dato de presentación: el dominio no necesita imágenes para atender episodios. Vive en `gui/media` y en `gui/data/`, asociada por ID, y así el modelo no depende de la vista. |
| ¿Cómo viaja una foto elegida? | `dialog.ShowFileOpen` → `media.Store.SavePhoto` (límite de tamaño + `DetectContentType` sobre el contenido real) → copia a `gui/data/photos/<ID>.<ext>` → `canvas.NewImageFromFile` en el detalle y en la lista. |
| ¿Por qué la GUI es otro módulo? | Para que el módulo principal siga siendo solo librería estándar y compile sin gcc, y para que la GUI solo pueda usar lo que `hospital` exporta. `replace` apunta al código local. |
| ¿Para qué es `fyne.Do`? | Los widgets solo se tocan desde el hilo de la interfaz. El ticker del reproductor corre en otra goroutine, así que pide la actualización con `fyne.Do`. |

Repasa tu `AI_USAGE.md` para que lo que digas coincida exactamente con lo escrito.
