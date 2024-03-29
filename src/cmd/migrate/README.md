# Применение миграций
```
make migrate <db_alias> <command> [<version>]
```
Команда `migrate` вызывает определенную команду миграции, например `up`. Параметры коннекта к БД берутся из файла `databases/<db_alias>/.env`.

`<db_alias>` - алиас базы данных для которой применяется миграция.  
`<version>` - версия, обязательна или нет - зависит от команды.  
`<command>`:
```
goto V       Migrate to version V
up [V]       Apply all or up to V version
down [V]     Apply all or up to V version
```

Если во время миграций произошла ошибка, например ошибка в sql при создании таблицы, то несмотря на то, что таблица не создалась, в БД пропишется версия сломанной миграции, но при этом выставится пометка, что эта миграция грязная. В таком случае надо руками поправить БД (см. [version](../../cmd/version/README.md)).