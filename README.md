# Шаблон golang приложения

Шаблон для начала проекта на golang с использованием chi, fx, cobra и viper в качестве основы.

## roadmap

- [X] Custom Slog [logrus](https://github.com/sirupsen/logrus)
- [X] CLI Command [cobra](https://github.com/spf13/cobra)
- [X] Configuration [viper](https://github.com/spf13/viper)
- [X] Web [chi](https://github.com/go-chi/chi/)
- [X] DI/IOC [fx](https://github.com/uber-go/fx)
- [X] Database postgres
- [X] Query builder [squerel](https://github.com/Masterminds/squirrel)
- [X] Swagger generator [swag](https://github.com/swaggo/swag)
- [X] Migrate [goose](https://github.com/pressly/goose)
- [ ] Seed [пропробую сделаю это](https://pressly.github.io/goose/blog/2021/no-version-migrations/#final-thoughts)
- [ ] Redis
- [ ] Temporal
- [ ] RabbitMQ
- [ ] docker compose файлы для деплоя и локлаьной разработки
- [ ] ......

## слои приложения

```shell
   - app        # application main
     - cmd
     - ... 
   - config              # config
   - deploy              # ci/cd
     - pgsql             # pgsql docker-compose
     - ...               # other     
   - docs                # swag gen swagger2.0 doc
   - internal            # core 
     - controller        # http handler（controller）
     - domain            # domain 
       - user            # домен пользователя 
          - entity       # модели пользователя 
          - repository   # репозитории пользователя 
          - service      # сервисы пользователя 
          - ...          # другие файлы по домену пользователя 
        - ...            # другой домен
   - pkg                 # переиспользуемые пакеты
   - migrations          # миграции
     - ... 
   - ...
```

## Гайд

1. основные команды

   Для запуска приложения запусти или сбилди main.go
   ```
   go run cmd/main.go - запуск приложения
   go run cmd/main.go migrate [аргумент] - запуск миграций goose аргументы в пункте 2
   ```
2. как пользоватся миграциями

```
Commands:
    up                   Выполнить все миграциии
    up-by-one            Выполнить 1 миграцию
    up-to VERSION        Выполнить миграции до определенной версии VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Повторно запустите последнюю миграцию
    reset                Отктить все миграции (опасная операция, лучше не трогать)
    status               Статус миграций
    version              Распечатайте текущую версию базы данных
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Переименовывает миграции согласно порядку
```

3. как пользоватся миграциями
