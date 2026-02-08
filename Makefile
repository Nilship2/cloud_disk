.PHONY: help build run test clean lint migrate docker-up docker-down

# 颜色定义
GREEN  := $(shell tput -Txterm setaf 2)
YELLOW := $(shell tput -Txterm setaf 3)
WHITE  := $(shell tput -Txterm setaf 7)
RESET  := $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM=20

help:
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target>${RESET}'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  ${YELLOW}%-$(TARGET_MAX_CHAR_NUM)s${RESET} ${GREEN}%s${RESET}\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

## 安装依赖
install:
	@echo "Installing dependencies..."
	go mod download

## 构建应用
build:
	@echo "Building application..."
	go build -o bin/server cmd/server/main.go

## 运行开发服务器
run:
	@echo "Starting development server..."
	go run cmd/server/main.go

## 运行测试
test:
	@echo "Running tests..."
	go test ./... -v

## 运行覆盖率测试
test-coverage:
	@echo "Running tests with coverage..."
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

## 清理构建文件
clean:
	@echo "Cleaning up..."
	rm -rf bin/
	rm -rf coverage.out coverage.html

## 代码检查
lint:
	@echo "Running linter..."
	golangci-lint run

## 格式化代码
fmt:
	@echo "Formatting code..."
	gofmt -w .

## 更新依赖
update:
	@echo "Updating dependencies..."
	go get -u ./...
	go mod tidy

## 启动Docker环境
docker-up:
	@echo "Starting Docker environment..."
	docker-compose up -d

## 停止Docker环境
docker-down:
	@echo "Stopping Docker environment..."
	docker-compose down

## 查看日志
docker-logs:
	@echo "Showing logs..."
	docker-compose logs -f

## 数据库迁移
migrate-up:
	@echo "Running database migrations..."
	migrate -path scripts/database/migrations -database "mysql://root:password@tcp(localhost:3306)/cloud_disk" up

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path scripts/database/migrations -database "mysql://root:password@tcp(localhost:3306)/cloud_disk" down

## 生成Swagger文档
swagger:
	@echo "Generating Swagger docs..."
	swag init -g cmd/server/main.go -o api/docs

## 查看帮助
.DEFAULT_GOAL := help