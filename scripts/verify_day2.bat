@echo off
chcp 65001 >nul
title 云盘系统-第二天功能验证
echo ===================================
echo    第二天开发内容验证脚本
echo ===================================
echo.
echo 请确保服务器已经启动 (http://localhost:8080)
echo.

:: 检查服务器是否运行
curl -s http://localhost:8080/health >nul
if %errorlevel% neq 0 (
    echo [31m❌ 错误: 无法连接到服务器！请先运行 go run cmd/server/main.go[0m
    echo.
    pause
    exit /b 1
)

echo [32m✅ 服务器连接成功[0m
echo.
echo ===================================
echo.

:: 生成唯一的用户名（使用时间戳）
set timestamp=%date:~0,4%%date:~5,2%%date:~8,2%%time:~0,2%%time:~3,2%%time:~6,2%
set timestamp=%timestamp: =0%
set USERNAME=testuser_%timestamp%
set EMAIL=%USERNAME%@example.com
set PASSWORD=123456

echo 正在使用测试账号:
echo   用户名: %USERNAME%
echo   邮箱: %EMAIL%
echo   密码: %PASSWORD%
echo.
echo ===================================
echo.

:: 1. 测试注册
echo [33m1. 测试用户注册...[0m
echo ------------------------
curl -X POST http://localhost:8080/api/v1/register ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"%USERNAME%\",\"email\":\"%EMAIL%\",\"password\":\"%PASSWORD%\",\"confirm_password\":\"%PASSWORD%\"}"
echo.
echo.
echo ------------------------
echo.

:: 暂停一下，让用户看清结果
pause
echo.

:: 2. 测试登录
echo [33m2. 测试用户登录...[0m
echo ------------------------
curl -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"%USERNAME%\",\"password\":\"%PASSWORD%\"}" > %temp%\login_response.txt
type %temp%\login_response.txt
echo.
echo.
echo ------------------------
echo.

:: 从登录响应中提取token（使用findstr）
set TOKEN=
for /f "tokens=2 delims=:," %%a in ('type %temp%\login_response.txt ^| findstr /c:"\"token\""') do (
    set TOKEN=%%a
    goto :got_token
)
:got_token
set TOKEN=%TOKEN:"=%
set TOKEN=%TOKEN: =%
set TOKEN=%TOKEN:,=%

if "%TOKEN%"=="" (
    echo [31m❌ 无法获取token，请检查登录响应[0m
    pause
    exit /b 1
)

echo [32m✅ 成功获取token[0m
echo.
pause
echo.

:: 3. 测试获取用户信息
echo [33m3. 测试获取用户信息...[0m
echo ------------------------
curl -X GET http://localhost:8080/api/v1/profile ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.
echo ------------------------
echo.
pause
echo.

:: 4. 测试更新个人资料
echo [33m4. 测试更新个人资料...[0m
echo ------------------------
curl -X PUT http://localhost:8080/api/v1/profile ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"avatar\":\"https://example.com/avatar.jpg\",\"bio\":\"这是一个测试用户的个人简介\"}"
echo.
echo.
echo ------------------------
echo.
pause
echo.

:: 5. 测试获取存储空间信息
echo [33m5. 测试存储空间信息...[0m
echo ------------------------
curl -X GET http://localhost:8080/api/v1/storage ^
  -H "Authorization: Bearer %TOKEN%"
echo.
echo.
echo ------------------------
echo.
pause
echo.

:: 6. 测试修改密码
echo [33m6. 测试修改密码...[0m
echo ------------------------
curl -X POST http://localhost:8080/api/v1/change-password ^
  -H "Content-Type: application/json" ^
  -H "Authorization: Bearer %TOKEN%" ^
  -d "{\"old_password\":\"%PASSWORD%\",\"new_password\":\"newpass123\",\"confirm_password\":\"newpass123\"}"
echo.
echo.
echo ------------------------
echo.
pause
echo.

:: 7. 使用新密码登录
echo [33m7. 使用新密码登录验证...[0m
echo ------------------------
curl -X POST http://localhost:8080/api/v1/login ^
  -H "Content-Type: application/json" ^
  -d "{\"username\":\"%USERNAME%\",\"password\":\"newpass123\"}"
echo.
echo.
echo ------------------------
echo.

echo ===================================
echo [32m✅ 所有测试执行完毕！[0m
echo ===================================
echo.
echo 测试账号: %USERNAME%
echo 测试邮箱: %EMAIL%
echo 最终密码: newpass123
echo.
echo 你可以保留此账号用于后续测试。
echo.

:: 清理临时文件
del %temp%\login_response.txt 2>nul

pause