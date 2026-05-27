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

---

# 登录失败如何不跳转页面（try/catch 错误处理）

## 问题

当前 `Login.vue` 的 `login` 函数中，`await axios.post(...)` 没有包在 `try/catch` 里。当后端返回 401（用户名或密码错误）时，axios 会抛出异常，控制台报错，用户看不到任何提示，体验很差。

## 当前代码的问题

```ts
const login = async () => {
    // ...
    const res = await axios.post('/auth/login', { username, password })
    // 如果上面抛异常，下面代码都不会执行
    token.value = res.data.token
    localStorage.setItem('token', token.value || '')
    router.push({ name: 'home' })
}
```

axios 请求失败（HTTP 状态码 ≥ 400）时，直接抛出异常，函数终止。但因为没有 `catch`，异常是**未捕获**的，浏览器只在控制台打印错误，用户看不到任何反馈。

## 修正方案

```ts
const login = async () => {
    const username = form.value.username
    const password = form.value.password
    if (!username || !password) {
        alert('请输入用户名和密码')
        return
    }

    try {
        const res = await axios.post('/auth/login', { username, password })
        // 只有请求成功（HTTP 2xx）才会执行到这里
        token.value = res.data.token
        localStorage.setItem('token', token.value || '')
        router.push({ name: 'home' })
    } catch (err: any) {
        // 登录失败，不跳转，显示错误信息
        const msg = err.response?.data?.error || '登录失败，请重试'
        alert(msg)
    }
}
```

## 关键点解释

| 代码 | 作用 |
|------|------|
| `try { ... }` | 包裹可能抛出异常的代码 |
| `catch (err: any)` | 捕获异常，拿到错误对象 |
| `err.response?.data?.error` | 从 axios 错误中取出后端返回的 error 字段（如"无效的用户名或密码"） |
| `?.` 可选链 | 防止 `response`/`data` 为 `undefined` 时再报错，安全回退到默认提示 |
| `alert(msg)` | 弹出错误提示，页面留在当前页，不跳转 |

## 执行流程对比

**没有 try/catch：**
```
请求失败 → 抛异常 → 函数中断 → 控制台报错 → 用户看到页面卡住
```

**有 try/catch：**
```
请求失败 → 抛异常 → catch 捕获 → alert 弹窗 → 用户看到错误信息，留在登录页
请求成功 → 不抛异常 → 执行 try 内代码 → 存 token → 跳转首页
```

## 同样的改法也适用于注册

`Register.vue` 的 `register` 函数结构完全一样，加上 `try/catch` 即可：

```ts
const register = async () => {
    // 校验非空...
    try {
        const res = await axios.post('/auth/register', { username, password })
        token.value = res.data.token
        localStorage.setItem('token', token.value || '')
        router.push({ name: 'home' })
    } catch (err: any) {
        const msg = err.response?.data?.error || '注册失败，请重试'
        alert(msg)
    }
}
```

## 进阶：用 ElMessage 替代 alert

Element Plus 提供了 `ElMessage`，比原生 `alert` 更好看：

```ts
import { ElMessage } from 'element-plus'

// catch 中：
ElMessage.error(msg)       // 错误提示（红色）
// try 中成功后跳转前：
ElMessage.success('登录成功')  // 成功提示（绿色）
```

`alert` 是浏览器原生弹窗，会阻塞页面；`ElMessage` 是页内通知，消失后自动关闭，体验更好。

---

# 带 token 的请求该怎么写

## 你的项目已经配置好了，不需要额外处理

你的 [axios.ts](frontend/src/axios.ts) 里已经写了一个**请求拦截器**（第 7-15 行）：

```ts
instance.interceptors.request.use(config => {
    const token = localStorage.getItem('token');
    if (token) {
        config.headers = config.headers || {};
        (config.headers as any).Authorization = 'Bearer ' + token;
    }
    return config;
})
```

