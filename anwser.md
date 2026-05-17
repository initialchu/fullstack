# config.go 结构体解释

```go
type Config struct {
    App struct {
        Name string
        Port string
    }
}
```

- `Config` 是一个结构体，用于组织应用的配置信息
- 内部的 `App` 是一个**匿名字段**（没有字段名，类型直接是 `struct`），可以直接通过 `cfg.App.Name` 访问
- `Name` 和 `Port` 是两个字符串字段，分别存储应用名称和端口号

**使用方式：**
```go
var cfg Config
cfg.App.Name = "myapp"
cfg.App.Port = "8080"
```

**结构树：**
```
Config
 └── App
      ├── Name  (string)
      └── Port  (string)
```

> 注意：`package comfig` 应改为 `package config`，Go 惯例是目录名即包名。

---

# InitConfig 函数解释

```go
func InitConfig() {
    viper.SetConfigName("config")   // 配置文件名，不需要后缀
    viper.SetConfigType("yaml")     // 配置文件类型
    viper.AddConfigPath("./config") // 配置文件所在目录
    if err := viper.ReadInConfig(); err != nil {
        log.Fatalf("Error reading config file: %v", err)
    }
}
```

- `viper.SetConfigName("config")` — 告诉 Viper 要读的文件名，不加后缀，Viper 会自动匹配 `.yaml`、`.json` 等
- `viper.SetConfigType("yaml")` — 指定格式是 YAML（如果文件名带后缀可省略）
- `viper.AddConfigPath("./config")` — 指定配置文件目录，即去 `./config/` 下找
- `viper.ReadInConfig()` — 真正执行读取，组合上面三条：去 `./config/` 下读取 `config.yaml`
- `log.Fatalf` — Go 标准库 `log` 包的函数：
  - `Fatal` = 输出日志后调用 `os.Exit(1)` 直接终止程序
  - `f` = 支持格式化字符串（类似 `fmt.Printf`）
  - 如果读配置失败，打印错误并立刻终止，不带着错误配置继续运行

**Fatalf 对比：**
| 函数 | 行为 |
|------|------|
| `log.Printf` | 只打印，不退出 |
| `log.Panicf` | 打印后触发 panic（可被 recover 捕获） |
| `log.Fatalf` | 打印后直接 `os.Exit(1)`，无法捕获 |

---

# import 引用错误排查

**现象：** `main.go` 中 `import "fullstack/config"` 报错

**原因：** Go 的 import 路径 = `go.mod` 中 `module` 声明的值 + 子目录路径。你的 `go.mod` 里写的是：

```
module exchangeapp
```

但 `main.go` 的 import 写的是：

```go
import "fullstack/config"
```

模块名 `exchangeapp` ≠ `fullstack`，所以 Go 找不到这个包。

**修复：** 将 import 路径改为与 `go.mod` 中的模块名一致：

```go
import "exchangeapp/config"
```

> 规则：Go import 路径的根是模块名（`go.mod` 的 `module` 行），不是外层文件夹名。

---

# Viper 读取 YAML 配置报错

**错误信息：**
```
Error reading config file: While parsing config: yaml: line 5: mapping values are not allowed in this context
```

**原因：** YAML 要求冒号后必须有空格才能识别为键值对。配置文件中第 4 行缺少空格：

```yaml
# ❌ 错误写法
app:
  name:CurrencyExchange   # 冒号后没空格，被当成普通字符串
  port: 3000              # 解析器在这里才报错

# ✅ 正确写法
app:
  name: CurrencyExchange   # 冒号后必须有空格
  port: 3000
```

> YAML 语法规则：`key: value` 的冒号后面必须跟一个空格，否则 `key:value` 会被整体视为一个字符串而不是一个映射项。

---

# `r.Run()` 报错：missing port in address

**错误信息：**
```
[GIN-debug] [ERROR] listen tcp: address 3000: missing port in address
```

**原因：** `r.Run()` 要求的参数格式是 `:端口号`（如 `:3000`），但传的是不带冒号的 `"3000"`：

