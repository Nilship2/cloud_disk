@echo off
chcp 65001 >nul
echo ====================================
echo     修复数据库密码问题
echo ====================================
echo.

cd /d C:\Users\lenovo\Documents\GitHub\cloud_disk

echo 1. 创建 .env.production 文件...
(
echo # MySQL配置
echo MYSQL_ROOT_PASSWORD=rootpassword
echo MYSQL_PASSWORD=cloudpassword
echo.
echo # Redis配置
echo REDIS_PASSWORD=redispassword
echo.
echo # JWT配置
echo JWT_SECRET=your-32-character-jwt-secret-key-change-this
echo.
echo # MinIO配置
echo MINIO_ROOT_USER=minioadmin
echo MINIO_ROOT_PASSWORD=minioadmin
) > .env.production

echo 2. 加载环境变量...
for /f "tokens=*" %%a in (.env.production) do set %%a

echo 3. 重新启动容器...
docker-compose down
docker-compose up -d

echo 4. 等待服务启动...
timeout /t 10

echo 5. 查看应用日志...
docker-compose logs --tail=20 app

echo.
echo [92m✅ 修复完成！[0m
echo 测试: curl http://localhost:8080/health
echo.
pause