@echo off
chcp 65001 >nul
echo ====================================
echo       快速功能测试
echo ====================================
echo.

:: 注册
echo 1. 注册用户...
set /p newuser=请输入新用户名: 
curl -X POST http://localhost:8080/api/v1/register -H "Content-Type: application/json" -d "{\"username\":\"%newuser%\",\"email\":\"%newuser%@test.com\",\"password\":\"123456\",\"confirm_password\":\"123456\"}"
echo.
pause

:: 登录
echo 2. 登录...
curl -X POST http://localhost:8080/api/v1/login -H "Content-Type: application/json" -d "{\"username\":\"%newuser%\",\"password\":\"123456\"}"
echo.
pause

:: 创建测试文件
echo 3. 创建测试文件...
echo test > test.txt

:: 上传
echo 4. 上传文件...
set /p token=请输入上面获取的token: 
curl -X POST http://localhost:8080/api/v1/files/upload -H "Authorization: Bearer %token%" -F "file=@test.txt"
echo.
pause

:: 查看列表
echo 5. 查看文件列表...
curl -X GET "http://localhost:8080/api/v1/files" -H "Authorization: Bearer %token%"
echo.
pause

echo 测试完成！
pause