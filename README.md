# 吉米AI - 后端开发文档

## 项目简介

吉米AI是一个基于 Go 语言开发的AI应用后端系统，采用 Gin 框架构建，提供用户端（App）和管理端（Admin）两个应用。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: MySQL
- **缓存**: Redis
- **支付**: 支付宝、微信支付、Apple支付
- **第三方服务**: 
  - 环信（即时通讯）
  - 七牛云（对象存储）
  - 高德地图（地理位置）
  - 数美（内容审核）
  - 阿里云（短信服务）

## 项目结构

```
jimiaigo/
├── api/                    # API层 - 处理HTTP请求
│   ├── adminChargeApi/     # 管理端充值API
│   ├── adminHomeApi/       # 管理端首页API
│   ├── adminInviteApi/     # 管理端邀请API
│   ├── adminMemberApi/     # 管理端会员API
│   ├── adminRouterApi/     # 管理端路由API
│   ├── adminSettingApi/    # 管理端设置API
│   ├── adminTransactionApi/# 管理端交易API
│   ├── adminUserApi/       # 管理端用户API
│   ├── appVersionApi/      # 应用版本API
│   ├── chargeApi/          # 充值API
│   ├── feedbackApi/        # 反馈API
│   ├── globalApi/          # 全局API
│   ├── inviteApi/          # 邀请API
│   ├── payApi/             # 支付API
│   └── userApi/            # 用户API
│
├── cmd/                    # 应用入口
│   ├── admin/              # 管理端入口
│   │   └── main.go
│   └── app/                # 用户端入口
│       └── main.go
│
├── config/                 # 配置管理
│   ├── config.go           # 配置结构定义
│   └── config.ini          # 配置文件
│
├── cron/                   # 定时任务
│   └── cron.go
│
├── global/                 # 全局变量和初始化
│   ├── db.go               # 数据库连接
│   ├── global.go           # 全局变量
│   └── redis/              # Redis连接
│       └── redis.go
│
├── middleware/             # 中间件
│   ├── banned.go           # 封禁检查
│   ├── context.go          # 上下文处理
│   ├── cors.go             # 跨域处理
│   ├── jwt.go              # JWT认证
│   └── trace.go            # 链路追踪
│
├── model/                  # 数据模型层
│   ├── adminRouterModel/   # 管理端路由模型
│   ├── adminUserModel/     # 管理端用户模型
│   ├── appVersionModel/    # 应用版本模型
│   ├── chargeModel/        # 充值模型
│   ├── contextModel/       # 上下文模型
│   ├── globalModel/        # 全局模型
│   ├── inviteModel/        # 邀请模型
│   ├── memberModel/        # 会员模型
│   ├── payModel/           # 支付模型
│   ├── positionModel/      # 位置模型
│   ├── settingModel/       # 设置模型
│   ├── userModel/          # 用户模型
│   └── model.go            # 通用模型
│
├── repo/                   # 数据访问层
│   ├── adminReportRepo/    # 管理端举报仓库
│   ├── adminRouterRepo/    # 管理端路由仓库
│   ├── adminUserRepo/      # 管理端用户仓库
│   ├── appVersionRepo/     # 应用版本仓库
│   ├── chargeRepo/         # 充值仓库
│   ├── feedbackRepo/       # 反馈仓库
│   ├── globalRepo/         # 全局仓库
│   ├── inviteRepo/         # 邀请仓库
│   ├── matchRepo/          # 匹配仓库
│   ├── memberRepo/         # 会员仓库
│   ├── payRepo/            # 支付仓库
│   ├── reportRepo/         # 举报仓库
│   ├── serverRepo/         # 服务器仓库
│   ├── settingRepo/        # 设置仓库
│   ├── userRepo/           # 用户仓库
│   └── versionRepo/        # 版本仓库
│
├── router/                 # 路由配置
│   ├── adminRouter.go      # 管理端路由
│   └── appRouter.go       # 用户端路由
│
├── service/                # 业务逻辑层
│   ├── adminInviteService/ # 管理端邀请服务
│   ├── adminRouterService/ # 管理端路由服务
│   ├── adminUserService/   # 管理端用户服务
│   ├── appVersionService/  # 应用版本服务
│   ├── chargeService/      # 充值服务
│   ├── feedbackService/    # 反馈服务
│   ├── inviteService/      # 邀请服务
│   ├── memberService/      # 会员服务
│   ├── payService/         # 支付服务
│   ├── settingService/     # 设置服务
│   ├── userService/        # 用户服务
│   └── versionService/     # 版本服务
│
├── util/                   # 工具类
│   ├── dypns/              # 阿里云一键登录
│   ├── dysms/              # 阿里云短信
│   ├── gaode/              # 高德地图
│   ├── http/               # HTTP工具
│   ├── hx/                 # 环信工具
│   ├── ishumei/            # 数美工具
│   ├── jpush/              # 极光推送
│   ├── pay/                # 支付工具
│   │   ├── alipay/         # 支付宝
│   │   └── wechat/         # 微信支付
│   ├── qiniu/              # 七牛云
│   ├── response/           # 响应工具
│   ├── jwt.go              # JWT工具
│   ├── md5.go              # MD5工具
│   └── util.go             # 通用工具
│
├── runtime/                # 运行时文件
│   ├── log/                # 日志文件
│   └── upload/             # 上传文件
│
├── test/                   # 测试文件
│
├── go.mod                  # Go模块定义
├── go.sum                  # Go依赖校验
├── Makefile                # 构建脚本
└── README.md               # 项目文档
```

