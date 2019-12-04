# Botany: API (Internal)

## 汉英名称规范

| zh         | en          | 说明                       | 避免               |
| ---------- | ----------- | -------------------------- | ------------------ |
| 用户       | User        | 帐号                       | Account            |
| 登录名     | Handle      | 帐号惟一标识符，字母数字   | Username, 用户名   |
| 昵称       | Nickname    | 可以随意更改的显示名       | Username, 用户名   |
| 签名       | Bio         | 个性签名，大家都懂         | Signature          |
| 选手       | Participant | 参加了某场比赛的帐号       | Player, User       |
| 提交       | Submission  | 选手提交代码的一次记录     | Code               |
| 比赛       | Contest     | 选手间的比赛               | Game               |
| 对局       | Match       | 程序间的对局               | Game, Versus, 评测 |
| 对局的一方 | Party       | 参加对局的提交             | Player, Side       |
| 站长权限   | Superuser   | 全站最高管理权限           | —                  |
| 主办权限   | Organizer   | 创建比赛的权限             | —                  |
| 比赛拥有者 | Owner       | 创建了某场比赛的帐号       | Creator, Organizer |
| 比赛管理员 | Moderator   | 可以修改某场比赛信息的帐号 | Manager            |
| 登录       | Log in      | —                          | 登陆               |
| 注册       | Sign up     | —                          | Register           |
| 简介       | Description | —                          | Brief, Abstract    |
| 报名       | Join        | —                          | Enter              |
| 排行       | Ranklist    | —                          | Ranking, Board     |
| 主页       | Profile     | —                          | Personal           |
| 发布比赛   | Publish     | —                          | Publicize, 公开    |

不用 Username 的原因是希望区分登录名 Handle 和昵称 Nickname。  
不用 Player 的原因是希望区分比赛的人 Participant 和对局的代码 Party。


## 普适

请求参数均为 application/x-www-form-urlencoded 格式 `key=val&key=val`，响应内容均为 JSON。

状态码：
- 任何时候返回 401 表示未登录，应该重定向到登录页面，登录后返回当前页面。
- 任何时候返回 404 表示内容不存在或无权访问。
- 任何时候返回 429 表示请求频率过高，人类快速操作并不会导致此情形。
- 任何时候返回 500 表示服务器内部错误。
- 其余状态码由每个接入点各自规定。


## 登录/注册

### 用户数据结构 User

- **id** (number) ID
- **handle** (string) 登录名
- **email** (string) email
- **privilege** (number) 权限
	- **0** 普通权限
	- **1** 主办权限
	- **2** 站长权限
- **joined_at** (number) 加入时刻的 Unix 时间戳，单位为秒
- **nickname** (string) 昵称
- **bio** (string) 个性签名

### 用户数据结构 UserShort

仅包含 User 的 **id**, **handle**, **privilege**，**nickname**

### 验证码 GET /captcha

响应 200
- **key** (string) 标识符
- **img** (string) Base64 编码的图像

### 注册 POST /signup

请求
- **handle** (string) 登录名
- **password** (string) 不可逆哈希后的密码
- **email** (string) 电子邮箱
- **nickname** (string) 昵称
- **captcha_key** (string) 此前获得的验证码 key
- **captcha_value** (string) 验证码填写值

响应 200
- **err** ([number]) 空数组 []

响应 400
- **err** ([number]) 错误列表
	- **1**	用户名重复
	- **2** 邮箱重复
	- **3** 验证码不正确
	- **-1** 任何一项信息过长、过短或包含不合法字符 —— 前端检查严格时不应出现此项

### 登录 POST /login

请求
- **handle** (string) 登录名
- **password** (string) 不可逆哈希后的密码

响应 200
- 一个 UserShort
- 会设置 Cookie，以后的请求不必特殊处理

响应 400
- 空对象 {}，登录名或密码错误


## 比赛

### 比赛数据结构 Contest

- **id** (number) ID
- **title** (string) 标题
- **banner** (string) Banner 图片链接（暂不使用）
- **start_time** (number) 开始时刻的 Unix 时间戳，单位为秒
- **end_time** (number) 结束时刻的 Unix 时间戳，单位为秒
- **desc** (string) 简要描述
- **details** (string) 长篇详细说明
- **is_reg_open** (boolean) 是否公开接受报名
- **owner** (UserShort) 创建者

### 比赛数据结构 ContestShort

仅包含 Contest 的 **id**, **title**, **banner**, **start_time**, **end_time**, **desc**, **is_reg_open**

### 比赛列表 GET /contest/list

响应
- 若干 ContestShort 组成的数组

### 比赛信息 GET /contest/{id}/info

响应
- 一个 Contest
