# 汇率兑换应用 (Currency Exchange App)

Go + Gin + GORM + MySQL 全栈实战项目

## 技术栈

**后端：**
- Go 1.26
- [Gin](https://github.com/gin-gonic/gin) — HTTP Web 框架
- [GORM](https://gorm.io/) — ORM 框架
- [Viper](https://github.com/spf13/viper) — 配置管理
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — 密码哈希
- [JWT](https://github.com/golang-jwt/jwt) — 身份认证令牌

**数据库：**
- MySQL

## 项目结构

```
fullstack/
├── backend/
│   ├── main.go              # 入口文件
│   ├── config/              # 配置相关
│   │   ├── config.go        # 配置结构体 & Viper 加载
│   │   ├── config.yml       # YAML 配置文件
│   │   └── db.go            # 数据库连接初始化
│   ├── router/              # 路由
│   │   └── router.go        # 路由注册
│   ├── middlewares/         # 中间件
│   │   └── auth_middleware.go # JWT 认证中间件
│   ├── controllers/         # 控制器（处理请求）
│   │   ├── auth_controller.go        # 用户认证（注册/登录）
│   │   └── auth_rate_controller.go   # 汇率数据接口
│   ├── models/              # 数据模型
│   │   ├── user.go          # 用户模型
│   │   └── exchange_rate.go # 汇率模型
│   ├── global/              # 全局变量
│   │   └── global.go        # 数据库连接实例
│   └── utils/               # 工具函数
│       └── utils.go         # 密码哈希、JWT 生成与验证
```

## 快速开始

### 1. 创建数据库

```sql
CREATE DATABASE fullstack CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 修改配置

编辑 `backend/config/config.yml`，根据你的环境修改数据库连接信息：

```yaml
database:
  dsn: root:你的密码@tcp(localhost:3306)/fullstack?charset=utf8mb4&parseTime=True&loc=Local
```

### 3. 安装依赖

```bash
cd backend
go mod tidy
```

### 4. 运行

```bash
go run main.go
```

服务默认运行在 `http://localhost:3000`

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

### 查询所有汇率

```bash
curl -X GET http://localhost:3000/api/exchangerate
```
