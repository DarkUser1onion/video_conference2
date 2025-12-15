# 🚀 Быстрый старт

## Запуск проекта на локальной машине

### Windows

1. Откройте PowerShell или CMD в папке проекта
2. Запустите:
```bash
.\start.bat
```

Или вручную:
```bash
docker-compose up --build
```

### Linux/Mac

1. Откройте терминал в папке проекта
2. Сделайте скрипт исполняемым:
```bash
chmod +x start.sh
```

3. Запустите:
```bash
./start.sh
```

Или вручную:
```bash
docker-compose up --build
```

## Первый запуск

1. **Подождите 30-60 секунд** пока все сервисы запустятся
2. Откройте браузер: `http://localhost`
3. Войдите с тестовыми данными:
   - Email: `ilya@example.com`
   - Пароль: `password123`

## Тестирование видеоконференции

1. **В первом браузере:**
   - Войдите в систему
   - Создайте комнату (например, "Тестовая комната")
   - Скопируйте ID комнаты
   - Нажмите "Присоединиться"
   - Разрешите доступ к камере и микрофону

2. **Во втором браузере (или режиме инкогнито):**
   - Откройте `http://localhost`
   - Войдите под другим пользователем (например, `ivan@example.com` / `password123`)
   - Вставьте ID комнаты из первого браузера
   - Нажмите "Присоединиться"
   - Разрешите доступ к камере и микрофону

3. **Проверьте:**
   - Видео обоих участников отображается
   - Аудио работает
   - Можно отправлять сообщения в чат
   - Можно включить/выключить микрофон и камеру
   - Можно начать демонстрацию экрана

## Просмотр логов

```bash
# Все логи
docker-compose logs -f

# Только backend
docker-compose logs -f backend

# Только LiveKit
docker-compose logs -f livekit
```

## Остановка

```bash
docker-compose down
```

## Устранение проблем

### Порты заняты

Если порты 80, 5432, 6379 или 7880 заняты, измените их в `docker-compose.yml` или остановите конфликтующие сервисы.

### Ошибка подключения к LiveKit

- Убедитесь что LiveKit контейнер запущен: `docker-compose ps`
- Проверьте логи: `docker-compose logs livekit`
- Убедитесь что порты 7880 и UDP порты не заняты

### Ошибка подключения к БД

- Подождите еще немного - PostgreSQL может запускаться дольше
- Проверьте логи: `docker-compose logs postgres`

### Фронтенд не загружается

- Проверьте что nginx запущен: `docker-compose ps`
- Проверьте логи: `docker-compose logs nginx`
- Убедитесь что папка `web/` существует с файлом `index.html`

## Полезные команды

```bash
# Пересборка только backend
docker-compose build backend
docker-compose up backend

# Перезапуск всех сервисов
docker-compose restart

# Очистка всех данных (БД, Redis)
docker-compose down -v

# Подключение к PostgreSQL
docker-compose exec postgres psql -U appuser -d app_database

# Подключение к Redis
docker-compose exec redis redis-cli
```

## Тестовые пользователи

Все с паролем `password123`:
- `ilya@example.com` - Илья
- `ivan@example.com` - Иван
- `matvey@example.com` - Матвей
- `artem@example.com` - Артем
- `alesya@example.com` - Алеся

## Что дальше?

- Проект готов к разработке!
- Каждый разработчик может работать над своим модулем
- Фронтенд можно улучшить (React/Vue/Angular)
- Можно добавить больше функций согласно документации

