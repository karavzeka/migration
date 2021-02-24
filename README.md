# Миграции

В основе данного проекта лежит эта библиотека, написанная на Go https://github.com/golang-migrate/migrate. Запуск команд осуществляется через `make`, который запускает скрипты внутри docker контейнеров. Так делается, чтобы отвязаться от окружения.

Для локальных тестов можно запустить контейнеры с postgres и mysql:
```
docker-compose up -d
```
В случае возникновения ошибки о несуществующей сети запустить команду:
```
make network

# или это же можно сделать через сам докер:
docker network create migration-network
```

### Make команды
[create](src/cmd/create/README.md) - создание up/down файлов с новой версией

[init](src/cmd/init/README.md) - инициализация первой миграции на основе существующей БД

[migrate](src/cmd/migrate/README.md) - применение или откатывание миграций

[snapshot](src/cmd/snapshot/README.md) - создание снимка БД в отдельную директорию со снимками

[version](src/cmd/version/README.md) - отображает текущую версию, или принудительно прописывает в БД указанную версию

### Требования к окружению
Каждая команда принимает на вход параметр `<db_alias>`, некий алиас для БД, он может быть любым, не обязан совпадать с именем БД. При запуске команд в проекте предварительно необходимо создать директорию `databases/<db_alias>`.

Для большинства команд необходимы переменные окружения, содержащие параметры коннекта к БД:
```
DB_ENGINE    postgres|mysql
DB_USER      user name
DB_PASSWORD  password
DB_HOST      host
DB_PORT      port
DB_DATABASE  database name
```
Но проще создать файл `databases/<db_alias>/.env` с этими параметрами. См. пример [.env-example](.env-example)

### IF [NOT] EXISTS
Зачастую при создании/удалении таблиц и других сущностей в БД удобно использовать конструкцию `IF EXISTS`/`IF NOT EXISTS`, чтобы не делать лишних действия для проверки существования/отсутствия сущности. Однако в миграциях лучше отказаться от использования данных конструкций, ибо они могут нарушить консистентность самих миграций. Может например случиться такое, что объявление создания таблицы повторится в нескольких миграциях, что логически неверно, при этом не возникнет никаких ошибок (или возникнут не там, где ожидалось).