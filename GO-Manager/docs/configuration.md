
# ENV Variable

NAME            DESCRIPTION                 MANDATORY   VALUES
___
SECRET_KEY      key to generate JWT Token   Yes
MANAGER_HOST    listen address host:port    Yes
DB_HOST         database address            Yes
DB_PORT         database port               Yes
DB_NAME         database name               Yes
DB_USER         database user               Yes
DB_PASS         database password           Yes
TRANSLATOR_HOST host of translator          Yes
MODE                                        NO          prod, debug
SMTP_HOST       SMTP server host            NO
SMTP_PORT       SMTP server port            NO
SMTP_USER       SMTP login                  NO
SMTP_PASS       SMTP password               NO
SMTP_FROM       address in the From header  NO
PUBLIC_SITE_URL public site origin, no path NO

Without `SMTP_HOST`, `SMTP_PORT`, `SMTP_FROM` or `PUBLIC_SITE_URL` a password-reset email request responds with `Unable to send email` and does not start the 15-minute send limit.

`MANAGER_HOST` — аргумент `http.ListenAndServe` (`host:port`). Это адрес, на котором слушает процесс, не публичный URL сайта. `localhost:8080` слушает только loopback этого процесса: другой контейнер и опубликованный порт хоста до него не подключаются. В `docker-compose.yaml` у сервиса `manager` задано `0.0.0.0:8080`. Значение `localhost:8080` в `.env` остаётся для запуска процесса на хосте; compose его в сервис `manager` не подставляет.

`PUBLIC_SITE_URL` — origin в ссылке письма сброса пароля, без path. Для UI, открытого с другой машины, это `http://<server-host>:<ui-port>`. К `MANAGER_HOST` и к base URL браузера не относится. Без него запрос письма по-прежнему отвечает `Unable to send email`.
