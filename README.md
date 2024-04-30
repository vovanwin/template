# Шаблон golang приложения

Шаблон для начала проекта на golang с использованием chi, fx, cobra и viper в качестве основы.

## roadmap

- [X] logger Slog
- [X] CLI Command [cobra](https://github.com/spf13/cobra)
- [X] Configuration [cleanEnv](https://github.com/ilyakaznacheev/cleanenv)
- [X] Web [chi](https://github.com/go-chi/chi/)
- [X] DI/IOC [fx](https://github.com/uber-go/fx)
- [X] Database postgres
- [X] ent ORM [ent](https://entgo.io/)
- [ ] Swagger codegen [oapi-codegen](https://github.com/deepmap/oapi-codegen)
- [X] Migrate [atlas](https://atlasgo.io/integrations/go-sdk)
- [ ] Seed  
- [ ] Redis
- [ ] Temporal
- [ ] RabbitMQ
- [x] docker compose файлы для локлаьной разработки
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
        - store          # ent ORM сгенерирование файлы и настройки, будут использоваться во всех репозиториях
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
   go run app/main.go migration  - запуск миграций atlas
    ```
      

2. как пользоватся миграциями

Миграции основаны на инструенте Atlas https://atlasgo.io/integrations/go-sdk

В task добавлены необходимые команды
```
task migrate:lint - запускает линтер миграций, проверяет что миграции накатятся без проблем, если такие есть то выведет предупреждения

task migrate:orm - генерирует из ent scheme (Модели) миграции сравнивая разницу между существующими миграциями и текущими моделями
   
task migrate:create - создать новый файл миграции
 
task migrate:hash - перехешировать хеш сумму миграций если что либо поменял
   
```

3. Swagger документация - В ПРОЦЕССЕ
   
Для генерации swagger документации по анотациям в контроллере
```
swag init -o docs -g app/main.go
```