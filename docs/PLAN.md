# GoArch Visualizer — план работ

> **Единый источник истины** для разработки и handoff между машинами.  
> Обновляй этот файл при завершении задачи или смене фокуса.

---

## Как работать с агентом

Скопируй в новый чат на любой машине:

```
Читай docs/PLAN.md в репозитории goarch-visualizer.
Остановился на задаче <ID>. Продолжаем её.
Не меняй код и не запускай команды без моей явной команды («примени», «запусти»).
```

Примеры:

- `Остановился на задаче C-01. Продолжаем.`
- `Задача C-02 done. Переключаемся на C-03.`
- `Покажи статус по docs/PLAN.md, что дальше?`

**Связанные файлы:**

| Файл | Назначение |
|------|------------|
| `docs/PLAN.md` | план, задачи, решения, текущий фокус |
| `readme.md` | публичное описание, установка, demo |
| `.cursor/mcp.json` | project-level MCP (может не подхватываться Cursor) |
| `~/.cursor/mcp.json` | **глобальный MCP** (рабочий вариант на Mac) |
| `~/.cursor/permissions.json` | allowlist для `analyze_project` |

---

## Текущий фокус

| Поле | Значение |
|------|----------|
| **Активная задача** | `C-01` — RuleAnalyzer (детерминированные правила слоёв) |
| **Фаза** | 1 — MVP (demo для экспертов) |
| **Обновлено** | 2026-09-17 |
| **Блокеры** | нет |

---

## Зафиксированные решения

Не менять без явного обсуждения.

| # | Решение | Почему |
|---|---------|--------|
| D-01 | **AI внутри Go** (`interface Analyzer`), один pipeline для CLI и MCP | Архитектура B; не два round-trip через Cursor |
| D-02 | **Claude API не обязателен** | Ollama дома, mock на работе; Cursor сам — LLM host |
| D-03 | **Package ID = import path** | `sample/internal/service`, не `pkg:service` |
| D-04 | **Mermaid рисует только packages** | funcs/types в графе есть, но не на диаграмме |
| D-05 | **Ollama получает `PackageGraph`** | полный граф → timeout на MCP |
| D-06 | **Validate отсекает галлюцинации LLM** | whitelist node ID из реального графа |
| D-07 | **MVP = rule-based validation**, LLM — опционально | mock/Ollama не дают честных violations |
| D-08 | **`go/packages` — после MVP** | WalkDir + module.go достаточно для demo |

---

## Pipeline — статус компонентов

```
AnalyzeDir (WalkDir + ast)  →  Analyzer  →  Validate  →  ToMermaid  →  CLI / MCP
```

| Компонент | Путь | Статус | Примечание |
|-----------|------|--------|------------|
| Graph model | `internal/analyzer/graph.go` | ✅ | Nodes, Edges, PackageGraph |
| Module resolve | `internal/analyzer/module.go` | ✅ | go.mod, import path |
| AST parse | `internal/analyzer/ast.go` | ✅ | packages, funcs, types, imports |
| WalkDir loader | `internal/analyzer/dir.go` | ✅ | skip vendor, _test.go |
| go/packages loader | — | ❌ | задача D-01 |
| MockAnalyzer | `internal/ai/mock.go` | 🟡 | всем packages `ok` |
| OllamaAnalyzer | `internal/ai/ollama.go` | ✅ | qwen2.5-coder:7b, PackageGraph |
| ClaudeAnalyzer | — | ❌ | задача E-01 |
| RuleAnalyzer | — | ❌ | **задача C-01 (активная)** |
| Validate | `internal/ai/validate.go` | ✅ | |
| Mermaid | `internal/renderer/mermaid.go` | ✅ | safeId: `/`, `:`, `.` → `_` |
| HTML report | — | ❌ | задача F-02 |
| MCP server | `internal/mcp/server.go` | ✅ | tool `analyze_project` |
| CLI | `cmd/goarch/main.go` | ✅ | `analyze <path>` |