## 架构设计

项目采用**分层架构**设计，主要分为以下层次：

1. **API层** (`api/`): 处理HTTP请求，参数验证，调用Service层
2. **Service层** (`service/`): 业务逻辑处理，调用Repo层
3. **Repo层** (`repo/`): 数据访问层，封装数据库操作
4. **Model层** (`model/`): 数据模型定义，对应数据库表结构

### 数据流向

```
HTTP请求 → API层 → Service层 → Repo层 → Database
                ↓
            Response ← API层 ← Service层 ← Repo层
```

## 数据库模型

### 核心业务模型

#### 1. 用户模型 (User)

**表名**: `users`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Phone | string | 手机号（唯一） |
| Avatar | string | 头像URL |
| City | string | 城市 |
| Nickname | string | 昵称 |
| Birthday | string | 出生日期 |
| Weight | int | 体重 |
| Tags | json | 用户标签 |
| Longitude | float64 | 经度 |
| Latitude | float64 | 纬度 |
| LastRequestTime | datetime | 最后请求时间 |
| VipTime | datetime | VIP到期时间 |
| SumVipDays | int | 累积会员天数 |
| HotBanned | bool | 是否被热度封禁 |
| HotConstValue | int | 热度常量值 |
| BanLevel | int | 封禁等级（0=不封禁, 1=禁言一天, 2=禁言一周, 3=禁言一个月, 4=禁言1年, 5=永久禁言, 6=永久封号） |
| BanTime | datetime | 封禁到期时间 |
| BanRecordId | int | 封禁记录ID |
| LastLoginIp | string | 最后登录IP |
| wallet_income | decimal | 钱包收入 |
| wallet_balance | decimal | 钱包余额 |
| wallet_withdrawal | decimal | 钱包提现 |
| payment_diamond | int | 钻石 |
| payment_jifen | int | 积分 |
| payment_jiyu_coin | int | 机遇币 |
| payment_accounts | json | 支付账号列表 |
| RegisterDevice | string | 注册设备 |
| RefId | int | 推荐人ID |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |
| DeletedAt | datetime | 删除时间 |

**关联模型**:
- `UserLike`: 用户喜好设置
- `UserBeLike`: 用户被喜好设置
- `Wallet`: 钱包信息（嵌入）
- `Payment`: 支付信息（嵌入）
- `Position`: 位置信息（嵌入）

