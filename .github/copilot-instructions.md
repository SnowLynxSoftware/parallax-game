# AI coding agents: parallax-game quick guide

Purpose: Help AI contributors be productive immediately in this Go monolith by capturing the project’s architecture, workflows, and house patterns.

## Big picture architecture

- Go 1.25 monolith with chi router. Entry: `main.go` → `cmd/handler.go` → default command `server` → `server/app.server.go`.
- Dependency wiring happens in `NewAppServer.Start()`:
  - Config: `config/app-config.go` (env-driven; panics if required vars missing).
  - DB: PostgreSQL via `sqlx` in `server/database`. Repos live in `server/database/repositories`.
  - Services live in `server/services` and encapsulate business logic; controllers call services.
  - Middleware: `server/middleware/auth.middleware.go` validates JWT from cookie and loads user.
  - UI: HTML templates + static assets are embedded via `embed` in services and served by controllers.
- External deps: chi, cors, sqlx, lib/pq, zerolog, golang-jwt, mailjet, heimdall (SQL migrations).

## Runtime & workflows

- Port: 3000. Start locally: `make run` (defaults to `server`). Migrate: `make migrate` (`go run . migrate`).
- Docker: `docker-compose.yml` provides Postgres (5432) and Redis (not yet used in code). Build/push images with `make docker-build`/`docker-push`.
- Tests: `make test`, coverage: `make test-coverage`. Lint/vet/format via `make fmt | vet | lint | check`.
- Required env (see `config/app-config.go`): `CLOUD_ENV`, `DB_CONNECTION_STRING`, `AUTH_HASH_PEPPER`, `JWT_SECRET_KEY`, `MJ_APIKEY_PUBLIC`, `MJ_APIKEY_PRIVATE`. Optional: `BASE_URL`, `CORS_ALLOWED_ORIGIN`, `COOKIE_DOMAIN`, `DEBUG_MODE`.

## Routing, auth, and responses

- Router mounts (see `server/app.server.go`):
  - API: `/api/health`, `/api/auth`, `/api/users`, `/api/links`.
  - UI pages: `/`, `/welcome`, `/register`, `/login`, `/dashboard`, `/terms`, `/privacy`.
  - Public redirect: `/ur/{short_code}` (unregistered short links).
- Auth: Access token is a JWT (HS512) stored in cookie `access_token`. `AuthMiddleware.Authorize(r)` returns immutable user context. Cookies are `Secure` when `CLOUD_ENV != local` and use `COOKIE_DOMAIN`.
- Controllers return JSON for API and render embedded templates for UI. Errors to clients are generic; details are logged via `server/util` and zerolog.

## House patterns to follow

- Layering: Controller → Service → Repository. Keep SQL only in repositories, mapping to entity structs (`repositories.*Entity`). Services return DTOs in `server/models`.
- Pagination: query params `page_size` and `page` → compute `offset`; see `UserController.getUsers` and `ILinksRepository.GetLinks*` for patterns.
- Dates/formatting: responses format times as RFC3339-like strings (e.g., `2006-01-02T15:04:05Z`).
- Feature flags: check via `FeatureFlagService.IsEnabled(key)`; missing flag returns false. Used in UI (e.g., `prelaunch_mode`).
- Short links: Displayed URL uses `/r/{customSlug}` if present else `/ur/{shortCode}`; generation via `common.LinkEngine`.

## Adding a new API feature (example)

1. Create a controller in `server/controllers/*.controller.go` implementing `IController.MapController()` and define routes with chi. Use `authMiddleware.Authorize(r)` for protected endpoints.
2. Add a repository in `server/database/repositories` using `sqlx` with parameterized SQL and entity structs.
3. Add a service in `server/services` to validate inputs, call the repo, and map to DTOs in `server/models`.
4. Wire dependencies and mount routes in `server/app.server.go` (instantiate repo/service/middleware, then `router.Mount("/api/thing", controller.MapController())`).
5. If DB schema changes, add a migration SQL file under `migrations/` and run `make migrate`.

## Gotchas

- App fails fast if required env vars are missing; `.env` is auto-loaded via `github.com/joho/godotenv/autoload` during development.
- `BASE_URL` affects short link URLs constructed in `LinkService.mapToResponseDTO`.
- Cookie domain must match the host you use in the browser or auth will fail silently.
- Redis is defined in compose but not used yet; avoid assuming a cache layer exists.
