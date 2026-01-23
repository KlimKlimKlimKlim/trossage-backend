# Trossage Backend

Backend для мессенджера с поддержкой real-time коммуникации через WebSocket.

## Возможности

- Регистрация и аутентификация (JWT access/refresh tokens)
- Приватные чаты между пользователями
- Real-time доставка сообщений и событий (WebSocket)
- Индикатор набора текста (typing indicator)
- Поиск пользователей

## Технологии

- **Go 1.25**, Gin, pgx
- **WebSocket** (coder/websocket)
- **JWT** для аутентификации
- **Argon2id** для хеширования паролей
- **PostgreSQL 16**

## Запуск

```bash
# Скопировать и настроить переменные окружения
cp .env.example .env

# Запуск через Docker Compose
docker-compose up -d
```

API будет доступен на `http://localhost:8080`, health-check на `http://localhost:8081/health`.

## API

Документация в формате Swagger находится в директории [`docs/swagger.json`](docs/swagger.json).

### Основные эндпоинты

| Метод | Путь                      | Описание             |
|-------|---------------------------|----------------------|
| POST  | `/api/auth/register`      | Регистрация          |
| POST  | `/api/auth/login`         | Вход                 |
| POST  | `/api/auth/refresh`       | Обновление токенов   |
| POST  | `/api/auth/logout`        | Выход                |
| GET   | `/api/users/me`           | Текущий пользователь |
| GET   | `/api/users/search`       | Поиск пользователей  |
| POST  | `/api/chats`              | Создать чат          |
| GET   | `/api/chats`              | Список чатов         |
| POST  | `/api/chats/:id/messages` | Отправить сообщение  |
| GET   | `/api/chats/:id/messages` | История сообщений    |
| GET   | `/api/ws`                 | WebSocket соединение |

## Структура проекта

```
cmd/                 # Точка входа
internal/
  application/       # Жизненный цикл приложения
  config/            # Конфигурация
  http/              # HTTP handlers, middleware, DTO
  postgres/          # Repository layer
  service/           # Бизнес-логика
  websocket/         # WebSocket hub и клиенты
  worker/            # Фоновые задачи
migrations/          # SQL миграции
docs/                # Swagger документация
```
