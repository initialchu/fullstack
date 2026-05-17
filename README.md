# 汇率兑换应用 (Currency Exchange App)

Go + Gin + GORM + MySQL 全栈实战项目

## 技术栈

**后端：**
- Go 1.26
- [Gin](https://github.com/gin-gonic/gin) — HTTP Web 框架
- [GORM](https://gorm.io/) — ORM 框架
- [Viper](https://github.com/spf13/viper) — 配置管理
- [Bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt) — 密码哈希

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
│   ├── controllers/         # 控制器（处理请求）
│   │   └── auth_controller.go
│   ├── models/              # 数据模型
│   │   └── user.go          # 用户模型
│   ├── global/              # 全局变量
│   │   └── global.go
│   └── utils/               # 工具函数
│       └── utils.go
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

| 方法 | 路径              | 说明     |
|------|-------------------|----------|
| POST | /api/auth/register| 用户注册 |
| POST | /api/auth/login   | 用户登录 |
| GET  | /ping             | 健康检查 |
