# GitHub Actions CI/CD 配置说明

本项目配置了完整的 CI/CD 工作流，包括持续集成（CI）和持续部署（CD）。

## 工作流文件

### 1. `ci.yml` - 持续集成
**触发条件**：
- 推送到 `main` 或 `develop` 分支
- 创建 Pull Request 到 `main` 或 `develop` 分支

**执行任务**：
- **Lint**: 代码静态检查（使用 golangci-lint）
- **Test**: 运行单元测试并生成覆盖率报告
- **Build**: 构建应用二进制文件

### 2. `cd.yml` - 持续部署
**触发条件**：
- 推送到 `main` 分支
- 创建版本标签（`v*`）
- 手动触发（workflow_dispatch）

**执行任务**：
- **Build and Push**: 构建 Docker 镜像并推送到 GitHub Container Registry
- **Deploy**: 自动部署到生产服务器（仅 main 分支）

### 3. `docker-build.yml` - Docker 构建测试
**触发条件**：
- 推送到 `main` 或 `develop` 分支
- 创建 Pull Request

**执行任务**：
- 构建 Docker 镜像（不推送）
- 测试镜像是否正常构建

## 配置要求

### GitHub Secrets 配置

在 GitHub 仓库设置中添加以下 Secrets：

#### 部署相关（用于 CD）
- `DEPLOY_HOST`: 服务器 IP 地址或域名
- `DEPLOY_USER`: SSH 用户名
- `DEPLOY_SSH_KEY`: SSH 私钥（用于服务器登录）
- `DEPLOY_PORT`: SSH 端口（可选，默认 22）

#### 其他 Secrets（可选）
- `CODECOV_TOKEN`: Codecov token（如果使用 Codecov 上传覆盖率）

### 配置步骤

1. **进入 GitHub 仓库设置**
   - 点击仓库的 `Settings` → `Secrets and variables` → `Actions`

2. **添加 Secrets**
   ```
   DEPLOY_HOST=your-server-ip
   DEPLOY_USER=deploy
   DEPLOY_SSH_KEY=-----BEGIN OPENSSH PRIVATE KEY-----
   ...
   -----END OPENSSH PRIVATE KEY-----
   DEPLOY_PORT=22
   ```

3. **生成 SSH 密钥对**（如果还没有）
   ```bash
   ssh-keygen -t ed25519 -C "github-actions" -f ~/.ssh/github_actions
   ```
   
   将公钥添加到服务器：
   ```bash
   ssh-copy-id -i ~/.ssh/github_actions.pub deploy@your-server-ip
   ```

4. **设置服务器权限**
   确保部署用户有权限执行 docker 和 docker-compose：
   ```bash
   sudo usermod -aG docker $USER
   ```

## 使用说明

### 自动触发

- **推送到 main 分支**：自动运行 CI，然后构建并推送 Docker 镜像，最后部署到服务器
- **推送到 develop 分支**：只运行 CI 和 Docker 构建测试
- **创建 Pull Request**：只运行 CI 和 Docker 构建测试

### 手动触发部署

1. 进入 GitHub Actions 页面
2. 选择 `CD` 工作流
3. 点击 `Run workflow`
4. 输入 Docker 镜像标签（默认：latest）
5. 点击 `Run workflow` 按钮

### 版本发布

创建版本标签会自动触发部署：

```bash
git tag -a v1.0.0 -m "Release version 1.0.0"
git push origin v1.0.0
```

这会：
1. 构建 Docker 镜像
2. 使用版本号作为标签（如 `v1.0.0`）
3. 推送到 GitHub Container Registry

## Docker 镜像

### 镜像位置

镜像会推送到 GitHub Container Registry：
```
ghcr.io/your-username/jimiaigo:latest
ghcr.io/your-username/jimiaigo:main
ghcr.io/your-username/jimiaigo:v1.0.0
```

### 拉取镜像

```bash
docker pull ghcr.io/your-username/jimiaigo:latest
```

### 使用镜像

在 `docker-compose.prod.yml` 中，可以修改为使用远程镜像：

```yaml
app:
  image: ghcr.io/your-username/jimiaigo:latest
  # 或者使用 build，本地构建
  # build:
  #   context: .
  #   dockerfile: Dockerfile
```

## 工作流状态

- ✅ **绿色**：所有检查通过
- ❌ **红色**：有检查失败，需要修复
- 🟡 **黄色**：正在运行中

## 故障排查

### CI 失败

1. **Lint 失败**：检查代码规范，运行 `golangci-lint run` 本地检查
2. **Test 失败**：检查测试用例，确保所有测试通过
3. **Build 失败**：检查 Go 代码编译错误

### CD 失败

1. **镜像构建失败**：检查 Dockerfile 是否正确
2. **推送失败**：检查 GitHub Token 权限
3. **部署失败**：
   - 检查 SSH 连接是否正常
   - 检查服务器上 docker-compose 是否可用
   - 检查服务器磁盘空间
   - 查看服务器日志：`journalctl -u docker`

## 最佳实践

1. **分支策略**：
   - `main`: 生产环境，自动部署
   - `develop`: 开发环境，只运行 CI

2. **代码审查**：
   - 所有代码变更通过 Pull Request
   - 确保 CI 通过后再合并

3. **版本管理**：
   - 使用语义化版本（Semantic Versioning）
   - 重要版本创建 Git 标签

4. **安全**：
   - 不要将 Secrets 提交到代码仓库
   - 定期轮换 SSH 密钥
   - 使用最小权限原则

## 相关文档

- [GitHub Actions 文档](https://docs.github.com/en/actions)
- [Docker 文档](https://docs.docker.com/)
- [项目部署文档](../doc/deploy.md)