Легенда: ✅ готово · 🟡 работает, но не для MVP · ❌ не начато

---

## Бэклог задач

### Фаза 0 — Инфраструктура ✅

| ID | Задача | Статус | DoD |
|----|--------|--------|-----|
| A-01 | Graph + AST + import path ID | ✅ | тесты зелёные, edges сходятся |
| A-02 | MCP tool `analyze_project` | ✅ | Cursor видит tool, зелёный статус |
| A-03 | CLI `go run ./cmd/goarch analyze` | ✅ | Mermaid в stdout |
| A-04 | AI interface + Mock + Ollama | ✅ | `GOARCH_AI=ollama\|mock` |
| A-05 | Validate (anti-hallucination) | ✅ | чужие node ID отбрасываются |
| A-06 | Mermaid renderer + classDef | ✅ | ok/warning/error цвета |
| A-07 | PackageGraph для Ollama | ✅ | bookings не timeout |
| A-08 | safeId fix (точки в ID) | ✅ | Mermaid парсится без «github» артеfact |
| A-09 | MCP в Cursor (global config) | ✅ | `go run -C ...`, allowlist |
| A-10 | README | ✅ | установка, demo, честные ограничения |

---

### Фаза 1 — MVP (demo для экспертов) 🔄

| ID | Задача | Статус | Файлы | DoD |
|----|--------|--------|-------|-----|
| **C-01** | **RuleAnalyzer** | 🔄 **активная** | `internal/ai/rules.go`, `main.go` | 3–5 правил слоёв; default analyzer; unit-тесты |
| C-02 | Детекция циклов | ⬜ | `internal/analyzer/cycle.go` или в rules | цикл A→B→A = error |
| C-03 | Fixture `testdata/violation-project` | ⬜ | `testdata/violation-project/` | намеренные нарушения; тест красный |
| C-04 | MCP: violations + Mermaid | ⬜ | `internal/mcp/server.go` | текст нарушений + diagram в одном ответе |
| C-05 | Mermaid labels (path, не `main`) | ⬜ | `internal/renderer/mermaid.go` | `[internal/domain]` вместо `[domain]` |
| C-06 | `go build` в MCP config | ⬜ | readme, `~/.cursor/mcp.json` | бинарь вместо `go run` |
| C-07 | Demo на `bookings` | ⬜ | — | `port→dto` = warning от rules, не LLM |

**Definition of Done для фазы 1:**

- [ ] `analyze_project` на `bookings` показывает ≥1 warning (rule-based)
- [ ] `violation-project` — красные узлы, тесты зелёные
- [ ] Работает offline без Ollama
- [ ] CLI и MCP дают одинаковый результат

---

### Фаза 2 — Качество графа

| ID | Задача | Статус | DoD |
|----|--------|--------|-----|
| D-01 | Loader на `go/packages` | ⬜ | build tags; env `GOARCH_LOADER=packages` |
| D-02 | Subgraph в Mermaid | ⬜ | группировка cmd / domain / adapter |
| D-03 | Параметр `level` в MCP | ⬜ | packages / types / functions |

---

### Фаза 3 — Дополнительные инструменты

| ID | Задача | Статус | DoD |
|----|--------|--------|-----|
| E-01 | Claude API analyzer | ⬜ | `GOARCH_AI=claude`, `ANTHROPIC_API_KEY` |
| F-01 | `.goarch.yaml` конфиг правил | ⬜ | кастомные layer rules |
| F-02 | HTML report generator | ⬜ | `internal/generator/html.go` |
| F-03 | MCP tool `check_violations` | ⬜ | текстовый отчёт без diagram |
| F-04 | Edges `call` / `implement` | ⬜ | из AST type info |

---

## Задача C-01 — детали (активная)

**Цель:** детерминированный анализатор архитектуры без LLM.

**Правила MVP (черновик):**

