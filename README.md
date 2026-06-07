<div align="center">

# 🚀 Go Service Template

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go Version" />
  <img src="https://img.shields.io/badge/Fiber-v2.52+-00ACD7?style=for-the-badge&logo=fiber&logoColor=white" alt="Fiber" />
  <img src="https://img.shields.io/badge/PostgreSQL-15+-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
</p>

<p align="center">
  <img src="https://img.shields.io/badge/License-MIT-green.svg?style=for-the-badge" alt="License" />
  <img src="https://img.shields.io/badge/PRs-Welcome-brightgreen.svg?style=for-the-badge" alt="PRs Welcome" />
  <img src="https://img.shields.io/badge/Made%20with-❤️-red.svg?style=for-the-badge" alt="Made with Love" />
</p>

<h3 align="center">🎯 Современный шаблон для создания микросервисов на Go</h3>

<p align="center">
  <strong>Готовое решение с RESTful API, PostgreSQL, Docker и Swagger документацией</strong><br/>
  Начните разработку микросервиса за считанные минуты! ⚡
</p>

<p align="center">
  <a href="#-быстрый-старт">Быстрый старт</a> •
  <a href="#-api-endpoints">API</a> •
  <a href="#-разработка">Разработка</a> •
  <a href="#-конфигурация">Конфигурация</a> •
  <a href="LICENSE">Лицензия</a>
</p>

---

</div>

## ✨ Возможности

<table>
<tr>
<td width="50%">

### 🏗️ **Архитектура**
- 🎯 **Clean Architecture** - чистая архитектура
- 📦 **Модульная структура** - легко расширяемая
- 🔌 **Service Layer** - бизнес-логика отделена от HTTP
- 🏭 **Repository Pattern** - абстракция данных

</td>
<td width="50%">

### 🚀 **Технологии**
- ⚡ **Fiber Web Framework** - быстрый HTTP сервер
- 🐘 **PostgreSQL** - надежная база данных
- 🐳 **Docker** - контейнеризация
- 📚 **Swagger/OpenAPI** - автодокументация

</td>
</tr>
<tr>
<td width="50%">

### 🛠️ **DevOps & Инструменты**
- 📈 **Structured Logging** - slog
- 🔍 **Health Checks** - мониторинг
- 🔄 **Database Migrations** - версионирование БД
- 🧪 **Testing Ready** - готов к тестированию

</td>
<td width="50%">

### 📊 **Производительность**
- ⚡ **Высокая скорость** - Fiber framework
- 💾 **Connection Pooling** - эффективные соединения
- 🔄 **Graceful Shutdown** - корректная остановка
- 📦 **Lightweight** - минимальные зависимости

</td>
</tr>
</table>

## 📁 Структура проекта

```
go-service-template/
├── cmd/
│   └── service/           # Точка входа приложения
├── internal/
│   ├── config/           # Конфигурация
│   ├── models/           # Модели данных
│   ├── server/           # HTTP сервер и роуты
│   ├── service/          # Бизнес-логика + Storage интерфейс
│   │   ├── service.go    # Service интерфейс
│   │   ├── example.go    # Service реализация
│   │   └── storage.go    # Storage интерфейс
│   └── storage/          # Реализации хранилищ
│       └── postgres/     # PostgreSQL реализация Storage
├── migrations/           # SQL миграции
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile           # Docker образ
└── README.md
```

## 🏛️ Архитектура

Проект построен на принципах **Clean Architecture** с четким разделением слоев:

```
┌─────────────────┐
│   HTTP Layer    │  ← Handlers (REST API, парсинг запросов)
├─────────────────┤
│ Business Logic  │  ← Services (валидация, бизнес-правила)  
│    + Interfaces │  ← Storage интерфейс (DIP принцип)
├─────────────────┤
│   Data Access   │  ← PostgreSQL (реализует Storage)
└─────────────────┘
```

## 🛠 Быстрый старт

### Предварительные требования

