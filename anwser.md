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

---

# `[]byte(pwd)` 类型转换

`[]byte(pwd)` 是 Go 的**类型转换**，把 `string` 转成 `[]byte`（字节切片）。

Go 中 `string` 和 `[]byte` 底层结构相似，可以互转：

```go
s := "hello"
b := []byte(s)     // string → []byte → [104 101 108 108 111]

b2 := []byte{104, 101, 108, 108, 111}
s2 := string(b2)   // []byte → string → "hello"
```

**Bcrypt 中为什么需要？** `bcrypt.GenerateFromPassword` 和 `bcrypt.CompareHashAndPassword` 的参数类型是 `[]byte`，不是 `string`，所以必须转换。这是该库的设计选择，因为密码这种敏感数据用 `[]byte` 可以在用完后手动清零，而 `string` 是不可变的，会一直留在内存中直到 GC 回收。

---

# GORM AutoMigrate 解释

```go
if err := global.Db.AutoMigrate(&user); err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
        "error": "Failed to migrate database",
    })
    return
}
```

**`AutoMigrate`** 是 GORM 的自动迁移功能，根据 Go 结构体同步数据库表结构。

**具体行为：**
1. 检查 `users` 表是否存在，不存在就自动创建（字段根据 `User` 结构体生成）
2. 表已存在时，检测结构体是否有**新增字段**，有则在表中添加对应列
3. **不会删除**已有列、**不会修改**已有列类型——只新增，不破坏

**为什么放在这里？** 在注册逻辑中调用，确保第一次注册时表一定存在。开发阶段省去手动建表 SQL。

**错误处理：** 如果建表/同步失败（如数据库断连），返回 500 并终止后续写入操作。

> 注意：生产环境建议用正式的迁移工具（如 golang-migrate），而不是在业务代码里跑 AutoMigrate。

---

# `ctx.ShouldBindJSON` 解释

```go
var input struct {
    Username string `json:"username"`
    Password string `json:"password"`
}
if err := ctx.ShouldBindJSON(&input); err != nil {
    ctx.JSON(http.StatusBadRequest, gin.H{
        "error": err.Error(),
    })
    return
}
```

**`ShouldBindJSON(&input)`** 做了三件事：

1. **读取请求体** — 从 HTTP 请求中读取 JSON 原始数据
2. **反序列化** — 把 JSON 解析并映射到结构体字段，映射依据是 `json` tag：
   ```json
   {"username": "alice", "password": "123456"}
   ```
   执行后 `input = {Username: "alice", Password: "123456"}`
3. **校验** — 如果字段类型不匹配（比如 `username` 传了数字），返回 error

**为什么要传 `&input`（指针）？** 函数需要**修改** `input` 的值（把解析结果填进去），值传递无法修改原变量。

**失败时：** 返回 400（Bad Request），表示客户端发的数据有问题。

> Gin 命名惯例：`ShouldBindJSON` 失败后返回 error 让你自己处理；`BindJSON` 失败会自动返回 400。`Should` 开头表示"我来判断要不要报错"。

---

# GORM 链式查询：根据用户名查用户

```go
if err := global.Db.Where("username = ?", input.Username).First(&user).Error; err != nil {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "无效的用户名或密码",
    })
    return
}
```

**逐步拆解：**

**`Where("username = ?", input.Username)`** — 添加查询条件：
- `?` 是占位符，GORM 用 `input.Username` 的值填充
- 生成 SQL：`SELECT * FROM users WHERE username = 'alice'`
- 用 `?` 而非字符串拼接，**防止 SQL 注入**

**`.First(&user)`** — 执行查询，取第一条匹配记录：
- 把查询结果**写入** `&user`（传指针改写原变量）
- 找到 → `user` 包含该用户所有字段（ID、Password 等）
- 没找到 → 返回 `gorm.ErrRecordNotFound`

**`.Error`** — 获取执行过程中的错误（找不到用户或数据库异常都为非 nil）

**为什么错误提示不区分"用户不存在"和"密码错误"？**

出于**安全考虑**。如果分开返回，攻击者可以批量尝试用户名枚举出已注册账号。模糊提示让其无法判断到底是用户名不存在还是密码不对，增加攻击难度。

---

