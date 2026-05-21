# Table 'fullstack.articles' doesn't exist 错误

---

## 错误原因

错误发生位置：[article_controller.go:55](backend/controllers/article_controller.go#L55)，即 `GetArticles` 函数中的 `global.Db.Find(&articles)`。

**根本原因**：`articles` 表从未被创建过。

你的代码里，`AutoMigrate(&Article)` 只在 `CreateArticle` 函数中调用（[article_controller.go:30](backend/controllers/article_controller.go#L30)）：

```go
func CreateArticle(ctx *gin.Context) {
    // ...
    if err := global.Db.AutoMigrate(&article); err != nil {  // ← 只有这里调了
        // ...
    }
    // ...
}

而 `GetArticles` 直接查表，没有调用 AutoMigrate：

```go
    // 直接查表，没有 AutoMigrate
    if err := global.Db.Find(&articles).Error; err != nil {  // ← 这里报错
        // ...
    }
}
```

所以如果你**先访问 GET /api/articles（获取文章列表），再访问 POST /api/articles（创建文章）**，表根本还不存在，查询就会报 `Error 1146: Table doesn't exist`。

另外，`InitDB()` 函数（[db.go](backend/config/db.go)）里也没有调用 AutoMigrate——它只是连接了数据库，没有迁移任何表。

---

## 解决方案

**把 AutoMigrate 统一放到 `InitDB()` 里，在程序启动时一次性完成所有表的迁移。**

修改 [db.go](backend/config/db.go)：

```go
func InitDB() {
    dsn := AppConfig.Database.Dsn
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to connect database:%v", err)
    }

    // 启动时自动迁移所有模型（加在这里）
    if err := db.AutoMigrate(
        &models.Article{},
        &models.User{},
        &models.ExchangeRate{},
        // 以后新增的模型都加到这里
    ); err != nil {
        log.Fatalf("failed to migrate database: %v", err)
    }

    // ... 其余代码不变
}
```

同时，把 `CreateArticle` 里的 `AutoMigrate` 删掉（它在 init 中已经执行过了，不需要每次创建文章都跑一次）。

---

## 为什么不应该在 CreateArticle 里调用 AutoMigrate？

1. **AutoMigrate 是启动时做的事**，不是每次处理请求时做的事。每次请求都执行一遍是浪费性能。
2. **更严重的是时序问题**：如果用户先访问 GET 接口，表还没建，直接报错。
3. GORM 的 AutoMigrate 会在启动时自动创建/更新表结构，迁移完成后就不再需要了。

---

## 现在的临时解决方案

如果你不想改代码，可以先发一个 POST 请求到 `/api/articles` 创建一个文章（这会触发 `CreateArticle` 中的 AutoMigrate，把表建起来），之后 GET 就不会报错了。但这只是临时绕过，**推荐还是按上面的方案修改**。

---

## 登录时提示 record not found + 401 的原因与处理

你这段日志：

```
record not found
SELECT * FROM `users` WHERE username = 'aaa' AND `users`.`deleted_at` IS NULL
401 POST "/api/auth/login"
```

**核心结论**：数据库里没有用户名为 `aaa` 且未被软删除的用户，所以登录被拒绝并返回 401。GORM 抛出 `record not found`，这是“查不到数据”的正常结果，不是数据库异常。

### 常见原因

1. **用户根本没注册**：表里没有 `username = 'aaa'` 的记录。
2. **软删除导致查不到**：记录存在但 `deleted_at` 不为 NULL，会被 GORM 默认过滤掉。
3. **大小写问题**：如果你的数据库/字段排序规则是区分大小写（例如 `utf8mb4_bin`），`aaa` 和 `AAA` 会被当成不同用户。
4. **前端字段名不一致**：前端传的是 `userName` 或 `name`，后端按 `username` 去取，导致实际查询条件为空或不匹配。
5. **注册接口没成功写入**：注册时报错或事务回滚，但前端以为注册成功。

### 你可以这样排查

1. **直接查库确认**：
  ```sql
  SELECT id, username, deleted_at FROM users WHERE username = 'aaa';
  ```
2. **确认注册流程**：注册接口返回值是否成功、是否真的插入数据库。
3. **检查软删除**：如果 `deleted_at` 有值，登录查询会默认查不到；可先把该字段置空验证。
4. **检查前端登录请求字段**：确认传的是 `username`，不要是 `userName`、`name`。

### 如果你想让错误更清晰

在登录接口里区分两种情况：

- 找不到用户 → 返回“用户不存在”
- 密码不匹配 → 返回“密码错误”

这样前端和你自己排查都会更直观。

---

## 问题

当前 `onSubmit` 中：

```ts
result.value = rate.rate * form.amount
```

乘法结果可能是 `7.25 * 100 = 725`（整数），也可能是 `7.251 * 100 = 725.1000000000001`（浮点数精度问题），没有控制小数位数。

---

## 方法 1：`toFixed(2)`（最简单，推荐）

```ts
result.value = Number((rate.rate * form.amount).toFixed(2))
```

- `toFixed(2)` 返回**字符串**，比如 `"725.10"`，所以用 `Number()` 包一层转回数字
- 如果你的 `result` 是 `ref<number>`，那么 `Number()` 是必须的

如果你不介意 `result` 的类型变成 `string`（模板里 `{{ }}` 展示效果一样），也可以省掉 `Number()`：

```ts
const result = ref<string>('')  // 类型改成 string

// onSubmit 中：
result.value = (rate.rate * form.amount).toFixed(2)
```

---

## 方法 2：`Math.round`（纯数字方案）

```ts
result.value = Math.round(rate.rate * form.amount * 100) / 100
```

先放大 100 倍取整，再缩回去。比如 `725.100000001 * 100 → 72510.0000001 → Math.round → 72510 → /100 → 725.1`。

缺点：如果结果本身就是整数（如 `725`），不会补零，显示 `725` 而不是 `725.10`。

---

## 方法 3：在模板里用过滤器/计算属性（不改 onSubmit 逻辑）

```ts
// 模板里
{{ result.toFixed(2) }}
```

或者用计算属性：

```ts
const displayResult = computed(() => {
  return result.value % 1 === 0 ? result.value.toFixed(2) : result.value.toFixed(2)
  // 简单写就是：
  // return result.value.toFixed(2)
})
```

---

## 推荐

**直接用方法 1**，在 `onSubmit` 里一行搞定，最省事：

```ts
result.value = Number((rate.rate * form.amount).toFixed(2))
```

---

# 关于 `currencies.value = [...new Set(...)]` 这行代码的解释

## 这行代码在干什么？

```ts
currencies.value = [...new Set(res.data.map((rate: ExchangeInfo) => [rate.fromCurrency, rate.toCurrency]).flat())]
```

**一句话**：从所有汇率数据中，提取出所有出现过的货币并去重。

这行代码可以拆成四步来理解：

### 第 1 步：`res.data.map((rate) => [rate.fromCurrency, rate.toCurrency])`

`res.data` 是一个数组，例如：
```json
[
  { "fromCurrency": "USD", "toCurrency": "CNY", "rate": 7.25 },
  { "fromCurrency": "CNY", "toCurrency": "USD", "rate": 0.138 },
  { "fromCurrency": "EUR", "toCurrency": "CNY", "rate": 7.85 }
]
```

`map` 把每个汇率对象变成一个包含两个元素的数组 `[fromCurrency, toCurrency]`：

```json
[
  ["USD", "CNY"],
  ["CNY", "USD"],
  ["EUR", "CNY"]
]
```

### 第 2 步：`.flat()`

把嵌套数组"拍平"成一维数组：

```json
["USD", "CNY", "CNY", "USD", "EUR", "CNY"]
```

### 第 3 步：`new Set(...)`

`Set` 自动去重，只保留唯一值：

```json
Set { "USD", "CNY", "EUR" }
```

### 第 4 步：`[... ]`

展开运算符把 Set 转回普通数组：

```json
["USD", "CNY", "EUR"]
```

最终这个数组赋值给 `currencies.value`，用于下拉框中显示所有可选货币。

---

## 为什么要这样写？

因为汇率表里每个货币可能出现多次（USD 出现在第一条的 from、第二条的 to；CNY 出现在三条记录里），但下拉框只需要每个货币出现一次。所以核心需求是：**提取所有货币 + 去重**。

---

## 还有别的写法吗？

有，下面列出几种等价写法，各有优缺点：

### 写法 1：用 `flatMap` 替代 `map + flat`（推荐，更简洁）

```ts
currencies.value = [...new Set(res.data.flatMap((rate: ExchangeInfo) => [rate.fromCurrency, rate.toCurrency]))]
```

`flatMap` 就是 `map` + `flat` 的合体，一步到位，语义更清晰。

### 写法 2：传统 for 循环（最直观，适合新手）

```ts
const set = new Set<string>()
for (const rate of res.data) {
  set.add(rate.fromCurrency)
  set.add(rate.toCurrency)
}
currencies.value = [...set]
```

优点：每一步都看得懂，不需要理解 `flatMap`/`flat` 这些高阶函数。

### 写法 3：用 `reduce`（函数式风格）

```ts
currencies.value = [...new Set(
  res.data.reduce<string[]>((acc, rate) => [...acc, rate.fromCurrency, rate.toCurrency], [])
)]
```

缺点：每次迭代都创建新数组（`[...acc, ...]`），性能不如 `flatMap`，不推荐。

### 写法 4：`reduce` + `Set`（性能更好的 reduce 版本）

```ts
currencies.value = [...res.data.reduce<Set<string>>((set, rate) => {
  set.add(rate.fromCurrency)
  set.add(rate.toCurrency)
  return set
}, new Set())]
```

直接用 `Set` 作为 accumulator，避免中间数组。

### 写法 5：两次 `map` + `flat`（过度工程化，不推荐）

```ts
currencies.value = [...new Set([
  ...res.data.map(r => r.fromCurrency),
  ...res.data.map(r => r.toCurrency)
])]
```

先分别提取所有 `fromCurrency` 和 `toCurrency`，再合并去重。

---

## 总结建议

| 写法 | 可读性 | 性能 | 推荐度 |
|------|--------|------|--------|
| 原始写法 (`map + flat`) | ★★★ | ★★★★ | 还行 |
| `flatMap`（写法 1） | ★★★★ | ★★★★ | **最推荐** |
| for 循环（写法 2） | ★★★★★ | ★★★★★ | 新手友好 |
| `reduce` + Set（写法 4） | ★★ | ★★★★★ | 想炫技时用 |

个人建议把原来的 `map + flat` 改成 `flatMap`，改动最小，语义也更好：

```ts
currencies.value = [...new Set(res.data.flatMap((rate: ExchangeInfo) => [rate.fromCurrency, rate.toCurrency]))]
```

你的后端返回值是：
```json
{"fromCurrency":"USD","toCurrency":"CNY","rate":7.25}
```

但当前代码存在 **3 个问题**，导致拿不到数据。

---

## 问题一（致命）：字段名完全对不上

你前端定义的接口：

```ts
interface ExchangeInfo {
    from: string;    // ❌ 后端返回的是 fromCurrency
    to: string;      // ❌ 后端返回的是 toCurrency
    amount: number;  // ❌ 后端返回的是 rate
}
```

后端实际返回的字段是 `fromCurrency`、`toCurrency`、`rate`，和你的 interface 名字完全不一样。TypeScript 的 interface 只是编译时的类型标注，它**不会自动做字段映射**——JSON 解析后拿到的 key 仍然是 `fromCurrency`、`toCurrency`，但你代码里到处在访问 `rate.from`、`rate.to`，自然全是 `undefined`。

**修正：interface 字段必须和后端 JSON key 一致**

```ts
interface ExchangeInfo {
    fromCurrency: string
    toCurrency: string
    rate: number
}
```

---

## 问题二（致命）：把单个对象当数组处理

后端返回的是一个**单独对象**：
```json
{"fromCurrency":"USD","toCurrency":"CNY","rate":7.25}
```

但你的代码写成了：
```ts
const res = await axios.get<ExchangeInfo[]>('/exchangerate')
rates.value = res.data     // res.data 是一个对象，不是数组
```

你给 `rates` 声明的是 `ref<ExchangeInfo[]>([])`（数组），后续 `rates.value.find(...)` 也是数组方法。一个对象上调用 `.find()` 会直接报错。

**根据后端实际接口，有两种可能：**

- 如果 `/exchangerate` 确实只返回**一个汇率**，那就不该用数组，改成 `ref<ExchangeInfo | null>(null)`，也不要用 `.find()`
- 如果 `/exchangerate` 应该返回**所有汇率的列表**，比如：
  ```json
  [{"fromCurrency":"USD","toCurrency":"CNY","rate":7.25}, ...]
  ```
  那你提供的示例 JSON 就不对，需要确认后端实际返回的是什么结构

---

## 问题三（附带）：URL 有多余空格

```ts
'/exchangerate	'  // ← 有 tab 和尾部空格
```

虽然浏览器会自动 trim，但这属于不规范的写法，建议修正为 `'/exchangerate'`。

---

## 完整修正建议

假设后端接口返回的是**单个汇率对象**（如你的示例）：

```ts
interface ExchangeInfo {
    fromCurrency: string
    toCurrency: string
    rate: number
}

// 改成单个对象，不是数组
const rateInfo = ref<ExchangeInfo | null>(null)
const currencies = ref<string[]>([])
const result = ref<number>(0)

const fetchCurrencies = async () => {
    const res = await axios.get<ExchangeInfo>('/exchangerate')
    rateInfo.value = res.data
    currencies.value = [res.data.fromCurrency, res.data.toCurrency]
}

const onSubmit = () => {
    if (rateInfo.value) {
        result.value = rateInfo.value.rate * form.amount
    }
}
```

如果后端返回的是**全部汇率数组**，那示例 JSON 应该是 `[{...}, {...}]` 这样的结构，请先确认后端实际返回的格式，再来调整前端代码。

---

## 总结一句话

**前端 interface 的字段名必须和后端 JSON 的 key 完全一致，TypeScript 不会帮你做字段映射。** 另外要搞清楚后端返回的是单个对象还是数组，两者的处理方式完全不同。
 
