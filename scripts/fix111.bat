@echo off
chcp 65001 >nul
echo ====================================
echo     紧急修复脚本
echo ====================================
echo.

echo 1. 停止应用容器...
docker stop cloud-disk-app
echo.

echo 2. 创建正确的配置文件...
(
echo server:
echo   env: "production"
echo   port: 8080
echo   host: "0.0.0.0"
echo.
echo database:
echo   driver: "mysql"
echo   host: "mysql"
echo   port: 3306
echo   username: "root"
echo   password: "rootpassword"
echo   database: "cloud_disk"
echo   charset: "utf8mb4"
echo   max_idle_conns: 10
echo   max_open_conns: 100
echo   conn_max_lifetime: 3600
echo.
echo jwt:
echo   secret: "your-32-character-jwt-secret-key-change-this"
echo   expires_hours: 24
echo.
echo log:
echo   level: "info"
echo   format: "json"
echo   output: "stdout"
echo.
echo storage:
echo   type: "minio"
echo   minio:
echo     endpoint: "minio:9000"
echo     access_key: "minioadmin"
echo     secret_key: "minioadmin"
echo     use_ssl: false
echo     bucket_name: "cloud-disk"
) > config.fixed.yaml

echo 3. 复制配置文件到容器...
docker cp config.fixed.yaml cloud-disk-app:/app/config.yaml
echo.

echo 4. 启动应用...
docker start cloud-disk-app
echo.

echo 5. 等待启动...
timeout /t 5
echo.

echo 6. 查看日志...
docker logs cloud-disk-app --tail 30
echo.

echo 7. 如果还有问题，请运行：docker exec -it cloud-disk-app sh
echo.
pause