# 为什么登录和注册都生成 JWT？

```go
// 注册接口最后
token, err := utils.GenerateJWT(user.Username)
// ...
ctx.JSON(200, gin.H{"token": token})

// 登录接口最后
token, err := utils.GenerateJWT(user.Username)
// ...
ctx.JSON(200, gin.H{"token": token})
```

**注册时返回 token = 注册即自动登录**

注册成功后如果还要跳回登录页重新输入账号密码，体验很差。直接下发 token，客户端拿到后等同于已登录，无缝进入应用。

**登录时返回 token = 正常的身份验证凭证**

验证用户名密码通过后，发 token 作为后续请求的身份凭证（后续请求带 token，服务端就知道是谁在操作）。

**本质上**：注册 = 创建用户 + 登录，两步合二为一，所以接口行为一致，都返回 token。

---

# GORM `Find` 查询全表

```go
if err := global.Db.Find(&exchangeRates).Error; err != nil {
    ctx.JSON(http.StatusInternalServerError, gin.H{
        "error": err.Error(),
    })
    return
}
```

**`Find(&exchangeRates)`** — 查询表中**所有**记录，写入切片。

等同于 SQL：
```sql
SELECT * FROM exchange_rates;
```

**逐层解释：**
- `Find()` — 不带 `Where` 条件，查全表
- `&exchangeRates` — 传切片指针，GORM 把查询结果**填充**进去
- `.Error` — 获取执行错误（连接中断、表不存在等）

**和 `First()` 对比：**

| 方法 | 效果 | 没找到时 |
|------|------|----------|
| `First(&user)` | 查**一条** | 返回 `ErrRecordNotFound` |
| `Find(&slice)` | 查**全部** | 返回空切片，不报错 |

---

# Gin 中间件写法：AuthMiddleware

```go
func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token := ctx.GetHeader("Authorization")
        if token == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{
                "error": "Authorization header is missing",
            })
            ctx.Abort()
            return
        }
    }
}
```

**`gin.HandlerFunc`** — 本质是 `func(*gin.Context)` 的别名，Gin 中间件和处理器都遵循这个签名。

**为什么要包一层闭包（`return func`）？**

- 闭包 = 外层函数 + 内层返回 `func(*gin.Context)`
- 好处：后续想加参数时不需要改调用方的代码，比如 `AuthMiddleware("admin")` 按角色鉴权
- 外层传参，内层使用，对外接口保持不变

**`ctx *gin.Context` 参数怎么来的？**

Gin 框架收到请求时自动创建 `*gin.Context`，封装了该次请求的全部信息（请求头、请求体、参数、响应方法），然后传入中间件链。开发者只负责从 `ctx` 取数据、写响应。

**`ctx.GetHeader("Authorization")`** — 从 HTTP 请求头取 `Authorization` 字段，客户端 JWT 通常放在这里。

**`ctx.Abort()`** — 核心！只返回 JSON 不会中断流程，必须调用 `Abort()` 才能阻止后续中间件和处理器继续执行。`return` 退出当前函数，`Abort()` 阻止 Gin 继续往下调。

**整体流程：**
```
请求进来 → Gin 创建 ctx → 进入 AuthMiddleware
                              → 取 Authorization 头
                              → 为空？返回 401 + Abort()，请求终止
                              → 不为空？验证 token，通过则 c.Next()
```

---

# Bearer Token 格式修复

在 `utils.go` 中生成 JWT 时：

```go
// ❌ 错误：Bearer 和 token 之间没有空格
return "Bearer" + Token, err

// ✅ 正确：必须有空格
return "Bearer " + Token, err
```

**原因：** HTTP `Authorization` 头的标准格式是 `Bearer <token>`，中间必须有空格。没有空格的话，中间件按空格分割取第二部分时，拿到的不是纯 token，导致验证失败。

---

# `ParseJWT` 解析与验证令牌

```go
func ParseJWT(tokenString string) (string, error) {
    if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
        tokenString = tokenString[7:]
    }
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("unexpected signing method")
        }
        return []byte("secret"), nil
    })
    if err != nil {
        return "", err
    }
    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        username, ok := claims["username"].(string)
        if !ok {
            return "", errors.New("username claim is not a string")
        }
        return username, nil
    }
    return "", err
}
```

