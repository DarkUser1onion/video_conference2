# Video Conference Application

Система видеоконференций на Go с использованием LiveKit для медиа-потоков.

## Структура проекта

```
video_conference/
├── cmd/
│   └── server/          # Точка входа приложения
├── internal/
│   ├── config/         # Конфигурация
│   ├── domain/         # Доменные модели
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # Middleware (auth, CORS, rate limit)
│   ├── repository/     # Репозитории для работы с БД
│   └── service/        # Бизнес-логика
├── pkg/
│   ├── jwt/            # JWT утилиты
│   └── logger/         # Логирование
├── migrations/         # Миграции БД (TODO)
├── docs/               # Документация
├── docker-compose.yml  # Docker Compose конфигурация
└── go.mod              # Go модули
```

## Технологический стек

- **Backend**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Media Server**: LiveKit
- **WebSocket**: gorilla/websocket
- **Authentication**: JWT (golang-jwt/jwt)

## Установка и запуск

1. Установите зависимости:
```bash
go mod download
```

2. Настройте переменные окружения (создайте `.env` файл):
```env
ENVIRONMENT=development
SERVER_PORT=8080
DATABASE_DSN=postgres://appuser:apppass123@localhost:5432/app_database?sslmode=disable
REDIS_ADDR=localhost:6379
JWT_ACCESS_SECRET=your-access-secret-key-change-in-production
JWT_REFRESH_SECRET=your-refresh-secret-key-change-in-production
LIVEKIT_URL=ws://localhost:7880
LIVEKIT_API_KEY=your-livekit-api-key
LIVEKIT_API_SECRET=your-livekit-api-secret
```

3. Запустите PostgreSQL и Redis через Docker:
```bash
docker-compose up -d
```

4. Запустите приложение:
```bash
make run
# или
go run cmd/server/main.go
```

## API Endpoints

### Аутентификация
- `POST /api/v1/auth/register` - Регистрация
- `POST /api/v1/auth/login` - Вход
- `POST /api/v1/auth/refresh` - Обновление токена

### Пользователи
- `GET /api/v1/users/me` - Получить текущего пользователя
- `PUT /api/v1/users/me` - Обновить профиль
- `GET /api/v1/users/me/settings` - Получить настройки
- `PUT /api/v1/users/me/settings` - Обновить настройки

### Комнаты
- `POST /api/v1/rooms` - Создать комнату
- `GET /api/v1/rooms` - Список комнат
- `GET /api/v1/rooms/:id` - Получить комнату
- `PUT /api/v1/rooms/:id` - Обновить комнату
- `DELETE /api/v1/rooms/:id` - Удалить комнату
- `POST /api/v1/rooms/:id/join` - Присоединиться к комнате
- `POST /api/v1/rooms/:id/leave` - Покинуть комнату
- `POST /api/v1/rooms/:id/invite` - Создать приглашение
- `GET /api/v1/rooms/:id/participants` - Список участников

### Чат
- `GET /api/v1/rooms/:roomId/chat/messages` - Получить сообщения
- `POST /api/v1/rooms/:roomId/chat/messages` - Отправить сообщение
- `PUT /api/v1/rooms/:roomId/chat/messages/:messageId` - Редактировать сообщение
- `DELETE /api/v1/rooms/:roomId/chat/messages/:messageId` - Удалить сообщение

### Медиа
- `POST /api/v1/rooms/:roomId/media/token` - Получить LiveKit токен

### Статистика
- `GET /api/v1/rooms/:roomId/stats` - Статистика комнаты
- `GET /api/v1/rooms/:roomId/stats/participants/:participantId` - Статистика участника

## Разработка

### Структура модулей по команде

- **Артем**: Аутентификация (JWT), Frontend основной комнаты
- **Алеся**: Frontend регистрации, Frontend главного меню
- **Илья и Иван**: Backend создания комнат, Backend основной комнаты, Backend демонстрации экрана и чата
- **Матвей**: Демонстрация экрана и чат (фиолетовый на диаграмме)

## TODO

- [ ] Реализовать миграции БД
- [ ] Добавить полную реализацию Waiting Room
- [ ] Реализовать WebSocket для чата
- [ ] Добавить интеграционные тесты
- [ ] Настроить CI/CD
- [ ] Добавить документацию API (Swagger/OpenAPI)
