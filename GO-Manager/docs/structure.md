# GO-Manager — структура сервиса

Справочник по коду backend-сервиса `GO-Manager`. Документ описывает назначение сервиса и состав модулей. Контракт HTTP API — в Swagger (`/swagger/`). Переменные окружения — в [`configuration.md`](configuration.md).

## Назначение

`GO-Manager` — публичный backend VSRT-LANG. Клиенты ходят только сюда.

Сервис:

- регистрирует пользователей и выдаёт JWT;
- хранит учебные сессии и карточки слов;
- при сохранении новой фразы запрашивает перевод, синонимы, антонимы, базовую форму и контексты у сервиса Translator;
- оркестрирует доступ к PostgreSQL.

Сервис не переводит слова сам и не рисует UI.

В режиме `MODE=debug` вместо Translator подключается заглушка. В `prod` нужен `TRANSLATOR_HOST`.

## Точка входа

Модуль: `cmd`

Реализует: запуск процесса, подключение к PostgreSQL, миграции, сборку зависимостей и HTTP-сервера.

Содержит:

- `main.go` — JSON-логгер, цикл подключения к БД, `migrations.Run`, `DependensFromEnv`, регистрацию handler'ов, CORS, лог запросов, `ListenAndServe` на `MANAGER_HOST`.

## Модули

### `internal/user`

Реализует: доменную модель пользователя, контракт хранилища, хеширование пароля и валидацию email/пароля.

Содержит:

- `user.go` — `User`, `UserId`, `Repository` (`Create`, `FindByEmail`, `FindByUsername`, `FindByID`), `NewUser`, bcrypt-хеш и сравнение пароля.
- `validate.go` — `ValidateEmail`, `ValidatePassword` (минимум 8 символов, верхний и нижний регистр, цифра).

### `internal/auth`

Реализует: регистрацию, вход, пару access/refresh токенов, отзыв refresh-токена.

Содержит:

- `service.go` — `Register`, `Login`, `Refresh`, `Logout`; поиск пользователя по username, затем по email; SHA-256 хеш refresh-токена перед записью в БД.
- `jwt.go` — HS256 access-токен (TTL 2 часа), разбор claims (`sub`, email, username), генерация refresh-токена (TTL 30 суток).
- `repository.go` — `RefreshToken` и интерфейс `RefreshTokenRepository` (`Save`, `FindByHash`, `RevokeByHash`, `DeleteByHash`).
- `password.go` — локальные хелперы bcrypt и валидации. Регистрация и вход вызывают функции пакета `internal/user`.

HTTP-маршруты `/refresh` и `/logout` в роутере не зарегистрированы. Методы `Refresh` и `Logout` есть только в сервисе.

### `internal/session`

Реализует: учебную сессию, карточку слова и прикладные операции над ними. Перевод делегирует интерфейсу `Translator`.

Содержит:

- `session.go` — структура `Session` (`ID`, `User`, `Name`, карта `Records` по фразе); `NewSession`; `SaveRecord` (новая карточка или увеличение `Count`); `GetRecords`.
- `record.go` — `Record` и `Context`; `NewRecord` собирает переводы, контексты, базовую форму, синонимы и антонимы через `Translator`.
- `service.go` — `NewSession`, `GetSessions`, `GetSession`, `DeleteSession`, `AddRecord`, `DeleteRecord`; проверка владельца сессии (`ErrUnauthorized`, `ErrSessionNotFound`).
- `repository.go` — интерфейс хранилища сессий: `Save`, `Take`, `FindByUser`, `Delete`.
- `translator.go` — интерфейс лингвистического бэкенда: `Translate`, `TranslateBulk`, `TakeContexts`, `TakeSynonyms`, `TakeAntonyms`, `TakeBaseForm`.

### `internal/http`

Реализует: маршрутизацию и сборку зависимостей сервера из окружения.

Содержит:

- `router.go` — `NewServeMux`: публичные `POST /register`, `POST /login`; защищённые `POST /sessions`, `DELETE /sessions/`, `POST /sessions/`, `GET /sessions/`, `GET /users/`; раздача Swagger.
- `server_deps.go` — `DependensFromEnv`: postgres-репозитории, `auth.Service`, `session.Service`, JWT; выбор Translator по `MODE`.

### `internal/http/handlers`

Реализует: единый JSON-формат ошибок HTTP.

Содержит:

- `error.go` — `WriteError` с полями `code`, `message`, `details`.

### `internal/http/handlers/auth`

Реализует: HTTP-слой регистрации и входа.

Содержит:

- `auth.go` — конструктор `Handler` над `auth.Service`.
- `dto.go` — `RegisterRequest`, `LoginRequest`.
- `register.go` — `POST /register`: создаёт пользователя, отвечает `{id, email}`.
- `login.go` — `POST /login`: возвращает `AccessToken` и `RefreshToken`.