- Go 1.26+
- Docker и Docker Compose
- PostgreSQL (если запускаете локально)

### 🐳 Запуск с Docker (рекомендуется)

1. **Клонируйте репозиторий:**
```bash
git clone <repository-url>
cd go-service-template
```

2. **Запустите одной командой:**
```bash
make start
# или
docker compose up --build
```

3. **Проверьте работу:**
```bash
# Проверки состояния
curl http://localhost:8080/livez   # liveness (без зависимостей)
curl http://localhost:8080/readyz  # readiness (проверяет БД)

# Создание примера
curl -X POST http://localhost:8080/api/v1/examples \
  -H "Content-Type: application/json" \
  -d '{"name":"Test","description":"Тестовый пример","value":100.5,"is_active":true}'

# Получение всех примеров
curl http://localhost:8080/api/v1/examples
```

4. **Доступные адреса:**

<table>
<tr>
<th>🔗 Сервис</th>
<th>📍 URL</th>
<th>📋 Описание</th>
</tr>
<tr>
<td><strong>🌐 API</strong></td>
<td><code>http://localhost:8080</code></td>
<td>Основной REST API</td>
</tr>
<tr>
<td><strong>📚 Swagger</strong></td>
<td><code>http://localhost:8080/swagger/</code></td>
<td>Документация API (включается <code>ENABLE_SWAGGER=true</code>; в Docker генерируется при сборке)</td>
</tr>
<tr>
<td><strong>❤️ Health</strong></td>
<td><code>/livez</code>, <code>/readyz</code></td>
<td>Liveness (без зависимостей) и readiness (проверка БД). <code>/health</code> — алиас readiness</td>
</tr>
</table>

> **💡 Swagger:** выключен по умолчанию (`ENABLE_SWAGGER=false`). В Docker-образе спецификация генерируется при сборке — достаточно поднять стек с `ENABLE_SWAGGER=true`. Локально используйте `make swagger-serve` (генерация + запуск).

### Локальный запуск

1. Установите зависимости:
```bash
go mod download
```

2. Настройте переменные окружения:
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_NAME=service_db
export DB_USER=postgres
export DB_PASSWORD=password   # обязательна — без неё конфиг не пройдёт валидацию
export DEBUG_MODE=true        # человекочитаемые логи для разработки
export ENABLE_SWAGGER=true    # включить Swagger UI на /swagger/
```

3. Запустите PostgreSQL и выполните миграции

4. Запустите приложение:
```bash
go run ./cmd/service
```

## 📋 API Endpoints

<details>
<summary><strong>🔍 Нажмите, чтобы развернуть полный список API</strong></summary>

### 🏥 Проверки состояния
```http
GET /livez    # liveness — 200, пока процесс жив (без зависимостей)
GET /readyz   # readiness — 200, если БД доступна, иначе 503
GET /health   # алиас readiness (обратная совместимость)
```
**Ответ `/readyz`:**
```json
{
  "message": "ready"
}
```

### 📝 Examples (CRUD операции)

#### Создание записи
```http
POST /api/v1/examples
Content-Type: application/json

{
  "name": "Пример",
  "description": "Описание примера",
  "value": 100.50,
  "is_active": true
}
```

#### Получение всех записей
```http
GET /api/v1/examples?limit=10&offset=0
```

#### Получение записи по ID
```http
GET /api/v1/examples/1
```

#### Обновление записи
```http
PUT /api/v1/examples/1
Content-Type: application/json

