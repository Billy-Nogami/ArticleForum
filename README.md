# ArticleForum - GraphQL Forum Application

### Tребования

* Docker и Docker Compose
* Go 1.24+ (для локальной разработки)

### Вариант 1: Использование Docker (рекомендуется)

1. Клонирование репозитория
```bash
git clone https://github.com/Billy-Nogami/ArticleForum.git
cd ArticleForum
go mod tidy
```

2. Запуск с PostgreSQL (продакшен)
```bash
# Запуск PostgreSQL и приложения
docker run -d \
  --name postgres \
  -e POSTGRES_DB=articleforum \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=12345678 \
  -p 5432:5432 \
  postgres:13

# Сборка и запуск приложения
docker build -t articleforum .
docker run -d \
  --name articleforum-app \
  --link postgres:postgres \
  -e POSTGRES_HOST=postgres \
  -e POSTGRES_PORT=5432 \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=12345678 \
  -e POSTGRES_DB=articleforum \
  -p 8080:8080 \
  articleforum
```

3. Запуск с хранением в памяти 
```bash
# Сборка образа
docker build -t articleforum .

# Запуск с хранением в памяти
docker run -d \
  --name articleforum-app \
  -p 8080:8080 \
  articleforum ./articleforum -storage memory
```

### Вариант 2: Локальная разработка

1. Клонирование и настройка
```bash
git clone <your-repository-url>
cd ArticleForum
go mod download
```

2. Запуск с PostgreSQL
```bash
# Запуск PostgreSQL (если еще не запущен)
docker run -d \
  --name postgres \
  -e POSTGRES_DB=articleforum \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=12345678 \
  -p 5432:5432 \
  postgres:13

# Применение миграций
goose -dir migrations/goose postgres "user=postgres password=12345678 dbname=articleforum host=localhost port=5432 sslmode=disable" up

# Запуск приложения
go run ./cmd/server -storage postgres -postgres-dsn "postgres://postgres:12345678@localhost:5432/articleforum?sslmode=disable"
```

3. Запуск с хранением в памяти
```bash
go run ./cmd/server -storage memory
```

## Переменные окружения

### Переменные приложения
* `STORAGE_TYPE` - тип хранилища: memory или postgres (по умолчанию: memory)
* `PORT` - порт сервера (по умолчанию: 8080)

### Переменные PostgreSQL (требуются при использовании postgres storage)
* `POSTGRES_HOST` - хост PostgreSQL (по умолчанию: localhost)
* `POSTGRES_PORT` - порт PostgreSQL (по умолчанию: 5432)
* `POSTGRES_USER` - пользователь PostgreSQL (по умолчанию: postgres)
* `POSTGRES_PASSWORD` - пароль PostgreSQL
* `POSTGRES_DB` - имя базы данных (по умолчанию: articleforum)
* `POSTGRES_SSLMODE` - режим SSL (по умолчанию: disable)

## Использование API

### GraphQL Playground
После запуска приложения, доступ к GraphQL Playground:
```
http://localhost:8080
```

### Примеры запросов

**Создание поста**
```graphql
mutation {
  createPost(
    title: "My First Post"
    content: "This is the content of my post"
