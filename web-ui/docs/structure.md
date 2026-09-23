# Web-UI — структура сервиса

Справочник по коду React-приложения `web-ui`. Стили (CSS-модули) в этот документ не входят.

## Назначение

`web-ui` — браузерный клиент VSRT-LANG. Пользователь входит в систему, смотрит список сессий и наполняет активную сессию словами.

Приложение:

- вызывает HTTP API GO-Manager;
- хранит JWT и id активной сессии в `localStorage`;
- разделяет экран на слои API → Model → Hook → Component.

Приложение не ходит в Translator напрямую и не пишет в PostgreSQL.

## Точка входа

Модуль: `src`

Реализует: монтирование React-дерева, глобальный SWR и маршруты.

Содержит:

- `index.js` — `ReactDOM.createRoot`, `SWRConfig` (`revalidateOnFocus`, `dedupingInterval: 2000`), `RouterProvider`.
- `router.jsx` — `createBrowserRouter`:
  - `/` — страница входа; при наличии access-токена редирект на `/account`;
  - `/account` — требуется токен; вкладки `profile` и `session`;
  - страницы грузятся через `React.lazy` / `Suspense`.
- `routes/prefetch.js` — `prefetchAccount`, `prefetchSessionTab`: предварительная загрузка чанков кабинета.

## Слои страницы

Правило зависимостей (снизу вверх ничего не знает о верхнем слое):

1. **API** — HTTP и хранилище. Нет React и JSX.
2. **Model** — преобразование данных и ключи кэша. Нет JSX.
3. **Hook** — состояние экрана и обработчики для UI.
4. **Component** — разметка. Не вызывает backend напрямую.

## Модули

### `src/api`

Реализует: общий HTTP-клиент, токены и операции с сессиями.

Содержит:

- `client.js` — axios с `REACT_APP_API_URL` (по умолчанию `http://localhost:8080`); подстановка `Authorization: Bearer`; на `401` вне `/login` и `/register` очищает сессию и отправляет на `/`.
- `auth.js` — чтение/запись access и refresh токенов, `hasAccessToken`.
- `storage.js` — ключи `vsrt.auth.v1` и `vsrt.session.v1`; миграция со старых ключей `access_token`, `refresh_token`, `session_tab_session_id`.
- `sessions.js` — `GET /users/sessions`, создание и удаление сессии, активный `sessionId`, `GET/POST /sessions/{id}/records`, ключи SWR.

### `src/pages/LoginPage`

Реализует: экран входа и регистрации.

Содержит:

- `LoginPage.jsx` — переключение форм LOGIN / REGISTER.

#### `src/pages/LoginPage/api`

Реализует: запросы аутентификации.

Содержит:

- `login.js` — `POST /login` (`login`, `password`).
- `register.js` — `POST /register` (`username`, `email`, `password`).

#### `src/pages/LoginPage/model`

Реализует: начальное состояние форм.

Содержит:

- `loginFormModel.js` — `{ login, password }`.
- `registerFormModel.js` — `{ username, email, password }`.

#### `src/pages/LoginPage/hooks`

Реализует: состояние форм, отправку, ошибки, навигацию после успеха.

Содержит:

- `useLoginForm.js` — вход, `saveAuthTokens`, prefetch кабинета, переход на `/account`.
- `useRegisterForm.js` — регистрация; при успехе вызывает `onSuccess` (возврат к форме входа).

#### `src/pages/LoginPage/components`

Реализует: разметку форм.

Содержит:

- `LoginForm/LoginForm.jsx` — поля login/password, кнопка входа, ссылка на регистрацию.
- `LoginForm/FormField.jsx` — обёртка над `shared/ui/Input`.
- `RegisterForm/RegisterForm.jsx` — username/email/password, возврат к входу.

### `src/pages/AccountPage`

Реализует: каркас кабинета после входа.

Содержит:

- `AccountPage.jsx` — меню PROFILE / SESSION / LOGOUT, `Outlet` для вкладок; logout чистит токены, активную сессию и кэш SWR.

### `src/pages/AccountPage/components/ProfileTab`

Реализует: профиль и список учебных сессий.

Содержит:

- `ProfileTab.jsx` — карточка пользователя и таблица сессий. Карточка берёт имя, почту и аватар из `useProfileCard`; загрузка аватара тоже в этом хуке. Список сессий по-прежнему из `useProfileTab`.