{
  "name": "Обновленный пример",
  "description": "Новое описание",
  "value": 200.75,
  "is_active": false
}
```

#### Удаление записи
```http
DELETE /api/v1/examples/1
```

### 📚 Документация
```http
GET /swagger/*
```

</details>

## 🔧 Конфигурация

Приложение настраивается через переменные окружения:

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| `DB_HOST` | Хост базы данных | `localhost` |
| `DB_PORT` | Порт базы данных | `5432` |
| `DB_NAME` | Название базы данных | `service_db` |
| `DB_USER` | Пользователь БД | `postgres` |
| `DB_PASSWORD` | Пароль БД (**обязательна**) | — |
| `DB_SSLMODE` | SSL-режим (`disable`/`require`/`verify-full`/…) | `disable` |
| `DB_MAX_CONNS` | Максимум коннектов пула | `10` |
| `DB_MIN_CONNS` | Минимум коннектов пула | `1` |
| `DB_MAX_CONN_LIFETIME` | Срок жизни коннекта | `1h` |
| `DB_MAX_CONN_IDLE_TIME` | Idle-время коннекта | `30m` |
| `SERVER_HOST` | Хост сервера | `localhost` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `SERVER_READ_TIMEOUT` | Таймаут чтения запроса | `10s` |
| `SERVER_WRITE_TIMEOUT` | Таймаут записи ответа | `10s` |
| `SERVER_BODY_LIMIT` | Макс. размер тела запроса, байт | `4194304` |
| `SERVER_RATE_LIMIT` | Лимит запросов/мин на IP (0 — выкл.) | `100` |
| `CORS_ALLOW_ORIGINS` | Разрешённые CORS-источники | `*` |
| `DEBUG_MODE` | Текстовые debug-логи вместо JSON | `false` |
| `ENABLE_SWAGGER` | Включить Swagger UI на `/swagger/` | `false` |

> **ℹ️ Примечание:** в таблице — значения по умолчанию из кода. Локальный стек (`.env.example` / `docker-compose.yml`) переопределяет часть из них: `SERVER_HOST=0.0.0.0`, `DB_MAX_CONNS=20`, `DB_MIN_CONNS=2`, `DB_PASSWORD=password`, `ENABLE_SWAGGER=true`.

## 📊 База данных

Проект использует PostgreSQL с системой миграций. Пример таблицы:

```sql
CREATE TABLE examples (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value DOUBLE PRECISION NOT NULL DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);
```

## 🧪 Разработка

### 🎯 Добавление новых эндпоинтов

<details>
<summary><strong>📖 Пошаговое руководство</strong></summary>

#### 1️⃣ Создайте модель
```go
// internal/models/user.go
type User struct {
    ID       int       `json:"id"`
    Name     string    `json:"name"`
    Email    string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

type UserRequest struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

#### 2️⃣ Добавьте методы в Storage
```go
// internal/service/storage.go
type Storage interface {
    // ... существующие методы
    CreateUser(ctx context.Context, user *models.User) error
    GetUserByID(ctx context.Context, id int) (*models.User, error)
}
```

#### 3️⃣ Реализуйте в PostgreSQL
```go
// internal/storage/postgres/storage.go
func (s *PostgresStorage) CreateUser(ctx context.Context, user *models.User) error {
    query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id`
    if err := s.pool.QueryRow(ctx, query, user.Name, user.Email, user.CreatedAt).Scan(&user.ID); err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}
```

#### 4️⃣ Добавьте методы в Service
```go
// internal/service/service.go
type Service interface {
    // ... существующие методы
    CreateUser(ctx context.Context, req *models.UserRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id int) (*models.User, error)
}

// internal/service/errors.go — объявите sentinel-ошибки (паттерн уже есть в проекте)
var (
    ErrUserNameRequired = errors.New("name is required")
    ErrCreateUserFailed = errors.New("failed to create user")
)

// internal/service/user.go (или в example.go)
func (s *service) CreateUser(ctx context.Context, req *models.UserRequest) (*models.User, error) {
    if strings.TrimSpace(req.Name) == "" {
        return nil, ErrUserNameRequired
    }

    user := &models.User{
        Name:      strings.TrimSpace(req.Name),
        Email:     strings.TrimSpace(req.Email),
        CreatedAt: time.Now(),
    }

    if err := s.storage.CreateUser(ctx, user); err != nil {
        s.logger.Error("failed to create user", slog.String("error", err.Error()))
        return nil, ErrCreateUserFailed
    }
    return user, nil
}
```

#### 5️⃣ Создайте обработчики
```go
// internal/server/handlers.go
// @Summary Create user
// @Tags users
func (s *Server) createUser(c *fiber.Ctx) error {
    var req models.UserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(models.ErrorResponse{Error: "invalid request body"})
    }

    user, err := s.services.User.CreateUser(c.UserContext(), &req)
    if err != nil {
        return s.handleServiceError(c, err) // маппинг sentinel-ошибок в HTTP-статус
    }

    return c.Status(fiber.StatusCreated).JSON(user)
}
```

#### 6️⃣ Добавьте роуты
```go
// internal/server/server.go
users := api.Group("/users")
users.Post("/", s.createUser)
```

#### 7️⃣ Зарегистрируйте ресурс в двух местах
- Добавьте новый сервис в структуру `Services` и в `NewServices` — `internal/service/service.go`.
- Сопоставьте новые sentinel-ошибки с HTTP-статусами в `mapServiceErrorToHTTPStatus` — `internal/server/handlers.go`.

</details>

### 🔧 Полезные команды

```bash
# Разработка
make run              # Запуск в dev режиме
make test             # Запуск тестов
make test-coverage    # Тесты с покрытием
make fmt              # Форматирование кода
make lint             # Линтинг

# Docker
make start            # Запуск с Docker
make stop             # Остановка
make docker-logs      # Просмотр логов

# База данных
make migrate-up       # Применить миграции
make migrate-down     # Откатить миграции

# 📚 Swagger документация
make swagger-gen      # Генерировать документацию
make swagger-serve    # Запуск с обновленной документацией
make install-swag     # Установить swag CLI
```

### 📚 Swagger документация

#### Генерация документации

```bash
# Установка swag CLI (автоматически при первом использовании)
make install-swag

# Генерация документации из аннотаций в коде
make swagger-gen

# Запуск сервера с автоматической генерацией документации  
make swagger-serve
```

#### Добавление аннотаций

Добавляйте комментарии к обработчикам:

```go
// createExample создает новый пример
// @Summary Create example
// @Description Creates a new example record
// @Tags examples
// @Accept json
// @Produce json
// @Param example body models.ExampleRequest true "Example data"
// @Success 201 {object} models.Example
// @Failure 400 {object} models.ErrorResponse
// @Router /examples [post]
func (s *Server) createExample(c *fiber.Ctx) error {
    // реализация
}
```

После добавления новых аннотаций обязательно перегенерируйте документацию:
```bash
make swagger-gen
```

### 🧪 Тестирование

Проект использует стандартный `testing` (без сторонних assert-библиотек): моки за интерфейсами + хелперы `newTestServer`/`doRequest`.

```go
// internal/server/handlers_test.go
func TestCreateExample(t *testing.T) {
    mock := &mockExampleService{
        createFn: func(_ context.Context, req *models.ExampleRequest) (*models.Example, error) {
            return &models.Example{ID: 1, Name: req.Name}, nil
        },
    }
    s := newTestServer(mock, nil)

    resp := doRequest(s, http.MethodPost, "/api/v1/examples", models.ExampleRequest{Name: "test"})
    if resp.StatusCode != http.StatusCreated {
        t.Fatalf("ожидался 201, получен %d", resp.StatusCode)
    }
}
```

## 📝 Лицензия

```
MIT License

Copyright (c) 2025 Eduard Gagite
```

Этот проект распространяется под лицензией MIT. Подробности смотрите в файле [LICENSE](LICENSE).

---


<div align="center">

**Go Service Template** - современный шаблон микросервиса на Go  
Создано с ❤️ для Go сообщества

[![MIT License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.26+-blue.svg)](https://golang.org)
[![Docker](https://img.shields.io/badge/Docker-Ready-blue.svg)](Dockerfile)

</div> 
