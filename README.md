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
docker-compose up --build
```

3. **Проверьте работу:**
```bash
# Health check
curl http://localhost:8080/health

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
<td>Интерактивная документация API<br/><small>⚠️ Требует генерации: <code>make swagger-gen</code></small></td>
</tr>
<tr>
<td><strong>❤️ Health Check</strong></td>
<td><code>http://localhost:8080/health</code></td>
<td>Проверка работоспособности сервиса</td>
</tr>
</table>

> **💡 Важно:** Для работы Swagger документации необходимо сначала её сгенерировать командой `make swagger-gen` или использовать `make swagger-serve` для автоматической генерации и запуска.

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
export DB_PASSWORD=password
export DEBUG_MODE=true
```

3. Запустите PostgreSQL и выполните миграции

4. Запустите приложение:
```bash
go run ./cmd/service
```

## 📋 API Endpoints

<details>
<summary><strong>🔍 Нажмите, чтобы развернуть полный список API</strong></summary>

### 🏥 Health Check
```http
GET /health
```
**Ответ:**
```json
{
  "message": "Service is healthy"
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
| `DB_PASSWORD` | Пароль БД | `` |
| `DB_SSLMODE` | SSL режим для БД | `disable` |
| `SERVER_HOST` | Хост сервера | `localhost` |
| `SERVER_PORT` | Порт сервера | `8080` |
| `DEBUG_MODE` | Режим отладки | `false` |
| `DB_MAX_CONNS` | Максимум коннектов пула | `10` |
| `DB_MIN_CONNS` | Минимум коннектов пула | `1` |
| `DB_MAX_CONN_LIFETIME` | Срок жизни коннекта | `1h` |
| `DB_MAX_CONN_IDLE_TIME` | Idle-время коннекта | `30m` |

## 📊 База данных

Проект использует PostgreSQL с системой миграций. Пример таблицы:

```sql
CREATE TABLE examples (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    value DECIMAL(10,2) DEFAULT 0.00,
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
    CreateUser(user *User) error
    GetUserByID(id int) (*User, error)
}
```

#### 3️⃣ Реализуйте в PostgreSQL
```go
// internal/storage/postgres/storage.go
import "go-service-template/internal/service"

// PostgresStorage реализует service.Storage
func (s *PostgresStorage) CreateUser(user *User) error {
    query := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id`
    err := s.db.QueryRow(query, user.Name, user.Email, user.CreatedAt).Scan(&user.ID)
    return err
}
```

#### 4️⃣ Добавьте методы в Service
```go
// internal/service/service.go
type Service interface {
    // ... существующие методы
    CreateUser(req *UserRequest) (*User, error)
    GetUserByID(id int) (*User, error)
}

// internal/service/user.go (или в example.go)
func (s *service) CreateUser(req *UserRequest) (*User, error) {
    // Валидация
    if strings.TrimSpace(req.Name) == "" {
        return nil, errors.New("name is required")
    }
    
    // Создание модели
    user := &User{
        Name:      strings.TrimSpace(req.Name),
        Email:     strings.TrimSpace(req.Email),
        CreatedAt: time.Now(),
    }
    
    // Сохранение
    if err := s.storage.CreateUser(user); err != nil {
        return nil, errors.New("failed to create user")
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
    var req UserRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(ErrorResponse{Error: err.Error()})
    }
    
    user, err := s.services.Example.CreateUser(&req)
    if err != nil {
        return c.Status(400).JSON(ErrorResponse{Error: err.Error()})
    }
    
    return c.Status(201).JSON(user)
}
```

#### 6️⃣ Добавьте роуты
```go
// internal/server/server.go
users := api.Group("/users")
users.Post("/", s.createUser)
```

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

```go
// internal/server/handlers_test.go
func TestCreateExample(t *testing.T) {
    // Настройка тестового сервера
    app := fiber.New()
    
    // Тестовый запрос
    req := httptest.NewRequest("POST", "/api/v1/examples", nil)
    resp, _ := app.Test(req)
    
    assert.Equal(t, 201, resp.StatusCode)
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
