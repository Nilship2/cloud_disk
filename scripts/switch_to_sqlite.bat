@echo off
chcp 65001 >nul
echo ===================================
echo     切换到 SQLite 数据库
echo ===================================
echo.

echo 1. 安装 SQLite 驱动...
go get -u gorm.io/driver/sqlite

echo.
echo 2. 创建 SQLite 配置文件...
copy nul config.sqlite.yaml >nul
(
echo server:
echo   env: "development"
echo   port: 8080
echo   host: "localhost"
echo   read_timeout: 30
echo   write_timeout: 30
echo.
echo database:
echo   driver: "sqlite"
echo.
echo jwt:
echo   secret: "development-secret-key"
echo   expires_hours: 24
echo.
echo log:
echo   level: "debug"
echo   format: "console"
echo   output: "stdout"
echo.
echo storage:
echo   type: "local"
echo   local:
echo     base_path: "./storage/uploads"
echo     temp_path: "./storage/uploads/temp"
echo     max_size_mb: 100
) > config.sqlite.yaml

echo.
echo 3. 更新依赖...
go mod tidy

echo.
echo [32m✅ SQLite 配置完成！[0m
echo.
echo 现在运行：go run cmd/server/main.go
echo.
pause