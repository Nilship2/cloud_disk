@echo off
chcp 65001 >nul
title 第五天功能验证
echo ====================================
echo    回收站与监控功能快速验证
echo ====================================
echo.

:: 检查服务器
curl -s http://localhost:8080/health >nul
if %errorlevel% neq 0 (
    echo [31m❌ 服务器未运行！请先启动服务器[0m
    pause
    exit /b 1
)

:: 登录
set /p username=请输入用户名: 
set /p password=请输入密码: 

curl -s -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"%username%\",\"password\":\"%password%\"}" > %temp%\login.json

for /f "tokens=*" %%i in ('powershell -Command "$json = Get-Content '%temp%\login.json' | ConvertFrom-Json; $json.data.token"') do set TOKEN=%%i

if "%TOKEN%"=="" (
    echo [31m❌ 登录失败[0m
    pause
    exit /b 1
)

echo [32m✅ 登录成功[0m
echo.

:: 1. 查看回收站统计
echo [33m1. 回收站统计...[0m
curl -X GET "http://localhost:8080/api/v1/trash/stats" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo ------------------------
pause

:: 2. 获取回收站列表
echo [33m2. 回收站列表...[0m
curl -X GET "http://localhost:8080/api/v1/trash?page=1&page_size=5" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo ------------------------
pause

:: 3. 健康检查
echo [33m3. 健康检查...[0m
curl http://localhost:8080/health
echo.
echo ------------------------
pause

:: 4. 系统监控
echo [33m4. 系统统计...[0m
curl -X GET "http://localhost:8080/api/v1/monitor/stats" ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo ------------------------
pause

echo.
echo [32m✅ 验证完成！[0m
pause