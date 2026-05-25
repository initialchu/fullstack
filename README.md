# 全栈实战项目 (Fullstack App)

Go + Gin + GORM + MySQL + Redis + Vue 3 全栈实战项目，涵盖用户认证、汇率查询、文章管理和点赞功能。

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
- Vue 3（Composition API + `<script setup>`）
- TypeScript
- [Vite](https://vite.dev/) — 构建工具 & 开发服务器（默认 `localhost:5173`）
- [Pinia](https://pinia.vuejs.org/) — 状态管理
- [Element Plus](https://element-plus.org/) — UI 组件库
- [Axios](https://axios-http.com/) — HTTP 请求库

**数据库 & 缓存：**
- MySQL — 持久化存储
- Redis — 缓存 + 点赞计数

## 项目结构

```
fullstack/
├── backend/
│   ├── main.go                   # 入口文件（优雅关闭）
│   ├── config/                   # 配置
│   │   ├── config.go             # 配置结构体 & Viper 加载
│   │   ├── config.yml            # YAML 配置文件
│   │   ├── db.go                 # MySQL 连接初始化
│   │   └── redis.go              # Redis 连接初始化
│   ├── router/                   # 路由
│   │   └── router.go             # 路由注册 & CORS 配置
│   ├── middlewares/              # 中间件
│   │   └── auth_middleware.go    # JWT 认证中间件
│   ├── controllers/              # 控制器
│   │   ├── auth_controller.go        # 用户认证（注册/登录）
│   │   ├── auth_rate_controller.go   # 汇率数据接口
│   │   ├── article_controller.go     # 文章 CRUD + 旁路缓存
│   │   └── like_controller.go        # 点赞功能（Redis INCR）
│   ├── models/                   # 数据模型
│   │   ├── user.go               # 用户模型
│   │   ├── exchange_rate.go      # 汇率模型
│   │   └── article.go            # 文章模型
│   ├── global/                   # 全局变量
│   │   └── global.go             # DB & Redis 连接实例
│   └── utils/                    # 工具函数
│       └── utils.go              # Bcrypt 密码哈希、JWT 生成/验证
├── frontend/
│   ├── index.html                # HTML 入口
│   ├── package.json              # 前端依赖
│   ├── vite.config.ts            # Vite 配置
│   ├── tsconfig.json             # TypeScript 配置
│   └── src/
│       ├── main.ts               # Vue 应用入口（挂载 Element Plus + Router）
│       ├── App.vue               # 根组件（顶部导航 + router-view）
│       ├── axios.ts              # Axios 实例（baseURL + JWT 拦截器）
│       ├── style.css             # 全局样式
│       ├── shims-vue.d.ts        # Vue 类型声明
│       ├── router/
│       │   └── index.ts          # 路由配置（7 个路由）
│       ├── store/
│       │   └── auth.ts           # Pinia 认证状态（token/登录/注册/退出）
│       ├── views/                # 页面级组件
│       │   ├── Home.vue          # 首页（欢迎页）
│       │   ├── CurrencyExchange.vue  # 货币兑换
│       │   ├── News.vue          # 文章列表（含点赞交互）
│       │   └── Detail.vue        # 文章详情
│       ├── components/           # 功能组件
│       │   ├── Login.vue         # 登录表单
│       │   └── Register.vue      # 注册表单
│       └── assets/               # 静态资源
│           └── like.png          # 点赞图标
└── .gitignore
```

## 快速开始

### 1. 环境准备

确保已安装并启动以下服务：
- Go 1.26+
- Node.js 18+
- MySQL
- Redis

### 2. 创建数据库

```sql
CREATE DATABASE fullstack CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置

编辑 `backend/config/config.yml`，修改数据库连接信息：

```yaml
database:
  dsn: root:你的密码@tcp(localhost:3306)/fullstack?charset=utf8mb4&parseTime=True&loc=Local
```

Redis 默认连接 `localhost:6379`，无需密码（如需修改在 `backend/config/redis.go` 中调整）。

> **关于 CORS 和代理**：后端已配置 CORS 允许 `localhost:5173` 跨域，前端可直接调用后端接口。如需通过 Vite 代理转发（生产更推荐），可在 `frontend/vite.config.ts` 中添加 `server.proxy` 配置。

### 4. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

服务运行在 `http://localhost:3000`

### 5. 启动前端

```bash
cd frontend
npm install
npm run dev
```

浏览器自动打开 `http://localhost:5173`

### 6. 前后端联调

后端已配置 CORS 允许 `localhost:5173` 跨域请求（含 Cookie/Token）。开发环境建议配合 Vite proxy 使用，避免跨域问题。

## API 文档

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
| GET  | /api/articles          | 查询全部   | 否   |
| GET  | /api/articles/:id      | 查询单篇   | 是   |

### 点赞功能

| 方法 | 路径                     | 说明       | 认证 |
|------|-------------------------|------------|------|
| POST | /api/articles/:id/like   | 点赞文章   | 是   |
| GET  | /api/articles/:id/like   | 查询点赞数 | 否   |

## 架构说明

- **缓存策略**：文章列表采用旁路缓存（Cache-Aside），先查 Redis → 未命中则查 MySQL → 回填 Redis（TTL 10 分钟）；缓存命中时删除缓存中的旧数据，确保下次请求获取最新内容
- **点赞计数**：使用 Redis `INCR` 原子操作，并发安全，无需分布式锁
- **身份认证**：JWT（HS256），登录/注册后返回 token，后续请求放 `Authorization: Bearer <token>` 头
- **密码安全**：Bcrypt 哈希存储（cost=12），永不存明文
- **CORS**：开发环境允许 `localhost:5173` 跨域，生产环境交给 Nginx 处理
- **优雅关闭**：收到 SIGINT/SIGTERM 信号后等待 5 秒让现有请求处理完毕再退出

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

### 新增汇率（需认证）

```bash
curl -X POST http://localhost:3000/api/exchangerate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"fromCurrency":"USD","toCurrency":"CNY","rate":7.25}'
```

### 创建文章（需认证）

```bash
curl -X POST http://localhost:3000/api/articles \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"title":"我的第一篇文章","content":"文章正文","preview":"文章摘要"}'
```

### 点赞文章（需认证）

```bash
curl -X POST http://localhost:3000/api/articles/1/like \
  -H "Authorization: Bearer <token>"
```

### 查询文章点赞数

```bash
curl http://localhost:3000/api/articles/1/like
```

## 当前进度

### 后端（已完成）
- [x] 项目骨架搭建（分层架构）
- [x] MySQL + Redis 连接
- [x] 用户注册/登录（JWT 认证）
- [x] 汇率增查接口
- [x] 文章 CRUD + 旁路缓存
- [x] 文章点赞（Redis 原子计数）
- [x] JWT 鉴权中间件
- [x] 开发环境 CORS 配置
- [x] 优雅关闭

### 前端（基本完成）
- [x] Vite + Vue 3 + TypeScript 项目脚手架
- [x] Element Plus、Pinia、Axios 已安装
- [x] Axios 封装（baseURL + JWT 请求拦截器）
- [x] 顶部导航栏（Element Plus Menu + Vue Router）
- [x] 首页（欢迎页）
- [x] 登录页面（已接入 Pinia store，含错误处理）
- [x] 注册页面（已接入 Pinia store，含错误处理）
- [x] 文章列表页面（从后端获取、卡片展示、点赞交互）
- [x] 货币兑换页面（下拉选择币种、输入金额、计算兑换结果）
- [x] 路由配置（Vue Router，7 条路由，含根路径重定向）
- [x] Pinia 状态管理（auth store：token、登录/注册/退出、isLoggedIn 计算属性）
- [x] 退出登录功能（清除 token + 跳转首页）
- [x] 导航栏根据登录状态显示/隐藏菜单项（v-if + isLoggedIn）
- [x] 文章详情页（通过 query 参数传递文章数据）
- [x] 点赞交互（调用后端接口 + 前端实时更新点赞数）
- [ ] Vite 代理配置（目前通过后端 CORS 直连）
