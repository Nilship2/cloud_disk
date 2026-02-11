@echo off
chcp 65001 >nul
echo ===================================
echo     初始化数据库
echo ===================================
echo.

echo 1. 检查MySQL容器是否运行...
docker ps | findstr cloud-disk-mysql
if %errorlevel% neq 0 (
    echo [31m❌ MySQL容器未运行,正在启动...[0m
    docker-compose up -d mysql
    timeout /t 5
)

echo.
echo 2. 创建数据库 cloud_disk_dev...
docker exec -i cloud-disk-mysql mysql -uroot -prootpassword << EOF
CREATE DATABASE IF NOT EXISTS cloud_disk_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
SHOW DATABASES;
EOF

echo.
echo [32m✅ 数据库初始化完成！[0m
echo.
pause