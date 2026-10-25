# Architecture

labdon is a **Go** service following the layered model taught by
`.agents/skills/layered-architecture-types` (see DECISIONS.md ADR-0002).

## Layers & direction of dependency

```
cmd/sites        entrypoint: config, DI wiring, graceful shutdown
   │
   ▼
app              HTTP layer — handlers, middleware. Primitive JSON types only.
   │
   ▼
business         domain logic — strong types (business/types/*), services.
   │             No cross-domain imports. Extensions wrap via ExtBusiness seam.
   ▼
foundation       shared infra — logger, trace, database, buildinfo.
stores (*db)     persistence inside each domain — native SQL row types only.
```

Rules:
- **Primitives live at the edges**: API request/response structs and DB row
  structs use `string`, `int`, `bool`, `time.Time`, `json.RawMessage`, `sql.Null*`.
- **Strong types live only in the Business layer** (`business/types/*`).
- **Every crossing goes through a named converter** — never assign across a
  boundary directly:
  | Direction | Converter | Notes |
  |---|---|---|
  | App → Business | `toBus<Type>` | parses + validates, returns `errs.FieldErrors` |
  | Business → App | `fromBus<Type>Response` | strong → primitive, explicit |
  | Business → Storage | `toDB<Type>` | |
  | Storage → Business | `toBus<Type>` | native → strong, returns error |
- Handlers only decode/validate/call/respond; status mapping is centralized.
- Cross-cutting concerns (OTEL, logging, metrics, caching, auth) wrap the core
  `Business` through the `ExtBusiness`/`Extension` seam without modifying it.

Full details, examples and conformance checklists: the
`layered-architecture-types` and `business-layer-extensions` skills.