这个拦截器会在**每次请求发出前**自动执行：
1. 从 `localStorage` 读取 token
2. 如果 token 存在，自动加到请求头 `Authorization: Bearer <token>` 上

所以你在 [News.vue:35](frontend/src/views/News.vue#L35) 写的代码：

```ts
const like = async (id: string) => {
    await axios.post(`articles/${id}/like`)
}
```

**这样写就已经带上 token 了**，不需要做任何额外操作。因为你 import 的是 `../axios`（你自己封装的实例），不是原生的 axios 库。

## 注意：import 的是哪个 axios

在你的项目里：

```ts
import axios from '../axios'   // ✅ 用的是你封装的实例，自动带 token
import axios from 'axios'      // ❌ 用的是原生库，不会自动带 token
```

只要确保 import 路径是 `../axios`（指向你的 `axios.ts`），所有请求都会自动带上 token。

## 你当前代码还需要的改进

虽然 token 的发送没问题了，但 `like` 函数还有几点可以完善：

### 1. 加上错误处理

如果用户没登录就点了点赞，或者 token 过期了，后端会返回 401。目前你的代码没有 catch，出错时用户看不到任何反馈：

```ts
const like = async (id: string) => {
    try {
        await axios.post(`articles/${id}/like`)
        // 点赞成功，可以更新 UI（比如按钮变色、点赞数+1 等）
    } catch (err: any) {
        const msg = err.response?.data?.error || '点赞失败'
        alert(msg)
    }
}
```

### 2. 阻止事件冒泡

你的点赞按钮在 `el-card` 里面，而 `el-card` 也有 `@click` 事件（跳转详情页）。点击点赞按钮时会**同时触发**卡片的点击事件，导致跳转。需要用 `@click.stop` 阻止冒泡：

```html
<el-button class="likes" @click.stop="like(article.ID)">点赞</el-button>
```

### 3. 拼写错误

`calss="likes"` → 应该是 `class="likes"`。

## 总结

| 你关心的 | 现状 |
|----------|------|
| 请求带 token | ✅ 已自动处理，拦截器帮你加了 `Authorization` 头 |
| 错误处理 | ❌ 没有 try/catch，请求失败时用户无感知 |
| 事件冒泡 | ❌ 点赞会同时触发卡片点击，需加 `.stop` |

---

# 点赞请求报 401 的原因

## 根本原因：token 被加了两次 "Bearer " 前缀

全链路追踪如下：

**第一步**：后端 `GenerateJWT` 返回时已经带了 "Bearer " 前缀（[utils.go:25](backend/utils/utils.go#L25)）：

```go
func GenerateJWT(username string) (string, error) {
    // ...
    return "Bearer " + Token, err  // ← 返回 "Bearer eyJhbGciOi..."
}
```

**第二步**：前端把整个字符串（含 "Bearer "）存入 localStorage。

**第三步**：前端拦截器又加一次 "Bearer "（[axios.ts:12](frontend/src/axios.ts#L12)）：

```ts
config.headers.Authorization = 'Bearer ' + token;
// token 已经是 "Bearer eyJhbGciOi..."
// 最终 → "Bearer Bearer eyJhbGciOi..."  ← 两个 Bearer！
```

**第四步**：后端 `ParseJWT` 只去掉一个 "Bearer "（[utils.go:39-41](backend/utils/utils.go#L39-L41)）：

```go
if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
    tokenString = tokenString[7:]  // 去掉一个，剩下 "Bearer eyJhbGci..."
}
// 剩下的 "Bearer eyJhbGci..." 不是合法 JWT → 解析失败 → 401
```

## 解决方案

**改后端**（推荐，改一处即可）：`GenerateJWT` 不要加 "Bearer " 前缀，只返回纯 token：

```go
// backend/utils/utils.go 第 25 行，改为：
return Token, err  // 原来写的是 "Bearer " + Token
```

这样整个链路就各司其职了：

```
后端返回纯 token → 前端存纯 token → 拦截器加上 "Bearer " → 后端解析时去掉 "Bearer " → 正确
```

前后端各自只处理一次 "Bearer "，不会重复。

---

# News 到 Detail 页面的传参方式分析

## 当前实现

**News.vue 第 34 行**：
```ts
router.push({name:'detail', query:{id, title, content}})
```

**Detail.vue 第 8-9 行**：
```html
<h2>{{ route.query.title }}</h2>
<article>{{ route.query.content }}</article>
```

## 结论：方式不正确，存在三个问题

### 问题 1：URL 长度限制

query 参数拼接在 URL 上，例如：

```
/detail?id=1&title=xxx&content=这是一篇很长很长的文章内容...
```

浏览器 URL 长度上限约 **2000 字符**（各浏览器略有差异）。文章 `Content` 稍微长一点就会超出限制，导致参数被截断或跳转失败。

### 问题 2：安全问题

文章内容完全暴露在 URL 中：
- 用户浏览器历史记录会保存完整文章内容
- 如果用户复制 URL 分享给别人，对方直接看到全部内容
- 服务器访问日志也会记录完整 URL（含文章内容）

### 问题 3：编码问题

文章内容可能包含中文、空格、特殊符号（`&`、`=`、`#` 等），这些字符在 URL 中需要编码。Vue Router 会自动处理编码，但如果文章特别长，编码后的 URL 会更长，更容易触及长度上限。

## 正确做法：只传 id，Detail 页自己请求数据

### 为什么只传 id 就够了？

你的后端已经有现成的接口：`GET /api/articles/:id`（[router.go:49](backend/router/router.go#L49)），对应 `GetArticleByID` 控制器。所以 Detail 页面完全可以**只拿到 id，然后自己发请求获取文章数据**。

### 修改 News.vue

```ts
// 只传 id，不要传 title 和 content
const goDetail = (id: string) => {
    router.push({ name: 'detail', query: { id } })
}
```

模板里的点击事件也相应简化：

```html
<el-card @click="goDetail(article.ID)" class="article-card" v-for="article in articles" :key="article.ID">
```

### 修改 Detail.vue

在 `onMounted` 里用 `route.query.id` 发请求获取文章数据：

```html
<template>
    <el-container>
        <el-main>
            <div v-if="article">
                <el-card class="article-card">
                    <h2>{{ article.Title }}</h2>
                    <article>{{ article.Content }}</article>
                    <el-button @click="router.back()">返回</el-button>
                    <footer>&copy; 2026 My App</footer>
                </el-card>
            </div>
            <div v-else-if="loading">加载中...</div>
            <div v-else>文章不存在</div>
        </el-main>
    </el-container>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import axios from '../axios'

const router = useRouter()
const route = useRoute()

interface Article {
    ID: string
    Title: string
    Content: string
}

const article = ref<Article | null>(null)
const loading = ref(true)

onMounted(async () => {
    try {
        const id = route.query.id as string
        const res = await axios.get<Article>(`/articles/${id}`)
        article.value = res.data
    } catch (err: any) {
        console.error('获取文章失败', err)
    } finally {
        loading.value = false
    }
})
</script>
```

### 改后的流程对比

| | 当前做法 | 正确做法 |
|---|---|---|
| URL | `/detail?id=1&title=xxx&content=很长很长的内容...` | `/detail?id=1` |
| URL 长度 | 随文章内容增长，可能超限 | 固定很短 |
| 安全性 | 文章内容暴露在 URL | URL 中只有 id |
| 数据来源 | 从 URL 读取（不可靠） | 从后端 API 获取（可靠） |
| SEO/分享 | 无意义 | URL 干净，可分享 |

## 补充建议

如果想更规范，还可以用**路由 params**代替 query：

```ts
// router/index.ts 中把路由改成
{
    path: '/detail/:id',
    name: 'detail',
    component: () => import('../views/Detail.vue')
}

// News.vue 跳转
router.push({ name: 'detail', params: { id } })

// Detail.vue 读取
const id = route.params.id as string
```

用 params 的 URL 是 `/detail/1`，比 `/detail?id=1` 更简洁、更 RESTful。但这只是风格偏好，query 和 params 都可以，核心关键是**不要把文章内容放进 URL**。

---

# `article.value = res.data` 是直接赋值吗？有隐患吗？

## 是的，这就是一个普通的对象引用赋值

```ts
article.value = res.data
```

这行代码做的事情：把 `res.data` 这个对象的**引用**赋给 `article.value`。之后 `article.value` 和 `res.data` 指向**堆内存中的同一个对象**。

JavaScript 的赋值永远是"传引用"（对象类型），没有隐式的深拷贝。

## 这样做有问题吗？

**在这个场景下没有问题。** 而且是 Vue 3 里最标准的做法。

原因：

1. **axios 每次请求返回的是新对象**：`res.data` 是 axios 从 HTTP 响应体反序列化出来的全新对象，不是复用的缓存。所以不存在"多个地方持有同一个引用导致互相干扰"的问题。

2. **Vue 3 的 ref 会把赋值的对象自动变成响应式**：当你执行 `article.value = res.data`，Vue 内部的 Proxy 会拦截这个赋值，把 `res.data` 包装成响应式对象。之后模板里 `article.Title` 的展示、修改都会自动追踪。

3. **TypeScript 只做编译时检查**：`axios.get<Article>` 告诉 TypeScript "我期望 `res.data` 的形状是 `Article`"，编译器帮你做类型校验。但运行时真实的 `res.data` 可能带着后端返回的全部字段（比如后端多返回了 `Preview`、`CreatedAt` 等），这些字段也一样会出现在 `article.value` 上，不影响功能。

## 什么时候才需要关心这个问题？

| 场景 | 直接赋值 OK？ | 说明 |
|------|:---:|---|
| 从 API 获取新数据后替换展示 | ✅ | 标准做法 |
| 想把一个对象的修改"隔离"（编辑表单草稿） | ❌ | 需要拷贝，否则改草稿会直接改原对象 |
| 多个组件共享同一个状态对象 | ⚠️ | 应该用 store（Pinia），而不是通过 ref 互相传递引用 |

### 什么时候需要拷贝？

如果你在 Detail 页有一个"编辑"功能，用户修改表单时你不想直接改 `article` 对象（防止用户点了取消后数据已经变了），那时候才需要：

```ts
// 需要深拷贝的场景：编辑草稿
const draft = ref<Article>(JSON.parse(JSON.stringify(article.value)))
// 或者更规范的：
const draft = ref<Article>({ ...article.value })  // 浅拷贝够用的话
```

但对于你的 Detail 页——**只展示，不修改**——`article.value = res.data` 完全没问题。这是一次新的 HTTP 响应，一个全新的对象，直接赋值是最简洁正确的写法。

---

# 如果后端字段名和前端 interface 变量名不同怎么办？

## 核心认知：TypeScript interface 不做字段映射

```ts
// 假设后端返回的 JSON：
{ "user_name": "张三", "user_age": 25 }

// 你定义的前端 interface：
interface User {
    name: string    // ❌ 和后端的 user_name 对不上
    age: number     // ❌ 和后端的 user_age 对不上
}

// 执行赋值
user.value = res.data  // res.data 是 { user_name: "张三", user_age: 25 }
```

赋值之后，`user.value` 这个对象的 **key 仍然是后端的 `user_name` 和 `user_age`**，不会魔法般地变成 `name` 和 `age`。

此时你在模板里写：

```html
<p>{{ user.name }}</p>   <!-- 显示空白！因为对象上没有 name 这个 key -->
<p>{{ user.user_name }}</p>  <!-- 能显示"张三"，但你的 interface 没定义这个字段 -->
```

## 一句话总结

`axios.get<User>(...)` 的泛型 `<User>` 只是告诉 TypeScript "你帮我检查代码里有没有访问不存在的字段"。但**真正决定对象上有什么 key 的，是后端返回的 JSON**。

## 三种情况的处理方式

### 情况 1：字段名完全一致（你的项目当前就是）

后端返回 `{ ID, Title, Content }`，前端 interface 也写 `{ ID, Title, Content }`，模板用 `article.Title`。

✅ 完全没问题，直接赋值即可，TypeScript 检查 + 运行时取值都正确。

### 情况 2：字段名不同，你想在前端换名字

后端返回 `{ user_name, user_age }`，你想在前端用 `name` 和 `age`。

这时候**必须手动映射**，不能直接赋值：

```ts
const res = await axios.get('/user/1')
// res.data = { user_name: "张三", user_age: 25 }

// 手动映射
user.value = {
    name: res.data.user_name,
    age: res.data.user_age
}
```

TypeScript 不会帮你做这件事，JSON 解析也不会。

### 情况 3：用 GORM 的 JSON 标签控制后端字段名

你可以在 Go 的 model 里加 `json` 标签，让后端输出的字段名直接匹配前端：

```go
type Article struct {
    ID      uint   `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}
```

这样后端返回的就是 `{ "id": 1, "title": "xxx", "content": "yyy" }`，前端 interface 全部用小写即可。

## 你项目的建议

你后端 Article model 大概率没有写 `json` 标签，所以 GORM 默认用字段名（大写开头：`ID`、`Title`、`Content`）。前端 interface 用同样的名字就能对上，所以当前没有问题。

但更推荐的 Go 风格是：**后端 model 都加上 `json` 标签统一用小写驼峰**，前端 interface 也用小写开头。这样前后端风格一致，以后不会乱。

---

# 后端返回了 Preview 等额外字段，但 interface 里没写

## 实际返回 vs 你定义的 interface

后端 `Article` 模型（[article.go:5-11](backend/models/article.go#L5-L11)）嵌入了 `gorm.Model`，所以 `GetArticleByID` 接口实际返回的 JSON 是：

```json
{
    "ID": 1,
    "CreatedAt": "2026-05-27T...",
    "UpdatedAt": "2026-05-27T...",
    "DeletedAt": null,
    "Title": "文章标题",
    "Content": "文章内容...",
    "Preview": "文章摘要...",
    "Likes": 5
}
```

共 8 个字段。但你定义的 interface 只有 3 个：

```ts
interface Article {
    ID: string
    Title: string
    Content: string
}
```

## 这样做有问题吗？

**没有运行时问题。** `article.value = res.data` 之后，`article.value` 上的实际字段是后端返回的 8 个，不是 interface 定义的 3 个。模板只用到了 `article.Title` 和 `article.Content`，其余 6 个字段存在但没被访问，毫无影响。

TypeScript 的 interface 在这里扮演的是"我只需要这些"的角色，不是"只允许有这些"。这叫做 TypeScript 的**结构化类型**（structural typing）——只要对象包含 interface 要求的字段，多出来的字段不会报错。

## 要不要把 interface 补全？

**不是必须的，但建议加上你实际要用的字段：**

```ts
interface Article {
    ID: number        // ← 建议改成 number，后端返回的是数字
    Title: string
    Content: string
    Preview: string   // ← 加上，后续可能会用到
    // CreatedAt、UpdatedAt、DeletedAt、Likes 用不到就不加
}
```

不补也完全能跑，补了的好处是：
- 以后在代码里访问 `article.Preview` 或 `article.Likes` 时 TypeScript 不会报错
- 让读代码的人一眼就知道这个对象上有哪些字段

## 关键点

`axios.get<Article>` 中的 `<Article>` 泛型只影响 `res.data` 的**类型推断**，不影响运行时对象上的实际字段。你写 `<Article>` 或 `<any>`，`res.data` 上都是后端返回的那 8 个字段，一个不会多一个不会少。

