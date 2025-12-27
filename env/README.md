# 环境变量配置文件说明

## 目录结构

在生产服务器上，环境变量文件应放置在以下位置：

```
/opt/project/
├─ docker-compose.yml
├─ env/
│  ├─ frontend.env      # 前端服务环境变量（可选）
│  ├─ backend.env       # 后端服务环境变量
│  ├─ mysql.env         # MySQL 数据库环境变量
│  └─ redis.env         # Redis 环境变量
└─ data/
   ├─ mysql/            # MySQL 数据目录
   └─ redis/            # Redis 数据目录
```

## 使用方法

### 1. 创建环境变量文件

将项目中的示例文件复制到服务器对应位置：

```bash
# 在服务器上执行
mkdir -p /opt/project/env
mkdir -p /opt/project/data/mysql
mkdir -p /opt/project/data/redis

# 复制示例文件并重命名（去掉 .example 后缀）
cp env/mysql.env.example /opt/project/env/mysql.env
cp env/redis.env.example /opt/project/env/redis.env
cp env/backend.env.example /opt/project/env/backend.env
cp env/frontend.env.example /opt/project/env/frontend.env  # 如果有前端服务
```

### 2. 修改环境变量

编辑各个 env 文件，根据实际环境修改配置值：

- **mysql.env**: MySQL 数据库相关配置
- **redis.env**: Redis 相关配置
- **backend.env**: 后端应用相关配置（包括数据库连接、Redis 连接、API 配置等）
- **frontend.env**: 前端服务相关配置（如果有）

### 3. 设置文件权限

确保环境变量文件权限安全：

```bash
chmod 600 /opt/project/env/*.env
```

### 4. 启动服务

在 `/opt/project/` 目录下执行：

```bash
docker-compose -f docker-compose.prod.yml up -d
```

## 文件说明

### mysql.env
MySQL 数据库服务的环境变量，包括：
- `MYSQL_ROOT_PASSWORD`: root 用户密码
- `MYSQL_DATABASE`: 数据库名称
- `MYSQL_USER`: 普通用户名称
- `MYSQL_PASSWORD`: 普通用户密码
- `MYSQL_PORT`: 端口映射

### redis.env
Redis 服务的环境变量，包括：
- `REDIS_PORT`: 端口映射
- `REDIS_PASSWORD`: Redis 密码

### backend.env
后端应用服务的环境变量，包括：
- `APP_PORT`: 应用端口映射
- `DB_*`: 数据库连接配置
- `REDIS_*`: Redis 连接配置
- `APP_DYU_API_*`: API 相关配置
- `QINIU_*`: 七牛云配置
- `EMAIL_*`: 邮箱配置
- 等等...

### frontend.env
前端服务环境变量（如果项目中有前端服务）

## 注意事项

1. **安全性**: 生产环境请务必修改所有默认密码
2. **备份**: 建议定期备份 `/opt/project/env/` 目录
3. **权限**: env 文件不应提交到版本控制系统
4. **数据持久化**: 数据目录 `/opt/project/data/` 使用主机路径，确保数据持久化

