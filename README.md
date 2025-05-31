# [CRUD] Рецепты печенья (Project Cookies) v 1.0.0

Простой CRUD-сервис для управления рецептами печенья, написанный на Go с использованием PostgreSQL и контейнеризации через Docker (наверное).

## Запуск через docker

Сборка и запуск контейнеров:
```bash
docker-compose -p choco-cookie up --build
```

Остановка и удаление контейнеров:
```bash
docker-compose down
docker rmi -f (ID)
```

Загрузка образа на DockerHub:
```bash
docker push miggyabadocker/cookie-app:latest
```

Получение контейнеров с репозитория DockerHub:
```bash
docker pull miggyabadocker/cookie-app:latest
```

Запуск контейнера:
```bash
docker run -p 8080:8080 miggyabadocker/cookie-app:latest
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
