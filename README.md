# Шаблон golang приложения

Шаблон для начала проекта на golang с использованием chi, fx, cobra и sqlc в качестве основы.

## roadmap

- [X] logger Slog
- [X] CLI Command [cobra](https://github.com/spf13/cobra)
- [X] Configuration [cleanEnv](https://github.com/ilyakaznacheev/cleanenv)
- [X] Web [chi](https://github.com/go-chi/chi/)
- [X] DI/IOC [fx](https://github.com/uber-go/fx)
- [X] Database postgres
- [X] sqlc ORM [sqlc](https://docs.sqlc.dev/en/latest/)
- [x] codegen [ogen](https://github.com/ogen-go/ogen)
- [X] Migrate [goose](https://github.com/pressly/goose)
- [ ] Seed  
- [ ] Redis
- [ ] Temporal
- [ ] RabbitMQ
- [x] docker compose файлы для локлальной разработки
- [ ] docker compose файлы для прода
- [ ] ......

## слои приложения

```shell
   - app        # application main
     - cmd
     - config            # config
     - database
      - migrations          # миграции
     - internal          # core 
      - module            # domain 
       - shared          # общие файлы
        - types          # типы данных используемых по всему приложению, сейчас тут uuid сгенерирование из cmd/gen-types
       - user            # Модуль пользователя 
          - controller   # папка с контроллерами. 
          - repository   # репозитории пользователя 
          - service      # сервисы пользователя 
          - ...          # другие файлы по модулю пользователя 
        - ...            # другой домен
     - pkg                 # переиспользуемые пакеты     - ... 
   - deployments              # ci/cd
     - local             # docker-compose
     - ...               # other     
   - docs                # openapi для кодогенерации контроллеров
     - ... 
   - ...
```

## Гайд

1. основные команды

   Для запуска приложения запусти или сбилди main.go
   ```
   go run app/main.go - запуск приложения
   go run app/main.go migration up  - запуск миграций goose
    ```
      

2. как пользоватся миграциями

 

3. Swagger документация - В ПРОЦЕССЕ
   
