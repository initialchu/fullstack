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
