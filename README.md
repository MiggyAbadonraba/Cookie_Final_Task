# [CRUD] Рецепты печенья (Project Cookies)

Простой CRUD-сервис для управления рецептами печенья, написанный на Go с использованием PostgreSQL и контейнеризации через Docker (наверное).

## Запуск через docker

Сборка и запуск контейнеров:
```bash
docker-compose -p choco-cookie up --build
```

Остановка и удаление контейнеров:
```bash
docker-compose down
```

Приложение будет доступно по адресу:
```http://localhost:8080```

## Сборка и тестирование

Сборка проекта локально:
```terminal
go build -o cookies
```

Запуск модульных тестов:
```terminal
go test ./...
```