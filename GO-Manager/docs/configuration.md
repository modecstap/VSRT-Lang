
# ENV Variable

NAME            DESCRIPTION                 MANDATORY   VALUES
___
SECRET_KEY      key to generate JWT Token   Yes
MANAGER_HOST    host of current app         Yes
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