**第 1 步：去掉 `Bearer ` 前缀**

```go
if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
    tokenString = tokenString[7:]
}
```

- `len(tokenString) > 7` — 安全判断，防止空串或短字符串越界
- `tokenString[:7]` — 取前 7 个字符比对标不标准 `"Bearer "`
- `tokenString[7:]` — 截掉前 7 个字符，剩下纯 token

**第 2 步：解析并验证签名**

```go
token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
    if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
        return nil, errors.New("unexpected signing method")
    }
    return []byte("secret"), nil
})
```

- `jwt.Parse(待解析token, 密钥回调)` — 解析和验证二合一
- **回调函数**：库先读 header 中的 `alg` 字段，然后调用你的回调，你返回对应的密钥，库完成签名比对
- **`token.Method.(*jwt.SigningMethodHMAC)`** — 类型断言，确认使用 HMAC 算法。防止攻击者伪造 `"alg":"none"` 的 token 绕过验证
- **`return []byte("secret")`** — 返回与生成时相同的密钥

**第 3 步：提取载荷中的用户名**

```go
if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
    username, ok := claims["username"].(string)
    // ...
    return username, nil
}
```

- `token.Claims.(jwt.MapClaims)` — 类型断言，把 Claims 转成 map 方便取值
- `.Valid` — 库已自动校验 `exp`（过期），过期则为 `false`
- `claims["username"].(string)` — 从 map 取值再做类型断言，确保值确实是 string

---

# 为什么用 Redis？跟 MySQL 是什么关系？

**MySQL 和 Redis 不是替代关系，是协作关系。**

**MySQL 的问题：** 每次查询走磁盘 I/O，即便用 B+ 树索引也要几毫秒。高并发下几百个请求同时查汇率，MySQL 撑不住。

**Redis 的优势：**
1. **纯内存操作** — 数据在内存中，单次读取 0.1 毫秒，比 MySQL 快几十倍
2. **缓存高频数据** — 汇率短时间内不变，不需要每次都查 MySQL。启动时从 MySQL 加载进 Redis，后续直接读 Redis
3. **减轻 MySQL 压力** — 99% 的读请求被 Redis 拦截，MySQL 只处理写操作

**协作流程：**
```
客户端请求汇率
    → 先查 Redis（有就直接返回，0.1ms）
    → Redis 没命中 → 查 MySQL → 写入 Redis → 返回
```

**各自职责：**

| | MySQL | Redis |
|------|-------|-------|
| 职责 | 持久化存储（数据不丢） | 加速读取（扛并发） |
| 速度 | 毫秒级 | 0.1 毫秒级 |
| 存储 | 磁盘 | 内存 |
| 数据 | 永久 | 可设置过期 |

---

# Redis `INCR` 点赞计数

```go
likeKey := "article:" + articleID + ":likes"
if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
    // ...
}
```

**`INCR` 的行为：**
- key **不存在** → Redis 先自动设值为 `0`，再 `+1`，返回 `1`
- key **已存在** → 直接在现有值上 `+1`

**调用这个函数前，Redis 里可能根本没有** `article:1:likes` **这个键；调用后一定存在，值至少为 1。**

不需要先检查 key 是否存在、也不需要手动初始化为 0——Redis 的 `INCR` 原子操作天然帮你处理了。既省了代码，又保证了并发安全（两个请求同时点赞不会导致计数错误）。

---

# GORM 哪些方法可以接 `.Error`？

GORM 链式方法分两类：

**终结方法（可以接 `.Error`）** — 真正执行 SQL，`Error` 字段会被填充：

| 方法 | 作用 |
|------|------|
| `Find(&slice)` | 查询多条 |
| `First(&struct)` | 查询一条 |
| `Create(&struct)` | 插入 |
| `Save(&struct)` | 保存 |
| `Delete(&struct)` | 删除 |
| `Update("col", val)` | 更新 |

```go
// ✅ 有意义的 .Error
global.Db.Find(&articles).Error
global.Db.Create(&article).Error
```

**中间方法（接 `.Error` 无意义）** — 只构建查询条件，不访问数据库，`Error` 永远是 `nil`：

