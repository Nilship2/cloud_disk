# Cloud Disk 云盘系统

一个仿百度网盘的后端系统，支持文件存储、分享、收藏等功能。

## 功能特性
- 用户注册登录（JWT认证）
- 文件上传下载（支持本地/MinIO存储）
- 文件分享（生成分享链接）
- 文件收藏管理
- 回收站功能
- 存储空间管理

## 技术栈
- **语言**: Go 1.21+
- **框架**: Gin + GORM
- **数据库**: MySQL 8.0 + Redis
- **存储**: 本地/MinIO/S3
- **部署**: Docker + Docker Compose

## 快速开始

### 环境要求
- Go 1.21+
- MySQL 8.0+
- Redis 6.0+
- (可选) MinIO

### 安装
```bash
# 克隆项目
git clone <repository-url>

# 进入目录
cd cloud-disk

# 安装依赖
go mod download

# 配置环境变量
cp .env.example .env

# 启动服务
make run