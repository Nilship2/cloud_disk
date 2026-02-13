@echo off
chcp 65001 >nul
echo ====================================
echo     切换到 MySQL 数据库
echo ====================================
echo.

cd /d C:\Users\lenovo\Documents\GitHub\cloud_disk

echo 1. 检查MySQL容器...
docker-compose up -d mysql
timeout /t 5

echo.
echo 2. 创建数据库...
docker exec -i cloud-disk-mysql mysql -uroot -prootpassword -e "CREATE DATABASE IF NOT EXISTS cloud_disk_dev CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
docker exec -i cloud-disk-mysql mysql -uroot -prootpassword -e "SHOW DATABASES;"

echo.
echo 3. 修改配置文件...
copy config.mysql.yaml config.dev.yaml /Y

echo.
echo 4. 更新依赖...
go mod tidy

echo.
echo [32m✅ 切换完成！[0m
echo.
echo 现在运行：go run cmd/server/main.go
echo.
pause