```go
// ❌ 错误
r.Run(config.AppConfig.App.Port)  // "3000"

// ✅ 正确
r.Run(":" + config.AppConfig.App.Port)  // ":3000"
```

---

# Gin 路由分组解释

```go
r := gin.Default()
auth := r.Group("/api/auth")
{
    auth.POST("/login", func(ctx *gin.Context) {
        ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
            "message": "Login successful",
        })
    })
    auth.POST("/register", func(ctx *gin.Context) {
        ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
            "message": "Register successful",
        })
    })
}
return r
```

- **`gin.Default()`** — 创建 Gin 引擎，内置 Logger（请求日志）和 Recovery（panic 恢复）中间件
- **`r.Group("/api/auth")`** — 创建路由组，组内所有路由自动带 `/api/auth` 前缀，方便统一管理或加中间件
- **`{ ... }`** — 大括号只是视觉分组，无语法作用，纯粹提高可读性
- **`auth.POST("/login", ...)`** — 注册 POST 路由，完整路径为 `/api/auth/login`
- **`ctx.AbortWithStatusJSON(http.StatusOK, gin.H{...})`** — 两个效果：
  - `Abort` — 阻止后续中间件/处理器执行
  - `WithStatusJSON` — 返回 JSON 并设置 HTTP 状态码（这里 200）
- **`return r`** — 把配置好的路由引擎返回给调用方

**生成的路由表：**
```
POST  /api/auth/login     →  {"message": "Login successful"}
POST  /api/auth/register  →  {"message": "Register successful"}
```

---

# GORM 数据库连接解释

```go
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

这行代码分两层理解：

**内层 `mysql.Open(dsn)`** — MySQL 专用的拨号器（Dialector）：
- `dsn`（Data Source Name）是数据库连接字符串，格式如 `user:password@tcp(host:port)/dbname?charset=utf8&parseTime=true`
- 返回一个 `gorm.Dialector`，表示"准备好连 MySQL 了"，但此时还没真正连接

**外层 `gorm.Open(拨号器, &gorm.Config{})`** — GORM 的核心入口：
- 第一个参数：数据库驱动，告诉 GORM 连的是什么库
- 第二个参数：`&gorm.Config{}`，GORM 自身配置（日志、命名策略等），这里用默认值
- 真正执行数据库连接
- 返回 `db`（后续所有 CRUD 都通过它）和 `err`（连接失败的信息）

**`:=`（短变量声明）** — 同时声明并赋值，自动推导类型，等价于：
```go
var db *gorm.DB
var err error
db, err = gorm.Open(...)
```

**整体流程：**
```
mysql.Open(dsn)  →  准备 MySQL 驱动（"拨号盘"）
gorm.Open(...)   →  真正连接数据库（"按下拨号键"）
if err != nil    →  失败则 log.Fatalf 退出程序
```

---

# Bcrypt 是什么

**Bcrypt** 是一个专门用于密码存储的单向哈希算法。

**核心特点：**

1. **自带盐值（Salt）** — 每次哈希自动生成不同的随机盐，相同密码两次哈希结果完全不同，防止彩虹表攻击
2. **可调节计算成本** — 有一个 `cost` 参数（通常 10~14），越大越慢。cost 加 1，计算时间翻倍，能抵抗硬件升级带来的暴力破解
3. **单向不可逆** — 无法从哈希值反推出明文，只能通过"用户输入 → 哈希 → 比对"来验证

**Go 中使用（golang.org/x/crypto/bcrypt）：**

```go
import "golang.org/x/crypto/bcrypt"

// 加密：把明文密码变成哈希值
hash, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
// hash 形如: $2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
//          ↑    ↑
//       算法  cost值

// 验证：比对用户输入和数据库中存的哈希
err := bcrypt.CompareHashAndPassword(savedHash, []byte("123456"))
if err != nil {
    // 密码错误
}
```

**为什么不用 MD5/SHA256？** 那些是通用哈希，计算太快，攻击者一秒能试几十亿次。Bcrypt 故意做得慢（一次约 0.1~0.3 秒），用户登录没体感，但对暴力破解是致命打击。
