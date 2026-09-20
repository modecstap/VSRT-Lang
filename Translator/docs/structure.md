# Translator — структура сервиса

Справочник по коду сервиса `Translator`. Переменные окружения LibreTranslate — в [`Readme.md`](Readme.md).

## Назначение

`Translator` отдаёт лингвистические данные по английской фразе или слову. GO-Manager вызывает этот сервис при сохранении карточки.

Сервис:

- переводит en→ru через LibreTranslate (основной вариант и альтернативы);
- собирает синонимы, антонимы и примеры употребления через NLTK WordNet;
- приводит слово к базовой форме (`wordnet.morphy`).

Сервис не хранит пользователей и сессии. Он не является публичным API для web-ui.

Тяжёлые объекты переводчика живут в отдельных процессах. HTTP-слой берёт свободный воркер из пула на время запроса.

## Точка входа

Модуль: корень пакета `Translator`

Реализует: чтение конфига сервера и запуск FastAPI через uvicorn.

Содержит:

- `main.py` — `ServerConfig.from_env()`, создание `FastAPIServer`, `asyncio.run(server.start())`.

## Модули

### `process_pool.py`

Реализует: пул процессов с готовым экземпляром `ITranslator` в каждом воркере. Вызов метода переводчика идёт через очереди `multiprocessing.Queue`.

Содержит:

- `ProcessPool` — старт `pool_size` daemon-процессов, `get()` ждёт свободного воркера, `close()` останавливает воркеры.
- `ProcessPoolProxy` — прокси `translate`, `translate_bulk`, `take_synonyms`, `take_antonyms`, `take_context`, `take_base_form`; контекстный менеджер возвращает воркер в пул.
- исключения `ProcessPoolClosedError`, `ProcessPoolProxyReleasedError`.

Воркер создаёт переводчик через `TranslatorFactory.new_translator()` один раз при старте процесса.

### `server`

Реализует: HTTP-сервер и привязку адреса/порта.

Содержит:

- `fastapi_server.py` — сборка `FastAPI`, подключение роутера `/api/translator`, запуск uvicorn на `TRANS_HOST`:`TRANS_PORT`.
- `server_config.py` — `ServerConfig.from_env()`: `TRANS_HOST` (по умолчанию `0.0.0.0`), `TRANS_PORT` (по умолчанию `8080`).

### `server/routers`

Реализует: объявление маршрутов Translator API.

Содержит:

- `translator_router.py` — `TranslatorRouter`:
  - `POST /translate`
  - `POST /translate-bulk`
  - `GET /synonyms`
  - `GET /antonyms`
  - `GET /base-form`
  - `GET /context`

Префикс роутера: `/api/translator`.

### `server/handlers`

Реализует: обработку HTTP-запросов: взять воркер из пула, вызвать метод переводчика, вернуть JSON-совместимый результат.

Содержит:

- `translation_handler.py` — `TranslationHandler`:
  - создаёт `ProcessPool(TranslatorFactory(), TRANSLATOR_COUNT)` (по умолчанию 1 воркер);
  - `translate` / `translate_bulk` возвращают `Translation`;
  - `take_synonyms` / `take_antonyms` возвращают список строк;
  - `take_context` возвращает список `Translation`;
  - `take_base_form` возвращает строку.

### `translator`

Реализует: контракт и рабочую реализацию переводчика.

Содержит:

- `i_translator.py` — абстрактный `ITranslator` и базовое исключение `TranslatorError`.
- `translator.py` — `Translator`: LibreTranslate (`/translate`, `source=en`, `target=ru`) и WordNet (синонимы, антонимы, примеры, `morphy`). При создании объекта выполняется пробный `translate("initial")`.
- `translator_factory.py` — `TranslatorFactory.new_translator()`: `LIBRETRANSLATE_URL`, `LIBRETRANSLATE_API_KEY`; повторные попытки каждые 50 секунд, пока LibreTranslate не ответит.

### `translator/models`

Реализует: pydantic-модели ответа.

Содержит:

- `word.py` — `Word` с полем `content`.
- `tlanslation.py` — `Translation` с полями `original` и `translations`.

## Связи

```
HTTP /api/translator/*
  → TranslatorRouter
      → TranslationHandler
          → ProcessPool.get()
              → воркер-процесс
                  → TranslatorFactory.new_translator()
                      → Translator
                          → LibreTranslate (перевод)
                          → NLTK WordNet (синонимы, антонимы, контекст, базовая форма)
```

GO-Manager вызывает те же пути через `internal/net_translator`.
