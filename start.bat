@echo off
chcp 65001 >nul
echo 🚀 Запуск Video Conference Application...
echo.

REM Проверка Docker
where docker >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker не установлен. Установите Docker Desktop и повторите попытку.
    pause
    exit /b 1
)

where docker-compose >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ Docker Compose не установлен. Установите Docker Desktop и повторите попытку.
    pause
    exit /b 1
)

echo ✅ Docker найден
echo.

REM Остановка существующих контейнеров
echo 🛑 Остановка существующих контейнеров...
docker-compose down

echo.
echo 🔨 Сборка и запуск контейнеров...
docker-compose up --build -d

echo.
echo ⏳ Ожидание запуска сервисов...
timeout /t 10 /nobreak >nul

REM Проверка статуса
echo.
echo 📊 Статус сервисов:
docker-compose ps

echo.
echo ✅ Приложение запущено!
echo.
echo 🌐 Откройте в браузере: http://localhost
echo.
echo 📝 Тестовые учетные данные:
echo    Email: ilya@example.com
echo    Пароль: password123
echo.
echo 📋 Просмотр логов: docker-compose logs -f
echo 🛑 Остановка: docker-compose down
echo.
pause

