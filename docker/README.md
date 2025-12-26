# Docker 部署说明

本文档说明如何使用 Docker 和 Docker Compose 部署应用。

## 前置要求

- Docker 20.10 或更高版本
- Docker Compose 2.0 或更高版本

## 快速开始

### 1. 配置环境

在 `config/config.docker.ini` 文件中配置你的应用设置，特别是：
- 数据库连接信息（默认已配置为 Docker 服务名）
- Redis 连接信息（默认已配置为 Docker 服务名）
- 其他业务配置（如七牛云、邮件等）

### 2. 启动服务

```bash
# 构建并启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 3. 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止服务并删除数据卷（警告：会删除数据库数据）
docker-compose down -v
```

## 服务说明

### MySQL
- 容器名: `jimiai-mysql`
- 端口: `3306`
- 默认数据库: `jimiai`
- 默认用户名: `root`
- 默认密码: `123456`
- 数据持久化: `mysql_data` 数据卷

### Redis
- 容器名: `jimiai-redis`
- 端口: `6379`
- 默认密码: `123456`
- 数据持久化: `redis_data` 数据卷

### App (主应用)
- 容器名: `jimiai-app`
- 端口: `27355`
- 配置文件: `config/config.docker.ini`

### Admin (管理后台)
- 容器名: `jimiai-admin`
- 端口: `27356`
- 配置文件: `config/config.docker.ini`

## 数据持久化

以下数据会持久化存储：

1. **MySQL 数据**: `mysql_data` 数据卷
2. **Redis 数据**: `redis_data` 数据卷
3. **应用日志**: `./runtime/log` 目录
4. **上传文件**: `./runtime/upload` 目录

## 常见操作

### 查看日志

```bash
# 查看所有服务日志
docker-compose logs -f

# 查看特定服务日志
docker-compose logs -f app
docker-compose logs -f admin
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 进入容器

```bash
# 进入 MySQL 容器
docker exec -it jimiai-mysql bash

# 进入 Redis 容器
docker exec -it jimiai-redis sh

# 进入应用容器
docker exec -it jimiai-app sh
```

### 数据库操作

```bash
# 连接 MySQL
docker exec -it jimiai-mysql mysql -uroot -p123456 jimiai

# 连接 Redis
docker exec -it jimiai-redis redis-cli -a 123456
```

### 重启服务

```bash
# 重启所有服务
docker-compose restart

# 重启特定服务
docker-compose restart app
```

### 重建镜像

```bash
# 重新构建并启动
docker-compose up -d --build
```

## 生产环境部署

生产环境建议使用 `docker-compose.prod.yml`，并配置环境变量：

```bash
# 创建 .env 文件
cat > .env << EOF
MYSQL_ROOT_PASSWORD=your_strong_password
MYSQL_DATABASE=jimiai
MYSQL_USER=jimi
MYSQL_PASSWORD=your_strong_password
REDIS_PASSWORD=your_strong_redis_password
APP_PORT=27355
ADMIN_PORT=27356
MYSQL_PORT=3306
REDIS_PORT=6379
EOF

# 使用生产配置启动
docker-compose -f docker-compose.prod.yml up -d
```

## 注意事项

1. **首次启动**: MySQL 容器需要一些时间来初始化，应用容器会等待 MySQL 就绪后再启动
2. **配置文件**: Docker 环境使用 `config.docker.ini`，确保数据库和 Redis 主机名使用服务名（`mysql`、`redis`）
3. **端口冲突**: 确保主机上的 3306、6379、27355、27356 端口未被占用
4. **数据备份**: 定期备份 `mysql_data` 和 `redis_data` 数据卷
5. **日志管理**: 定期清理 `./runtime/log` 目录下的日志文件，避免占用过多磁盘空间

## 故障排查

### 应用无法连接数据库

1. 检查 MySQL 容器是否正常运行: `docker-compose ps mysql`
2. 检查网络连接: `docker network inspect jimiai_jimiaigo-network`
3. 查看 MySQL 日志: `docker-compose logs mysql`

### 应用无法连接 Redis

1. 检查 Redis 容器是否正常运行: `docker-compose ps redis`
2. 检查 Redis 密码是否正确
3. 查看 Redis 日志: `docker-compose logs redis`

### 查看应用日志

```bash
# 查看应用日志
docker-compose logs app
docker-compose logs admin

# 或者查看日志文件
tail -f runtime/log/app.log
tail -f runtime/log/admin.log
```