| 方法 | 作用 |
|------|------|
| `Where(...)` | 加条件 |
| `Order(...)` | 排序 |
| `Limit(n)` | 限制行数 |
| `Select(...)` | 选择字段 |

```go
// ❌ 无意义，永远是 nil
global.Db.Where("id=?", 1).Error
```

> 判断规则：只有真正访问数据库的方法才会填充 `Error`，纯拼 SQL 片段的方法永远不出错。

---

# go-redis 为什么需要 `.Result()`？

```go
cachedData, err := global.RedisDB.Get(cacheKey).Result()
```

**`Get()` 返回的是 `*redis.StringCmd`，不是直接的结果。**

**为什么这样设计？**

1. **统一返回模式** — 所有命令都返回 `(结果, error)` 对，写法一致
2. **多结果选择** — 同一个命令包装器提供多种取值方式：
   - `.Result()` → `(string, error)`，通用
   - `.Int()` → `(int, error)`，用于 `Incr` 等
   - `.Float64()` → `(float64, error)`
3. **延迟执行** — 可以先发命令，稍后再取结果（Pipeline 场景）

**常见对应关系：**

| Redis 方法 | 返回的包装类型 | 取值方法 |
|------------|--------------|----------|
| `Get(key)` | `*StringCmd` | `.Result()` |
| `Incr(key)` | `*IntCmd` | `.Result()` |
| `HGet(key, field)` | `*StringCmd` | `.Result()` |
| `SAdd(key, vals...)` | `*IntCmd` | `.Result()` |

> 判断规则：返回类型以 `*xxxCmd` 结尾就必须接 `.Result()`（或 `.Err()` 只要错误）。

---

# 旁路缓存（Cache-Aside）模式

```go
func GetArticles(ctx *gin.Context) {
    cachedData, err := global.RedisDB.Get(cacheKey).Result()
    if err == redis.Nil {
        // 缓存未命中 → 查 MySQL
        var articles []models.Article
        global.Db.Find(&articles)
        // 回填缓存
        articleJSON, _ := json.Marshal(articles)
        global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute)
        ctx.JSON(200, articles)
    } else if err != nil {
        // Redis 挂了
        ctx.JSON(500, gin.H{"error": err.Error()})
    } else {
        // 缓存命中 → 反序列化返回
        var articles []models.Article
        json.Unmarshal([]byte(cachedData), &articles)
        ctx.JSON(200, articles)
    }
}
```

**整体流程：**
```
请求 → 先查 Redis
         ├─ 命中 → json.Unmarshal → 直接返回（0.1ms）
         └─ 未命中 → 查 MySQL → json.Marshal → 写 Redis → 返回
```

**各函数职责：**

| 函数 | 作用 | 为什么需要 |
|------|------|-----------|
| `RedisDB.Get()` | 查缓存 | 内存操作 0.1ms，挡住大部分请求 |
| `json.Marshal()` | Go struct → JSON | Redis 只能存字节/字符串 |
| `RedisDB.Set(key, val, 10min)` | 回填缓存 | TTL 必须设，防内存膨胀和数据永远陈旧 |
| `json.Unmarshal()` | JSON → Go struct | 缓存中的 JSON 还原成 Go 结构体 |

**旁路缓存的优点：**
1. **读性能大幅提升** — Redis 0.1ms vs MySQL 3~10ms
2. **减轻数据库压力** — 99% 读请求被 Redis 拦截
3. **实现简单** — 三行判断命中/未命中，不需要额外中间件
4. **容错性好** — Redis 挂了只是多一次延迟，不会直接崩溃
5. **数据自动刷新** — TTL 过期后自动从 MySQL 拉最新数据

---

# 缓存删除后返回，数据会丢失吗？

在 `CreateArticle` 中：

```go
global.Db.Create(&article)        // 1. 先写 MySQL
global.RedisDB.Del(cacheKey)      // 2. 再删缓存
ctx.JSON(200, article)            // 3. 返回新文章
```

**不会丢失，因为三个步骤各走各的路：**

```
1. Create(&article)  →  MySQL 已持久化（数据安全）
2. Del(cacheKey)     →  删的是缓存中的旧列表（不包含新文章）
3. JSON(200, article) →  返回的是内存变量 article，不经过缓存
```

