# Сервис PR-Manager

## Описание
PR-Manager — это сервис для управления пулл-реквестами и командами, с возможностью отслеживания ревью, статуса PR и активности пользователей.
## Использовавшийся стек
Go, Postgres, Docker, Postman

## Использовавшиеся технологии
Clean Architecture, Makefile, REST API, Unit Tests, Graceful Shutdown, Backoff

## Использовавшиеся библиотеки и фреймворки
chi, pgx, gomock

## Сборка и запуск
Существует 2 способа запуска приложения:
1. Локально: через docker-compose.local.yaml. Через докер подтягивается бд и миграции, а сервис запускается локально через make
```bash
make up_local
```

1. Локально, но по-другому: через docker-compose.prod.yaml. Все то же самое, сам сервис подтягивается с моего DockerHub.
```bash
make up_prod
```
Тесты запускаются следующей командой:
```bash
make run_tests
```

## Слои
### Репозиторий
Имеется три объекта репозитория (user, team, pull_request). Со следующей зависимостью
```
userRepository          -> pgxpool
teamRepository          -> pgxpool + userRepository
pullRequestRepository   -> pgxpool + teamRepository + userRepository
```

### Сервис
Так же три объекта
```
userService          -> userRepository
teamService          -> teamRepository
pullRequestService   -> pullRequestRepository
```

### Контроллеры
Тут тоже без сюрпризов
```
userHandler          -> userService
teamHandler          -> teamService
pullRequestHandler   -> pullRequestService
```

## Эндпоинты
Все эндпоинты соответствуют OpenAPI документации описанной в openapi.yaml
Повзаимодействовать с ними можно через Postman:
```
https://web.postman.co/workspace/My-Workspace~d53d97d9-99e6-48a8-8c49-55ac2dc58ca5/collection/36633954-e4986fb6-ba41-4093-8ef5-787858606f05
```

## Переменные окружения
Пример хранится в .env в корневой папке проекта.