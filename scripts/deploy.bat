@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo =====================================
echo     云盘系统部署脚本
echo =====================================
echo.

:: 颜色定义
set GREEN=[92m
set RED=[91m
set YELLOW=[93m
set NC=[0m

echo %YELLOW%1. 检查Docker环境...%NC%
docker --version >nul 2>&1
if %errorlevel% neq 0 (
    echo %RED%❌ Docker未安装%NC%
    pause
    exit /b 1
)

docker-compose --version >nul 2>&1
if %errorlevel% neq 0 (
    echo %RED%❌ Docker Compose未安装%NC%
    pause
    exit /b 1
)
echo %GREEN%✓ Docker环境就绪%NC%
echo.

echo %YELLOW%2. 检查环境变量文件...%NC%
if not exist .env.production (
    echo %RED%❌ .env.production文件不存在%NC%
    pause
    exit /b 1
)
echo %GREEN%✓ 环境变量文件存在%NC%
echo.

echo %YELLOW%3. 构建Docker镜像...%NC%
docker-compose build --no-cache
if %errorlevel% neq 0 (
    echo %RED%❌ 镜像构建失败%NC%
    pause
    exit /b 1
)
echo %GREEN%✓ 镜像构建成功%NC%
echo.

echo %YELLOW%4. 停止旧容器...%NC%
docker-compose down
echo %GREEN%✓ 旧容器已停止%NC%
echo.

echo %YELLOW%5. 启动新容器...%NC%
docker-compose up -d
if %errorlevel% neq 0 (
    echo %RED%❌ 容器启动失败%NC%
    pause
    exit /b 1
)
echo %GREEN%✓ 容器启动成功%NC%
echo.

echo %YELLOW%6. 等待服务启动...%NC%
timeout /t 10 /nobreak >nul

echo %YELLOW%7. 健康检查...%NC%
curl -s http://localhost:8080/health | findstr "healthy" >nul
if %errorlevel% equ 0 (
    echo %GREEN%✓ 服务运行正常%NC%
) else (
    echo %RED%❌ 服务异常%NC%
    docker-compose logs --tail=20
    pause
    exit /b 1
)
echo.

echo %GREEN%
echo =====================================
echo    部署完成！
echo =====================================
echo API地址: http://localhost:8080
echo Swagger文档: http://localhost:8080/swagger/index.html
echo MinIO控制台: http://localhost:9001
echo =====================================
echo%NC%

pause