#### 2. 关注模型 (FollowLog)

**表名**: `follow_logs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | int | 用户ID |
| FollowUserID | int | 关注用户ID |
| CreatedAt | datetime | 创建时间 |

**索引**: 
- `idx_user_id`
- `idx_follow_user_id`
- `idx_user_id_follow_user_id` (唯一索引)

#### 3. 好友模型 (Friends)

**表名**: `friends`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | int | 用户ID |
| FriendUserID | int | 好友用户ID |
| CreatedAt | datetime | 创建时间 |

**索引**:
- `idx_user_id`
- `idx_friend_user_id`
- `idx_user_id_friend_user_id` (唯一索引)

#### 4. 充值配置模型 (ChargeInfo)

**表名**: `charge_infos`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Price | int | 价格（分） |
| OriginPrice | int | 原价（分） |
| Month | int | 月数（1=1个月, 3=3个月, 12=年会员） |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

#### 5. 支付日志模型 (PaymentLog)

**表名**: `payment_logs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Uid | int | 用户ID |
| Nickname | string | 昵称 |
| Rmb | float64 | 人民币金额 |
| Diamond | int | 钻石数量 |
| FreeDiamond | int | 赠送钻石 |
| OrderNo | string | 商户订单号 |
| TradeNo | string | 第三方订单号 |
| PaymentType | string | 支付类型（alipay/wechat/apple） |
| PaymentEnv | string | 支付环境（app/web/h5） |
| FreeJifen | int | 赠送积分 |
| FreeVipDays | int | 赠送VIP天数 |
| FreeGift | int | 赠送礼物 |
| OrderStatus | int | 订单状态（0=未支付, 1=已支付, 2=失败） |
| Remark | string | 备注 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**: `idx_uid`

#### 6. 提现日志模型 (WithdrawalLog)

**表名**: `withdrawal_logs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Uid | int | 用户ID |
| Nickname | string | 昵称 |
| Rmb | float64 | 提现金额 |
| Ticket | int | 兑换映票 |
| WithdrawalType | string | 提现类型（alipay/wechat） |
| WithdrawalName | string | 提现姓名 |
| WithdrawalAccount | string | 提现账号 |
| Status | int | 提现状态（0=未处理, 1=已处理, 2=已拒绝） |
| HandleTime | datetime | 处理时间 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**: `idx_uid`

#### 7. 会员配置模型 (MemberConfig)

**表名**: `member_configs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Type | int | 会员类型（1=青铜, 2=白银, 3=黄金, 4=钻石） |
| Name | string | 会员名称 |
| Icon | string | 会员图标 |
| MinDays | int | 最小天数 |
| MaxDays | int | 最大天数 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

#### 8. 邀请记录模型 (InviteLog)

