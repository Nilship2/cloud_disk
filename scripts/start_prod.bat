@echo off
chcp 65001 >nul
title 云盘系统-生产环境
echo ====================================
echo     启动云盘系统生产环境
echo ====================================
echo.

cd /d C:\cloud-disk

echo 1. 检查日志目录...
if not exist logs mkdir logs

echo 2. 加载依赖...
go mod download

echo 3. 编译应用...
set GO_ENV=production
go build -o cloud-disk.exe ./cmd/server/main.go

if %errorlevel% neq 0 (
    echo [31m❌ 编译失败[0m
    pause
    exit /b 1
)

echo [32m✅ 编译成功[0m
echo.

echo 4. 启动服务...
echo 日志文件: C:\cloud-disk\logs\app.log
echo.

start /B cloud-disk.exe > logs\console.log 2>&1

echo 5. 等待启动...
timeout /t 5

echo 6. 测试服务...
curl http://localhost:8080/health

echo.
echo [32m✅ 服务已启动！[0m
echo 访问地址: http://localhost:8080
echo 查看日志: type logs\app.log
echo.
pause