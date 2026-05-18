# 汇率兑换应用 (Currency Exchange App)

Go + Gin + GORM + MySQL + Redis 全栈实战项目

## 技术栈

**后端：**
- Go 1.26
- [Gin](https://github.com/gin-gonic/gin) — HTTP Web 框架
- [GORM](https://gorm.io/) — ORM 框架
- [Viper](https://github.com/spf13/viper) — 配置管理
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — 密码哈希
- [JWT](https://github.com/golang-jwt/jwt) — 身份认证令牌
- [go-redis](https://github.com/go-redis/redis) — Redis 客户端
- [gin-contrib/cors](https://github.com/gin-contrib/cors) — CORS 跨域中间件

**前端：**
- Vue 3 + Vite（开发服务器默认 `localhost:5173`）

**数据库：**
- MySQL
- Redis

## 项目结构

```
fullstack/
├── backend/
│   ├── main.go                   # 入口文件
│   ├── config/                   # 配置相关
│   │   ├── config.go             # 配置结构体 & Viper 加载
│   │   ├── config.yml            # YAML 配置文件
│   │   ├── db.go                 # MySQL 连接初始化
│   │   └── redis.go              # Redis 连接初始化
│   ├── router/                   # 路由
│   │   └── router.go             # 路由注册 & CORS 配置
│   ├── middlewares/              # 中间件
│   │   └── auth_middleware.go    # JWT 认证中间件
│   ├── controllers/              # 控制器（处理请求）
│   │   ├── auth_controller.go        # 用户认证（注册/登录）
│   │   ├── auth_rate_controller.go   # 汇率数据接口
│   │   ├── article_controller.go     # 文章 CRUD + 旁路缓存
│   │   └── like_controller.go        # 点赞功能（Redis 计数）
│   ├── models/                   # 数据模型
│   │   ├── user.go               # 用户模型
│   │   ├── exchange_rate.go      # 汇率模型
│   │   └── article.go            # 文章模型
│   ├── global/                   # 全局变量
│   │   └── global.go             # MySQL & Redis 连接实例
│   └── utils/                    # 工具函数
│       └── utils.go              # 密码哈希、JWT 生成与验证
```

## 快速开始

### 1. 环境准备

确保已安装并启动 MySQL 和 Redis。

### 2. 创建数据库

```sql
CREATE DATABASE fullstack CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 修改配置

编辑 `backend/config/config.yml`，根据你的环境修改数据库连接信息：

```yaml
database:
  dsn: root:你的密码@tcp(localhost:3306)/fullstack?charset=utf8mb4&parseTime=True&loc=Local
```

Redis 配置在 `backend/config/redis.go` 中，默认连接 `localhost:6379`，无需密码。

### 4. 安装依赖

```bash
cd backend
go mod tidy
```

### 5. 运行

```bash
go run main.go
```

服务默认运行在 `http://localhost:3000`

### 6. 前后端联调

后端已配置 CORS，允许前端 `http://localhost:5173` 跨域请求（含 Cookie/Token）。

## API

### 用户认证

| 方法 | 路径              | 说明     | 认证 |
|------|-------------------|----------|------|
| POST | /api/auth/register| 用户注册 | 否   |
| POST | /api/auth/login   | 用户登录 | 否   |

### 汇率管理

| 方法 | 路径              | 说明     | 认证 |
|------|-------------------|----------|------|
| GET  | /api/exchangerate | 查询全部 | 否   |
| POST | /api/exchangerate | 新增汇率 | 是   |

### 文章管理

| 方法 | 路径                   | 说明       | 认证 |
|------|-----------------------|------------|------|
| POST | /api/articles          | 创建文章   | 是   |
| GET  | /api/articles          | 查询全部   | 是   |
| GET  | /api/articles/:id      | 查询单篇   | 是   |

### 点赞功能（Redis 原子计数）

| 方法 | 路径                     | 说明       | 认证 |
|------|-------------------------|------------|------|
| POST | /api/articles/:id/like   | 点赞文章   | 是   |
| GET  | /api/articles/:id/like   | 查询点赞数 | 是   |

## 架构说明

- **缓存策略**：文章列表采用旁路缓存（Cache-Aside），先查 Redis → 未命中则查 MySQL → 回填 Redis（TTL 10 分钟）
- **点赞计数**：使用 Redis `INCR` 原子操作，并发安全
- **身份认证**：JWT（HS256），登录/注册后返回 token，后续请求放 `Authorization: Bearer <token>` 头
- **CORS**：开发环境允许 `localhost:5173` 跨域，生产环境交给 Nginx 处理

## 请求示例

### 注册

```bash
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"123456"}'
```

### 登录

```bash
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"123456"}'
```

### 新增汇率

```bash
curl -X POST http://localhost:3000/api/exchangerate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"fromCurrency":"USD","toCurrency":"CNY","rate":7.25}'
```

### 创建文章

```bash
curl -X POST http://localhost:3000/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"我的第一篇文章","content":"文章正文","preview":"文章摘要"}'
```

### 点赞文章

```bash
curl -X POST http://localhost:3000/api/articles/1/like \
  -H "Authorization: Bearer <token>"
```

### 查询文章点赞数

```bash
curl -X GET http://localhost:3000/api/articles/1/like \
  -H "Authorization: Bearer <token>"
```
