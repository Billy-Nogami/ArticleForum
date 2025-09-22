# ArticleForum - GraphQL Forum Application

### Tребования

* Docker
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
    title: "Тестовый пост", 
    content: "Содержание тестового поста", 
    commentsEnabled: true
  ) {
    id
    title
    commentsEnabled
    createdAt
  }
}
```
**Комментарии поста**
```graphql
query {
  comments(postID: "5a1e2aaa-34fa-430a-96e9-72f2817cb720", limit: 10, offset: 0) {
    id
    postID
    content
    createdAt
  }
}
```

**Создать комментарий**
```graphql
mutation {
  createComment(
    postID: "5a1e2aaa-34fa-430a-96e9-72f2817cb720", 
    content: "Это мой первый комментарий!"
  ) {
    id
    postID
    content
    createdAt
  }
}
```

**Вернуть комментарии поста**
```graphql
query {
  comments(postID: "ID_ВАШЕГО_ПОСТА", limit: 10, offset: 0) {
    id
    postID
    content
    createdAt
  }
}
```

**Вернуть посты**
```graphql
query {
  posts {
    id
    title
    content
    commentsEnabled
    createdAt
  }
}
```
