# VSRT-LANG

Spaced-repetition language learning. User builds study sessions of words; the system fills translation, meanings, synonyms, antonyms, and usage contexts.

## Commands

Needs a `.env` at repo root. Variable lists: [GO-Manager config](GO-Manager/docs/configuration.md), [Translator env](Translator/docs/Readme.md).

```bash
docker-compose build
docker-compose up
```

First boot can kill `translator` if LibreTranslate is still starting. Then:

```bash
docker-compose restart translator
```

Tests (run in that service directory):

```bash
go test ./...          # GO-Manager
npm test               # web-ui
pytest                 # Translator; needs LibreTranslate
```

## Layout

Clients talk only to GO-Manager. web-ui never calls Translator and never writes PostgreSQL.

- `GO-Manager/` — public Go API (auth, sessions, cards); orchestrates Translator and PostgreSQL
- `Translator/` — Python linguistics for GO-Manager only (not a public client API)
- `web-ui/` — React client of GO-Manager
- `docker-compose.yaml` — manager, translator, web-ui, PostgreSQL, LibreTranslate

web-ui layers, bottom-up (lower never knows upper): API (no React) → Model (no JSX) → Hook → Component (no direct backend).

## Docs

- [GO-Manager structure](GO-Manager/docs/structure.md) — read when changing manager modules, HTTP, sessions, or auth.
- [GO-Manager configuration](GO-Manager/docs/configuration.md) — read when adding or changing manager env vars.
- [Translator structure](Translator/docs/structure.md) — read when changing translator workers, HTTP, or linguistics.
- [Translator env](Translator/docs/Readme.md) — read when changing LibreTranslate env.
- [web-ui structure](web-ui/docs/structure.md) — read when changing UI layers, routes, or the API client.
- [Agent onboarding](.agents/docs/onboarding.md) — read when wiring a new agent or explaining this layout to a human.

## Plans

Save working plans under `.agents/plans/` (gitignored). Do not commit them.
