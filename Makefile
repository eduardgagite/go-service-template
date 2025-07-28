.PHONY: build run test clean docker-build docker-up docker-down migrate-up migrate-down deps help

# 🎨 Цвета для красивого вывода
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[0;33m
BLUE=\033[0;34m
PURPLE=\033[0;35m
CYAN=\033[0;36m
WHITE=\033[0;37m
BOLD=\033[1m
NC=\033[0m # No Color

# 📋 Переменные
APP_NAME=service
BINARY_DIR=bin
DOCKER_COMPOSE=docker-compose

# 🎯 Цель по умолчанию
.DEFAULT_GOAL := help

# ═══════════════════════════════════════════════════════════════════════════════
# 🏗️  СБОРКА И ЗАПУСК
# ═══════════════════════════════════════════════════════════════════════════════

build: ## 🔨 Собрать приложение
	@echo "$(CYAN)🔨 Сборка $(APP_NAME)...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(APP_NAME) ./cmd/service
	@echo "$(GREEN)✅ Сборка завершена: $(BINARY_DIR)/$(APP_NAME)$(NC)"

run: ## 🚀 Запустить приложение локально
	@echo "$(BLUE)🚀 Запуск $(APP_NAME)...$(NC)"
	@go run ./cmd/service

# ═══════════════════════════════════════════════════════════════════════════════
# 🧪 ТЕСТИРОВАНИЕ И КАЧЕСТВО КОДА
# ═══════════════════════════════════════════════════════════════════════════════

test: ## 🧪 Запустить тесты
	@echo "$(YELLOW)🧪 Запуск тестов...$(NC)"
	@go test -v ./...
	@echo "$(GREEN)✅ Тесты завершены$(NC)"

test-coverage: ## 📊 Запустить тесты с покрытием
	@echo "$(YELLOW)📊 Запуск тестов с покрытием...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✅ Отчет о покрытии создан: coverage.html$(NC)"

fmt: ## 🎨 Форматировать код
	@echo "$(PURPLE)🎨 Форматирование кода...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✅ Код отформатирован$(NC)"

lint: ## 🔍 Проверить код линтером
	@echo "$(PURPLE)🔍 Запуск линтера...$(NC)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
		echo "$(GREEN)✅ Линтинг завершен$(NC)"; \
	else \
		echo "$(RED)❌ golangci-lint не установлен. Установите: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.54.2$(NC)"; \
	fi

check: fmt lint test ## ✅ Полная проверка (форматирование + линтинг + тесты)
	@echo "$(GREEN)$(BOLD)🎉 Все проверки пройдены успешно!$(NC)"

# ═══════════════════════════════════════════════════════════════════════════════
# 🐳 DOCKER КОМАНДЫ
# ═══════════════════════════════════════════════════════════════════════════════

docker-build: ## 🐳 Собрать Docker образ
	@echo "$(CYAN)🐳 Сборка Docker образа...$(NC)"
	@$(DOCKER_COMPOSE) build
	@echo "$(GREEN)✅ Docker образ собран$(NC)"

docker-up: ## ⬆️  Запустить Docker контейнеры
	@echo "$(BLUE)⬆️ Запуск Docker контейнеров...$(NC)"
	@$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✅ Контейнеры запущены$(NC)"

docker-down: ## ⬇️  Остановить Docker контейнеры
	@echo "$(YELLOW)⬇️ Остановка Docker контейнеров...$(NC)"
	@$(DOCKER_COMPOSE) down
	@echo "$(GREEN)✅ Контейнеры остановлены$(NC)"

docker-logs: ## 📋 Показать логи Docker контейнеров
	@echo "$(CYAN)📋 Показ логов Docker...$(NC)"
	@$(DOCKER_COMPOSE) logs -f

start: docker-build docker-up ## 🌟 Полный запуск с Docker (сборка + запуск)
	@echo ""
	@echo "$(GREEN)$(BOLD)🎉 Сервис успешно запущен!$(NC)"
	@echo "$(CYAN)┌──────────────────────────────────────────────┐$(NC)"
	@echo "$(CYAN)│  🌐 API:     http://localhost:8080          │$(NC)"
	@echo "$(CYAN)│  📚 Swagger: http://localhost:8080/swagger/  │$(NC)"
	@echo "$(CYAN)│  ❤️  Health:  http://localhost:8080/health   │$(NC)"
	@echo "$(CYAN)└──────────────────────────────────────────────┘$(NC)"

stop: docker-down ## 🛑 Остановить сервис
	@echo "$(RED)🛑 Сервис остановлен$(NC)"

# ═══════════════════════════════════════════════════════════════════════════════
# 🗄️  БАЗА ДАННЫХ
# ═══════════════════════════════════════════════════════════════════════════════

migrate-up: ## ⬆️  Применить миграции базы данных
	@echo "$(BLUE)⬆️ Применение миграций...$(NC)"
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgres://postgres:password@localhost:5432/service_db?sslmode=disable" up; \
		echo "$(GREEN)✅ Миграции применены$(NC)"; \
	else \
		echo "$(RED)❌ migrate CLI не установлен. Установите: https://github.com/golang-migrate/migrate$(NC)"; \
	fi

migrate-down: ## ⬇️  Откатить миграции базы данных
	@echo "$(YELLOW)⬇️ Откат миграций...$(NC)"
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "postgres://postgres:password@localhost:5432/service_db?sslmode=disable" down; \
		echo "$(GREEN)✅ Миграции откачены$(NC)"; \
	else \
		echo "$(RED)❌ migrate CLI не установлен. Установите: https://github.com/golang-migrate/migrate$(NC)"; \
	fi

# ═══════════════════════════════════════════════════════════════════════════════
# 🛠️  УТИЛИТЫ
# ═══════════════════════════════════════════════════════════════════════════════

deps: ## 📦 Установить зависимости
	@echo "$(PURPLE)📦 Установка зависимостей...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)✅ Зависимости установлены$(NC)"

install-swag: ## 🔧 Установить swag CLI для генерации Swagger документации
	@echo "$(PURPLE)🔧 Установка swag CLI...$(NC)"
	@if command -v swag > /dev/null; then \
		echo "$(GREEN)✅ swag уже установлен$(NC)"; \
	else \
		go install github.com/swaggo/swag/cmd/swag@latest; \
		echo "$(GREEN)✅ swag установлен$(NC)"; \
	fi

swagger-gen: install-swag ## 📚 Генерировать Swagger документацию
	@echo "$(CYAN)📚 Генерация Swagger документации...$(NC)"
	@swag init -g cmd/service/main.go -o docs --parseInternal
	@echo "$(GREEN)✅ Swagger документация сгенерирована в папке docs/$(NC)"

swagger-serve: swagger-gen ## 🌐 Запустить сервер с обновленной документацией
	@echo "$(BLUE)🌐 Запуск сервера с актуальной Swagger документацией...$(NC)"
	@go run ./cmd/service &
	@sleep 2
	@echo "$(GREEN)✅ Сервер запущен!$(NC)"
	@echo "$(CYAN)┌──────────────────────────────────────────────┐$(NC)"
	@echo "$(CYAN)│  📚 Swagger: http://localhost:8080/swagger/  │$(NC)"
	@echo "$(CYAN)│  🌐 API:     http://localhost:8080          │$(NC)"
	@echo "$(CYAN)└──────────────────────────────────────────────┘$(NC)"

clean: ## 🧹 Очистить артефакты сборки
	@echo "$(YELLOW)🧹 Очистка...$(NC)"
	@rm -rf $(BINARY_DIR)
	@rm -rf docs/
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)✅ Очистка завершена$(NC)"

# ═══════════════════════════════════════════════════════════════════════════════
# ℹ️  СПРАВКА
# ═══════════════════════════════════════════════════════════════════════════════

help: ## 📖 Показать справку
	@echo "$(BOLD)$(BLUE)"
	@echo "╔══════════════════════════════════════════════════════════════════════════════╗"
	@echo "║                          🚀 Go Service Template                             ║"
	@echo "║                     Система команд для разработки                           ║"
	@echo "╚══════════════════════════════════════════════════════════════════════════════╝"
	@echo "$(NC)"
	@echo "$(CYAN)🏗️  СБОРКА И ЗАПУСК:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /СБОРКА И ЗАПУСК/ || ($$0 ~ /(build|run):/ && !printed)) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)🧪 ТЕСТИРОВАНИЕ И КАЧЕСТВО:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /(test|fmt|lint|check):/) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BLUE)🐳 DOCKER:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /(docker|start|stop):/) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(PURPLE)🗄️  БАЗА ДАННЫХ:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /migrate/) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(YELLOW)📚 SWAGGER ДОКУМЕНТАЦИЯ:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /swagger/) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(CYAN)🛠️  УТИЛИТЫ:$(NC)"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*?##/ { if ($$0 ~ /(deps|install-swag|clean|help):/) printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "$(BOLD)💡 Примеры использования:$(NC)"
	@echo "  $(YELLOW)make start$(NC)           - Быстрый запуск с Docker"
	@echo "  $(YELLOW)make run$(NC)             - Локальная разработка"
	@echo "  $(YELLOW)make swagger-serve$(NC)   - Запуск с генерацией документации"
	@echo "  $(YELLOW)make swagger-gen$(NC)     - Только генерация Swagger"
	@echo "  $(YELLOW)make check$(NC)           - Полная проверка кода"
	@echo "  $(YELLOW)make test-coverage$(NC)   - Тесты с отчетом о покрытии"
	@echo "" 