- 第 2 步是**缓存失效**——旧的 `articles` 列表里没有这条新文章，删掉
- 下次有人 `GET /api/articles` → 缓存未命中 → 从 MySQL 重新查全量 → 回填 Redis，新列表就包含新文章了
- 第 3 步返回的 `article` 是第 1 步写入时就存在的内存变量，和缓存无关

> 类比：更新纸质原档（MySQL）后，把前台旧复印件扔掉（Redis），明天有人来拿时重新复印最新版。你手里正在看的那页不受影响。

---

# 开发/生产环境跨域处理方案

**生产环境：** Nginx 在反向代理层处理 CORS，Go 代码不加 CORS 中间件。

**开发环境（三个方案）：**

**方案 1：Vite 代理（推荐）**

在 Vue3 前端 `vite.config.js` 中配置代理：

```js
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true
      }
    }
  }
})
```

原理：前端发请求到同源（`localhost:5173/api/...`），Vite 自动转发到 Go 后端，浏览器认为是同源，不存在跨域。不改 Go 代码，不引入额外依赖，是前端开发标配。

**方案 2：Go 代码条件加载 CORS**

```go
if os.Getenv("ENV") == "dev" {
    r.Use(cors.Default())
}
```

**方案 3：浏览器插件**

装 Chrome 的 "Allow CORS" 插件，最简单但不推荐团队使用——每个人都要装，且可能忘记关。

> 推荐方案 1，开发环境用 Vite proxy，生产环境用 Nginx，Go 代码始终保持干净。

---

# Go 包级别变量的作用域

```go
var AppConfig *Config
```

**作用域不是目录，是按包（package）。**

Go 强制同目录下所有 `.go` 文件声明同一个 `package`，所以效果上等同于同目录都能用，但本质是按包：

```
backend/config/
├── config.go   →  package config  →  var AppConfig *Config（定义）
├── db.go       →  package config  →  直接用 AppConfig ✓
└── redis.go    →  package config  →  直接用 AppConfig ✓
```

**同目录 = 同包 = 自然共享所有变量，不用 import，不用传参。**

**大写 vs 小写：**

```go
var AppConfig *Config  // 大写 → 跨包可见，config.AppConfig
var appConfig *Config  // 小写 → 仅包内可见，config 包外访问不了
```

> 核心规则：小写 = 包内私有，大写 = 跨包公开。同目录下所有文件属于同一个包，自然共享一切。

---

# `os.Getenv` 与 `.env` 文件

**`os.Getenv` 不能读取 `.env` 文件。**

它只读取操作系统级别的环境变量（shell 里 `export` 设置或启动时传入的）。要读 `.env` 文件，Go 需要引入 `github.com/joho/godotenv`：

```go
import "github.com/joho/godotenv"

func main() {
    godotenv.Load()           // 加载 .env 文件
    env := os.Getenv("ENV")   // 现在才能拿到 .env 里的值
}
```

**一个项目可以有多个 `.env` 文件：**

前后端是各自独立进程，各有各的目录和 `.env`：

```
fullstack/
├── backend/
│   └── .env    ← Go 后端：数据库密码、JWT 密钥等
├── frontend/
│   └── .env    ← Vue 前端：VITE_API_BASE_URL 等
```

这是标准做法——两个应用配置需求不同，不共享一个 `.env`。

**Vite vs Go：** 前端 Vite 自动读取 `.env`（变量需 `VITE_` 前缀），Go 的 `os.Getenv` 不会自动读，必须用 godotenv 手动加载。

---

# 用 `.env` 区分开发/生产环境（router.go 实践）

**1. `main.go` 最前面加载 `.env`：**

```go
func main() {
    godotenv.Load()            // 注入 .env 到进程环境变量
    config.InitConfig()
    r := router.SetupRouter()  // 此时 os.Getenv("ENV") 能读到了
    r.Run(...)
}
```

> 必须在 `main()` 最开头调用——Go 的 `init()` 函数在 `main()` 之前执行，如果 init 里用 `os.Getenv` 会拿不到 `.env` 里的值。

**2. `backend/.env`：**

```
ENV=dev
```

**3. `router.go` 中判断环境：**