| # | Правило | Severity |
|---|---------|----------|
| R-01 | `domain` не импортирует другие пакеты модуля (кроме stdlib) | error |
| R-02 | `port` не импортирует `adapter/*`, `cmd/*` | error |
| R-03 | `port` → `dto` / `adapter/http/*` | warning |
| R-04 | `adapter/postgres` → `dto` / HTTP слой | warning |
| R-05 | цикл между packages | error |

**Шаги реализации:**

1. Создать `internal/ai/rules.go` — `type RuleAnalyzer struct{}`
2. Реализовать `Analyze(ctx, g, focus) (*Analysis, error)`
3. Использовать `analyzer.PackageGraph(g)` для import edges
4. Определять слой по substring в import path (`/domain`, `/port`, `/adapter/`, `/cmd/`)
5. Default в `main.go`: `RuleAnalyzer` вместо `MockAnalyzer`
6. Env override: `GOARCH_AI=ollama|mock|rules` (rules = default)
7. Тест: `bookings` → warning на `port→dto`

**Не делать в C-01:**

- `.goarch.yaml` (hardcode правил достаточно)
- Claude API
- go/packages

---

## Окружение

### Переменные

| Env | Значения | Default |
|-----|----------|---------|
| `GOARCH_AI` | `rules`, `mock`, `ollama` | `rules` (после C-01; сейчас `mock`) |

### MCP — глобальный конфиг (`~/.cursor/mcp.json`)

```json
{
  "mcpServers": {
    "goarch-visualizer": {
      "command": "go",
      "args": ["run", "-C", "/ABS/PATH/goarch-visualizer", "./cmd/goarch"],
      "env": { "GOARCH_AI": "ollama" }
    }
  }
}
```

> Замени `/ABS/PATH/` на свой путь. Флаг `-C` обязателен.  
> Для быстрого старта (C-06): `"command": "/ABS/PATH/goarch-visualizer/bin/goarch"`.

### Allowlist (`~/.cursor/permissions.json`)

```json
{
  "mcpAllowlist": ["goarch-visualizer:analyze_project"]
}
```

### Ollama (если `GOARCH_AI=ollama`)

```bash
ollama pull qwen2.5-coder:7b
ollama serve
```

---

## Тестовые проекты

| Path | Назначение | Ожидание |
|------|------------|----------|
| `testdata/sample-project` | минимальный happy path | 2 packages, 1 edge, быстро |
| `/Users/ermakov/GO/bookings` | реальный pet project | 13 packages, hexagonal |
| `testdata/violation-project` | нарушения (C-03) | ❌ ещё не создан |

---

## Журнал прогресса

| Дата | Задача | Что сделано |
|------|--------|-------------|
| 2026-09-07 | A-01..A-04 | analyzer, AI interface, Ollama smoke |
| 2026-09-15 | A-02, A-09 | MCP в Cursor, global mcp.json |
| 2026-09-16 | A-06, A-07 | Mermaid pipeline в MCP, Ollama env |
| 2026-09-17 | A-07, A-08 | PackageGraph filter, safeId dots |
| 2026-09-17 | A-10 | README |
| 2026-09-17 | — | Создан `docs/PLAN.md`, фокус → C-01 |

---

## Что НЕ в плане ( явно отложено )

- Два round-trip AI через Cursor (архитектура A+)
- MCP sampling (deprecated / ненадёжно)
- Отправка полного AST в LLM (только PackageGraph)
- Платный Claude API как requirement

---

## Быстрая проверка «всё работает»

```bash
# 1. Analyzer
go test ./internal/analyzer/...

# 2. CLI mock (быстро)
go run ./cmd/goarch analyze testdata/sample-project

# 3. CLI ollama (~15 сек)
GOARCH_AI=ollama go run ./cmd/goarch analyze testdata/sample-project

# 4. MCP — в Cursor:
# «Вызови analyze_project с path: .../testdata/sample-project»
```
