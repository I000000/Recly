# Recly — Cross‑media Recommendation System

Recly — это кросс-медийная рекомендательная система, которая рекомендует фильмы на основе любимых книг пользователя и наоборот, а также выдаёт внутридоменные рекомендации. Система использует мультимодальный подход, объединяя текстовые описания, жанровые метки и визуальные признаки (изображения обложек и постеров) для построения эмбеддингов элементов каталога. Рекомендации вычисляются через позднее слияние (Late Fusion) трёх модальностей с оптимальными весами, подобранными экспериментально.

Проект реализован в виде микросервисного веб-приложения с Go-бэкендом, Python-воркером для ML-вычислений, поисковым движком Meilisearch, очередью RabbitMQ и адаптивным фронтендом на Next.js. Все компоненты контейнеризированы и управляются через Docker Compose.

## Технологический стек
| Компонент | Технологии |
|-----------------|------------|
| **Backend API** | Go (Gin, JWT, PostgreSQL, Redis, RabbitMQ, Meilisearch, zap) |
| **ML‑воркер** | Python (PyTorch, NumPy, pandas, h5py, scikit‑learn) |
| **Frontend** | Next.js (React, TypeScript, Tailwind CSS, React Query) |
| **Контейнеризация** | Docker, Docker Compose |

## Архитектура
Система состоит из микросервисов:
- `frontend` — клиентское приложение Next.js.
- `go-api` — REST API на Go (аутентификация, библиотека, рекомендации, поиск).
- `ml-worker` — Python‑воркер, подписанный на RabbitMQ, вычисляет рекомендации асинхронно.
- `meilisearch` — поисковый движок и batch‑запросы метаданных.
- `postgres` — реляционная БД для пользователей, их библиотеки и истории.
- `redis` — кэш статусов задач.
- `rabbitmq` — очередь сообщений для асинхронной обработки.

![Архитектура](assets/architecture.png)

## Установка и запуск
### Требования
- Docker и Docker Compose
- Git
### Шаги
1. Клонируйте репозиторий:
```bash
git clone https://github.com/I000000/Recly.git
cd Recly
```
2.  Скопируйте файл с переменными окружения и заполните его:
```bash
cp .env.example .env
# Отредактируйте .env, укажите свои пароли и ключи
```
3. Запустите все сервисы (первый запуск требует индексацию)
```bash
docker compose --profile manual up -d
```
4. Остановить и запустить контейнеры в будущем можно так
```bash
docker compose down
docker compose up -d
```
5. Откройте приложение
-   Фронтенд: [http://localhost:3000](http://localhost:3000/)
-   Swagger API: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)


## Основные эндпоинты API

| Метод | Путь | Описание |
|-------|------|----------|
| `POST` | `/api/auth/register` | Регистрация |
| `POST` | `/api/auth/login` | Вход |
| `POST` | `/api/book/:id/like` | Добавить книгу в библиотеку |
| `POST` | `/api/movie/:id/like` | Добавить фильм в библиотеку |
| `GET` | `/api/search` | Поиск |
| `POST` | `/api/recommend` | Запрос рекомендаций |
| `GET` | `/api/result/:taskId` | Получение результата |

Полная документация доступна через Swagger.

## Тестирование

Юнит-тесты:
```bash
cd backend
go test -v ./...
```
Интеграционные тесты (требуют Docker):
```bash
go test -v -tags=integration ./internal/repository/postgres/...
```