**表名**: `invite_logs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserID | int64 | 用户ID |
| InviteUserID | int64 | 被邀请用户ID |
| InviteType | int | 邀请类型（0=免费, 1=VIP） |
| InviteCoins | int | 邀请奖励币 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

#### 9. 举报模型 (Report)

**表名**: `reports`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserId | int | 举报人ID |
| Nickname | string | 举报人昵称 |
| Phone | string | 举报人手机号 |
| TargetUid | int | 被举报用户ID |
| TargetNickname | string | 被举报用户昵称 |
| Sign | string | 个性签名 |
| Avatar | string | 头像 |
| BgImage | string | 背景图 |
| Content | text | 举报文字内容 |
| ReportScene | int | 举报场景（0=漂流瓶, 1=动态, 2=消息, 3=个人消息） |
| ReportType | int | 举报类型（0=低俗色情, 1=垃圾广告, 2=骚扰/不文明, 3=涉嫌欺诈, 4=政治, 5=恐暴, 6=其他） |
| Detail | text | 举报补充说明 |
| Process | int | 处理状态（0=未处理, 1=已处罚, 2=未违规） |
| ProcessTime | datetime | 处理时间 |
| ProcessType | int | 处理类型（0=系统处理, 1=人工处理） |
| BanLevel | int | 封禁等级 |
| BanTime | datetime | 封禁到期时间 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**:
- `idx_user_id`
- `idx_phone`
- `idx_target_uid`

#### 10. 反馈模型 (Feedback)

**表名**: `feedbacks`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserId | int | 用户ID |
| Nickname | string | 反馈人昵称 |
| SystemVersion | string | 系统版本 |
| PhoneModel | string | 手机型号 |
| FeedbackType | int | 反馈类型 |
| Content | text | 反馈内容 |
| Status | int | 处理状态（0=未处理, 1=已处理） |
| ProcessTime | datetime | 处理时间 |
| Replay | text | 回复内容 |
| IsRead | bool | 是否已读 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**: `idx_user_id`

#### 11. 常见问题模型 (QA)

**表名**: `qas`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Question | string | 问题 |
| Answer | text | 答案 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

#### 12. 应用版本模型 (AppVersion)

**表名**: `app_versions`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Version | string | 版本号（唯一） |
| BuildVersion | string | 构建版本号 |
| Platform | string | 平台（ios/android） |
| Status | bool | 是否强制更新 |
| Url | string | 下载地址 |
| Content | text | 更新内容 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**: `idx_version` (唯一索引)

#### 13. 系统设置模型 (Setting)

**表名**: `settings`

包含多个嵌入配置：
- `Login`: 登录设置（短信验证码配置）
- `Withdraw`: 提现设置
- `IM`: 即时通讯设置
- `Payment`: 支付设置
- `Invite`: 邀请设置
- `App`: 应用设置
- `Report`: 举报设置
- `HX`: 环信设置
- `Aliyun`: 阿里云设置
- `Qiniu`: 七牛云设置
- `Match`: 匹配设置
- `Gaode`: 高德地图设置
- `PushChat`: 推送聊天设置
- `Bottle`: 漂流瓶设置
- `HotUser`: 热门用户设置
- `Shumei`: 数美设置
- `Ad`: 广告设置
- `Member`: 会员设置

#### 14. 文档模型 (Doc)

**表名**: `docs`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UserAgreement | mediumtext | 用户协议 |
| PrivacyPolicy | mediumtext | 隐私政策 |
| MemberAgreement | mediumtext | 会员协议 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

#### 15. 短信验证码模型 (SMSCode)

**表名**: `sms_codes`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Phone | string | 手机号 |
| Code | string | 验证码 |
| CreatedAt | datetime | 创建时间 |

**索引**: `idx_phone`

#### 16. 推荐ID-IP关联模型 (RefIdIp)

**表名**: `ref_id_ips`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| RefId | int | 推荐ID |
| Ip | string | IP地址 |
| CreatedAt | datetime | 创建时间 |

**索引**:
- `idx_ref_id`
- `idx_ip`

#### 17. 图片模型 (Image)

**表名**: `images`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| UUID | string | 图片UUID（唯一） |
| Width | int | 宽度 |
| Height | int | 高度 |
| CreatedAt | datetime | 创建时间 |

**索引**: `idx_uuid` (唯一索引)

### 管理端模型

#### 18. 管理员用户模型 (AdminUser)

**表名**: `admin_users`

| 字段 | 类型 | 说明 |
|------|------|------|
| ID | uint | 主键 |
| Username | string | 用户名（唯一） |
| Password | string | 密码（加密） |
| Role | int | 角色（0=管理员, 1=超级管理员） |
| Avatar | string | 头像 |
| Name | string | 姓名 |
| Introduction | string | 简介 |
| CreatedAt | datetime | 创建时间 |
| UpdatedAt | datetime | 更新时间 |

**索引**: `idx_username` (唯一索引)

#### 19. 管理端路由模型 (Router)

**表名**: `routers`

用于管理前端路由权限配置。

## 开发指南

### 环境要求

- Go 1.21+
- MySQL 5.7+
- Redis 6.0+

### 配置文件

配置文件位于 `config/config.ini`，需要配置以下部分：

```ini
[Database]
DbHost = localhost
DbPort = 3306
UserName = root
Password = password
DBName = jiyu
TimeZone = Asia/Shanghai
TablePrefix = jiyu_