#### `ProfileTab/api`

Реализует: реэкспорт операций сессий из `src/api/sessions.js`.

Содержит:

- `profileTabApi.js` — `fetchUserSessions`, `createSession`, `deleteSession`, `setActiveSessionId`, ключи SWR.
- `getCurrentUser.js` — `GET /users/me`, ключ SWR `['current-user']`.
- `saveAvatar.js` — `POST /users/avatar`, multipart-поле `avatar`.

#### `ProfileTab/model`

Реализует: представление сессии для таблицы.

Содержит:

- `profileTabModel.js` — `mapSessionToViewModel` / `buildSessionsViewModel`: `id`, имя, число сохранённых записей, дата (если API её отдаёт).
- `profileCardModel.js` — `profileCardView` собирает загрузку, ошибку без тела, успешное тело и ошибку повтора при уже полученном теле. `avatarUploadError` переводит код ошибки сохранения аватара в строку карточки.

#### `ProfileTab/hooks`

Реализует: загрузку списка, создание, выбор и удаление сессии.

Содержит:

- `useProfileTab.js` — SWR по `SESSIONS_SWR_KEY`; создание сессии с переходом на `/account/session`; выбор сессии записывает активный id; удаление обновляет кэш.
- `useProfileCard.js` — SWR по `['current-user']`; загрузка аватара; наружу `uploading`, `uploadError`, `handleAvatarChange` и поля карточки. После HTTP 204 обновляет тот же ключ. `uploadError` сам очищается через 5 секунд.

#### `ProfileTab/components`

Реализует: визуальные блоки вкладки.

Содержит:

- `ProfileUserCard.jsx` — аватар: иконка загрузки по hover и focus, клик открывает выбор файла. Ошибка загрузки на 5 секунд заменяет username. Кнопка CHANGE без обработчика.
- `SessionList.jsx` — поле имени, кнопка NEW, таблица сессий, удаление без перехода в сессию.

### `src/pages/AccountPage/components/SessionTab`

Реализует: работа с карточками активной сессии: список слов, детали, ввод слова и контекста.

Содержит:

- `SessionTab.jsx` — список сохранённых слов, перевод/синонимы/антонимы, значение и контексты, поля Word/Context, кнопка WRITE.

#### `SessionTab/api`

Реализует: загрузку и сохранение карточек текущей сессии.

Содержит:

- `sessionTabApi.js` — `loadSavedWords` (`getOrCreateSessionId` + `fetchSessionRecords`); `saveWordEntry` (`phrase`, `context`).

#### `SessionTab/model`

Реализует: нормализацию записей API к модели слова на экране.

Содержит:

- `sessionTabModel.js` — `mapRecordToWord`, `buildWordMap`, `mergeSavedWord`, `buildWordDetails` (строка переводов, списки синонимов/антонимов/контекстов, базовая форма как «значение»).

#### `SessionTab/hooks`

Реализует: кэш записей, выбор слова, отправку новой карточки.

Содержит:

- `useSessionTab.js` — SWR по `sessionRecordsKey(activeSessionId)`; синхронизация поля ввода с выбранным словом; `WRITE` вызывает API и обновляет список.

### `src/shared/ui`

Реализует: переиспользуемые элементы без знания страниц и API.

Содержит:

- `Button/Button.jsx` — кнопка, вариант `primary` по умолчанию.
- `Input/Input.jsx` — поле или textarea (`multiline`).

## Связи

```
браузер
  → router.jsx (токен? LoginPage : AccountPage)
      LoginPage
        → useLoginForm / useRegisterForm
            → pages/LoginPage/api
                → api/client.js → GO-Manager /login, /register
      AccountPage
        → ProfileTab → useProfileTab → api/sessions.js → /users/sessions, /sessions
        → ProfileTab → useProfileCard → getCurrentUser.js → GET /users/me
        → ProfileTab → useProfileCard → saveAvatar.js → POST /users/avatar
        → SessionTab → useSessionTab → sessionTabApi → /sessions/{id}/records
```

Активная сессия: `ProfileTab` пишет id в `localStorage`. `SessionTab` читает его; если id нет, `getOrCreateSessionId` создаёт сессию на backend.