```go
func SetupRouter() *gin.Engine {
    r := gin.Default()
    if os.Getenv("ENV") == "dev" {
        r.Use(cors.New(cors.Config{}))  // 开发环境加 CORS
    }
    // ... 路由注册
    return r
}
```

**4. `.gitignore` 必须加 `.env`**（避免密码泄露到 GitHub）

**启动时区分环境：**

```bash
# 开发（读 .env）
go run main.go

# 生产（不依赖 .env，设真实环境变量）
ENV=prod go run main.go
```

---

# CORS `AllowOriginFunc` 参数

`AllowOriginFunc` 是一个动态判断函数，对每个请求来源做自定义逻辑：

```go
r.Use(cors.New(cors.Config{
    AllowOriginFunc: func(origin string) bool {
        // origin = 浏览器发来的完整域名，如 "https://example.com"
        // 返回 true = 允许，false = 拒绝
        return strings.HasSuffix(origin, ".trusted.com")
    },
}))
```

**与静态 `AllowOrigins` 对比：**

| 参数 | 用法 | 适用场景 |
|------|------|----------|
| `AllowOrigins` | 写死 `["https://a.com"]` | 域名固定、数量少 |
| `AllowOriginFunc` | 回调动态判断 | 正则匹配、多子域名、查库白名单 |

**典型场景：**

```go
// 多租户：允许所有 .myapp.com 子域名
AllowOriginFunc: func(origin string) bool {
    return strings.HasSuffix(origin, ".myapp.com")
}

// 开发通配，生产严格
AllowOriginFunc: func(origin string) bool {
    if os.Getenv("ENV") == "dev" {
        return true
    }
    return origin == "https://production.com"
}
```

> `AllowOriginFunc` 和 `AllowOrigins` 不能同时用，二选一。

---

# go-redis Get 函数参数差异

## 问题

官方 Redis 文档示例中 `Get` 需要两个参数：`rdb.Get(ctx, "bike:1").Result()`，但我项目中只传了一个参数：`global.RedisDB.Get(cacheKey).Result()`，为什么？

## 原因

**版本不同。** 你的项目用的是 v6 老版本，官方文档展示的是 v7+ 新版本。

| | v6（项目在用） | v7+（官方文档） |
|---|---|---|
| 包路径 | `github.com/go-redis/redis` | `github.com/redis/go-redis/v9` |
| `Get` 签名 | `Get(key string) *StringCmd` | `Get(ctx context.Context, key string) *StringCmd` |
| 参数个数 | **1个** | **2个** |
| 其他方法 | `Set(key, val, ttl)` | `Set(ctx, key, val, ttl)` |

你 go.mod 中的版本：
```
github.com/go-redis/redis v6.15.9+incompatible
```

## 为什么新版要加 context？

Go 1.7 之后，`context` 成为网络库的标准实践，它可以：

- **超时控制**：`ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)`，超过 2 秒自动取消 Redis 请求
- **请求级联取消**：用户关闭浏览器 → Gin 的 `ctx.Request.Context()` 被取消 → 传播到 Redis 操作 → 不再浪费资源
- **链路追踪**：在 context 中携带 trace ID，方便排查问题

## 类比理解

就像寄快递：
- **v6**：直接扔给快递员，不管他送没送到（没有超时和取消机制）
- **v9**：给快递员一个倒计时闹钟（超时）+ 一个对讲机（取消信号），超时或你反悔了可以立刻通知他

## 需要升级吗？

暂时**不需要**。v6 功能完全够用。如果以后要升级到 v9，需要注意：
1. 包路径变了，import 要改
2. 所有 Redis 方法都要加 `ctx context.Context` 作为第一个参数
3. 初始化方式也有变化（`redis.NewClient` → `redis.NewClient(&redis.Options{...})`）

---

# 项目开发思路总结

## 一、整体架构

采用 Go 后端标准分层架构，关注点分离：

```
请求 → Router（路由 + CORS）
        → Middleware（JWT 鉴权）
          → Controller（处理业务）
            → Model（数据结构）
              → MySQL（持久化）/ Redis（缓存 + 计数）
```

**六层职责：**