### `internal/http/handlers/session`

Реализует: HTTP-слой учебных сессий и записей. Пользователь берётся из JWT-контекста.

Содержит:

- `session.go` — `Handler` и локальный интерфейс `Service`.
- `create.go` — создание сессии, ответ `201` с `id` и `name`.
- `save_record.go` — разбор пути `/sessions/{id}/records`, вызов `AddRecord`.
- `get_records.go` — список записей сессии.
- `delete.go` — удаление сессии, `204` при успехе.

### `internal/http/handlers/user`

Реализует: HTTP-слой списка сессий текущего пользователя.

Содержит:

- `handler.go` — `Handler` с `user.Repository` и репозиторием сессий.
- `get_sessions.go` — `GET /users/sessions`: сессии по `user_id` из токена.

### `internal/http/middleware`

Реализует: сквозную обработку HTTP-запросов.

Содержит:

- `auth.go` — проверка `Authorization: Bearer …`, запись `user_id` и email в context, `UserFromContext`.
- `cors.go` — заголовки CORS и ответ на preflight `OPTIONS`.
- `request_logger.go` — лог метода, пути, статуса, длительности; тело JSON с маскированием `password` и токенов.

### `internal/database/postgres`

Реализует: чтение параметров PostgreSQL из окружения.

Содержит:

- `config.go` — `LoadDBConfig` (`DB_HOST`, `DB_PORT`, `DB_NAME`, `DB_USER`, `DB_PASS`) и строка подключения без SSL.

### `internal/database/postgres/migrations`

Реализует: последовательное применение SQL-миграций и таблицу `schema_migrations`.

Содержит:

- `migrations.go` — `Run`: пользователи, сессии, записи, refresh-токены; identity для `sessions.id`; колонка `count`; уникальность `(session_id, phrase)`.

### `internal/database/postgres/user_repository`

Реализует: postgres-реализацию `user.Repository`.

Содержит:

- `repository.go` — обёртка над `*sql.DB`.
- `create.go` — вставка пользователя.
- `find_by_id.go` — поиск по id.
- `find_by_email.go` — поиск по email.
- `find_by_username.go` — поиск по username.

### `internal/database/postgres/session_repository`

Реализует: postgres-реализацию `session.Repository`. Сохраняет сессию и записи в одной транзакции.

Содержит:

- `repository.go` — конструктор.
- `models.go` — JSON-контексты и загрузка записей сессии (`applyRecords`).
- `save.go` — upsert сессии, полная перезапись `records`.
- `take.go` — сессия по id вместе с записями.
- `get_by_user.go` — все сессии пользователя.
- `delete.go` — удаление сессии (записи удаляются каскадом).

### `internal/database/postgres/refresh_token_repository`

Реализует: postgres-реализацию `auth.RefreshTokenRepository`.

Содержит:

- `repository.go` — сохранение по `token_hash`, поиск, отзыв (`revoked = true`), удаление.

### `internal/database/memory`

Реализует: in-memory репозитории для тестов и локальной разработки. Продакшен-сборка из `cmd/main.go` их не подключает.

Содержит:

- `user_repository.go` — пользователи в картах по id/email/username.
- `session.go` — сессии в памяти с клонированием записей.
- `token_repository.go` — refresh-токены по хешу.

### `internal/net_translator`

Реализует: HTTP-клиент к Translator. Базовый URL: `{TRANSLATOR_HOST}/api/translator`, таймаут 30 с.

Содержит:

- `translator.go` — нормализация хоста и клиент.
- `translate.go` — `POST /translate`.
- `translate_bulk.go` — `POST /translate-bulk`.
- `take-context.go` — `GET /context`.
- `take_synonyms.go` — `GET /synonyms`.
- `take_antonyms.go` — `GET /antonyms`.
- `take_base_form.go` — `GET /base-form`.

### `internal/translators/stub`

Реализует: детерминированную заглушку `session.Translator` для `MODE=debug`.

Содержит:

- `translator.go` — фиксированные переводы, синонимы, антонимы, контекст и базовая форма без сети.

### `internal`

Реализует: дополнительную тестовую заглушку Translator.

Содержит:

- `mock_translator.go` — `MockTranslator` с русскими фиктивными данными. Точка входа сервиса его не использует.

## Связи

```
клиент
  → CORS / RequestLogger / ServeMux
      → handlers/auth → auth.Service → user.Repository + RefreshTokenRepository + JWTService
      → middleware.Auth → handlers/session → session.Service
            → session.Repository (postgres)
            → Translator (net_translator | stub)
      → handlers/user → session.Repository.FindByUser
```

Создание карточки: `POST /sessions/{id}/records` → `AddRecord` → `Session.SaveRecord` → `NewRecord` → Translator → `Repository.Save`.