[Redis]
Addr = localhost:6379
Pwd = 
DB = 0

[Local]
IsMaster = true
RunPort = :8080

[AdminApp]
RunPort = :8081
```

### 启动项目

#### 用户端应用

```bash
# 编译
go build -o ./bin/app ./cmd/app/main.go

# 运行
./bin/app
```

#### 管理端应用

```bash
# 编译
go build -o ./bin/admin ./cmd/admin/main.go

# 运行
./bin/admin
```

### 数据库迁移

项目启动时会自动执行数据库迁移，所有模型定义在 `global/db.go` 的 `AutoMigrate()` 函数中。

### 中间件说明

1. **Cors**: 跨域处理
2. **Trace**: 链路追踪，为每个请求生成唯一TraceID
3. **JWTAuth**: JWT认证，验证用户token
4. **Context**: 上下文处理，注入用户信息
5. **Banned**: 封禁检查，验证用户是否被封禁

### API路由

#### 用户端路由 (`router/appRouter.go`)

- `/api/ping`: 服务响应测试
- `/api/login`: 登录
- `/api/user/*`: 用户相关接口
- `/api/charge/*`: 充值相关接口
- `/api/payment/*`: 支付相关接口
- `/api/feedback/*`: 反馈相关接口
- `/api/invite/*`: 邀请相关接口

#### 管理端路由 (`router/adminRouter.go`)

- `/api/login`: 管理员登录
- `/api/user/*`: 用户管理接口
- `/api/setting/*`: 系统设置接口
- `/api/charge/*`: 充值配置接口
- `/api/member/*`: 会员管理接口
- `/api/invite/*`: 邀请管理接口

### 开发规范

1. **命名规范**:
   - 包名使用小写
   - 结构体使用驼峰命名
   - 常量使用大写

2. **代码组织**:
   - API层只负责参数验证和调用Service
   - Service层处理业务逻辑
   - Repo层只负责数据库操作
   - Model层定义数据模型

3. **错误处理**:
   - 使用 `errors` 包处理错误
   - 统一使用 `util/response` 返回响应

4. **日志记录**:
   - 日志文件位于 `runtime/log/`
   - 使用标准库 `log` 记录日志

## 部署说明

### 构建

```bash
# 构建用户端
go build -o ./bin/app ./cmd/app/main.go

# 构建管理端
go build -o ./bin/admin ./cmd/admin/main.go
```

### 运行

```bash
# 后台运行用户端
nohup ./bin/app > out.log 2>&1 &

# 后台运行管理端
nohup ./bin/admin > admin.log 2>&1 &
```

### 注意事项

1. 确保数据库和Redis服务已启动
2. 配置文件路径正确
3. 日志目录有写入权限
4. 上传目录有写入权限

## 主要功能模块

### 用户模块
- 用户注册/登录
- 用户信息管理
- 关注/好友功能
- 用户封禁管理

### 充值模块
- 充值配置管理
- 订单创建
- 支付回调处理（支付宝/微信/Apple）

### 会员模块
- 会员等级配置
- VIP权限管理
- 会员等级计算

### 支付模块
- 支付日志记录
- 提现申请
- 提现审核

### 邀请模块
- 邀请记录
- 邀请排行榜
- 邀请奖励

### 反馈模块
- 意见反馈
- 常见问题管理

### 举报模块
- 用户举报
- 举报处理

## 许可证

[根据实际情况填写]

## 联系方式

[根据实际情况填写]