| 层 | 目录 | 职责 |
|------|------|------|
| 入口 | `main.go` | 启动、优雅退出 |
| 路由 | `router/` | URL 映射、CORS、中间件挂载 |
| 中间件 | `middlewares/` | JWT 鉴权、跨切面逻辑 |
| 控制器 | `controllers/` | 请求处理、业务编排 |
| 数据模型 | `models/` | 结构体定义、GORM 映射 |
| 工具 | `utils/` | 密码哈希、JWT 生成/验证 |

---

## 二、开发思路

### 1. 配置驱动启动

`main.go` 只做三件事：加载配置 → 初始化路由 → 启动服务器。不写业务逻辑，保证入口文件简洁。

```go
config.InitConfig()       // 一切配置从这里开始
r := router.SetupRouter() // 路由集中管理
srv.ListenAndServe()      // 启动
```

### 2. 全局变量集中管理

通过 `global` 包统一持有 MySQL 和 Redis 连接，避免在控制器之间传递依赖：

```go
global.Db      // 任何地方直接用
global.RedisDB // 任何地方直接用
```

### 3. 配置与初始化合并

`InitConfig()` 不仅读 YAML，还顺带调用 `InitDB()` 和 `InitRedis()`，一个调用完成全部初始化，减少 `main.go` 的复杂度。

### 4. 路由分组区分权限

```go
api.GET("/exchangerate", ...)           // 公开接口
api.Use(middlewares.AuthMiddleware())     // 从此往下全部需要 JWT
{
    api.POST("/exchangerate", ...)      // 需认证
    api.POST("/articles", ...)          // 需认证
}
```

通过 `Group` + `Use` 把公开和认证接口隔开，结构清晰，不会漏加认证。

### 5. Write-Through 缓存失效

创建文章时：写 MySQL → 立即删除 Redis 缓存。下次读取自动回填最新数据。保证**缓存和数据库最终一致**。

### 6. 注册即登录

注册接口完成后直接返回 JWT Token，减少用户操作步骤，前端拿 Token 后直接进入应用。

---

## 三、项目亮点

### 1. 优雅关闭

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
// 5 秒超时优雅关闭
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
srv.Shutdown(ctx)
```

收到 `Ctrl+C` 或 `kill` 信号后不是暴力退出，而是等待现有请求处理完（最多等 5 秒），防止请求中断导致数据不一致。

### 2. 双中间件链

```
CORS 中间件（全局）
    → JWT 鉴权中间件（认证路由组）
        → 业务处理器
```

全局 CORS + 按路由组 JWT，中间件挂载层次分明。

### 3. 旁路缓存 + Cache Invalidation

- 读：Redis → 未命中 → MySQL → 回填 Redis（TTL 10 分钟）
- 写：MySQL → 删 Redis 缓存 → 下次读自动回填

保证了数据最终一致性，同时大幅提升读性能。

### 4. Redis 原子点赞

使用 `INCR` 命令，天然支持并发安全——两个用户同时点赞不会导致计数错误。不需要分布式锁。

### 5. Bcrypt 密码存储

Cost 设为 12（2¹² = 4096 轮哈希），兼顾安全性与登录体验。密码永不落明文。

### 6. JWT 无状态认证

HS256 对称签名，24 小时过期。中间件统一解析 → 注入 `ctx.Set("username", ...)` → 后续处理器通过 `ctx.Get("username")` 获取当前用户，跨中间件传递上下文。

### 7. CORS 配置明确

`AllowCredentials: true` 支持前端携带 Token，`AllowOrigins` 精确限制 `localhost:5173`，不开放通配符 `*`，避免安全隐患。

---

## 四、技术栈一览

| 组件 | 库 | 用途 |
|------|------|------|
| Web 框架 | `gin-gonic/gin` | HTTP 路由、中间件、JSON 响应 |
| ORM | `gorm.io/gorm` | MySQL 数据访问 |
| 配置 | `spf13/viper` | YAML 配置加载 |
| 缓存 | `go-redis/redis` | Redis 缓存 + 点赞计数 |
| 密码 | `golang.org/x/crypto` | Bcrypt 哈希 |
| 认证 | `golang-jwt/jwt` | JWT 令牌生成/验证 |
| 跨域 | `gin-contrib/cors` | 开发环境 CORS |
| MySQL 驱动 | `gorm.io/driver/mysql` | GORM MySQL 驱动 |