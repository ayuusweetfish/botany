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


## 用户

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
- 空对象 {}
- 登录名或密码错误

### 当前帐号 GET /whoami

响应 200
- 一个 UserShort

响应 401
- 空对象 {}
- 未登录

### 个人主页 GET /user/{handle}/profile

响应
- **user** (User) 帐号信息
- :construction: **contests** ([ContestShort]) 参与的比赛列表
- :construction: **matches** ([MatchShort]) 最近对局列表

### :construction: 赋予或撤回主办权限 POST /user/{handle}/promote

请求
- **set** (boolean) true 表示赋予权限，false 表示取消权限

响应 200
- 空对象 {}
- 无论前后权限是否变化，都正常返回

响应 403
- 空对象 {}
- 无站长权限 —— 前端检查严格时不应出现此项


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
- **my_role** (number) 自己的参加情况
	- **-1** 未登录或未报名
	- **0** 作为选手参加
	- **1** 拥有管理权限（管理员或创建者）

### 比赛数据结构 ContestShort

仅包含 Contest 的 **id**, **title**, **banner**, **start_time**, **end_time**, **desc**, **is_reg_open**

### 提交记录数据结构 Submission

- **id** (number) ID
- **participant** (UserShort) 参赛者
- **created_at** (number) 提交时刻的 Unix 时间戳，单位为秒
- **status** (number) 状态
	- **0** 等待处理
	- **1** 正在编译
	- **9** 接受
	- **-1** 编译错误
	- **-9** 系统错误（请联系管理员）
- **msg** (string) 编译信息
- **contents** (string) 代码

### 提交记录数据结构 SubmissionShort

不包含 Submission 的 **msg** 和 **contents**

### 对局数据结构 Match

- **id** (number) ID
- **parties** 参与对局的各方，每个元素如下
	- **submission** (SubmissionShort) 提交记录
	- **score** (number) 本场得分
	- **is_winner** (boolean) 是否获胜
- **report** (string) 对局报告，交给动画播放器

### 对局数据结构 MatchShort

不包含 Match 的 **details**

### :construction: 创建比赛 POST /contest/create

请求
- **title** (string) 标题
- **banner** (string) Banner 图片链接（暂不使用）
- **start_time** (number) 开始时刻的 Unix 时间戳，单位为秒
- **end_time** (number) 结束时刻的 Unix 时间戳，单位为秒
- **desc** (string) 简要描述
- **details** (string) 长篇详细说明
- **is_visible** (boolean) 是否公开显示
- **is_reg_open** (boolean) 是否公开接受报名

响应 200
- 空对象 {}
- 无论是否变更可见性都正常返回

响应 400
- 空对象 {}
- 任何一项信息过长、过短或不合法 —— 前端检查严格时不应出现此项

响应 403
- 空对象 {}
- 无主办权限 —— 前端检查严格时不应出现此项

### :construction: 公开或隐藏比赛 POST /contest/{cid}/publish

请求
- **set** (boolean) true 表示设为公开，false 表示设为隐藏

响应 200
- 空对象 {}
- 无论前后可见性是否变化，都正常返回

响应 403
- 空对象 {}
- 无站长权限 —— 前端检查严格时不应出现此项

### 比赛列表 GET /contest/list

响应
- 若干 ContestShort 组成的数组

### 比赛信息 GET /contest/{cid}/info

响应
- 一个 Contest

### 比赛报名 POST /contest/{cid}/join

响应 200
- 空对象 {}
- 报名成功或此前已报名

响应 400
- 空对象 {}
- 比赛不开放报名 —— 前端处理正确时不应出现此项

### 自己的提交历史 GET /contest/{cid}/my

响应
- 若干 SubmissionShort 组成的数组，从最新到最旧排序

### 提交详情 GET /contest/{cid}/submission/{sid}

响应 200
- 一个 Submission

响应 403
- 空对象 {}
- 普通参赛者不能在比赛期间查看不属于自己的提交 —— 前端处理正确时不应出现此项

### 提交代码 POST /contest/{cid}/submit

请求
- **code** (string) 代码

响应 200
- **err** (number) 0
- **submission** (SubmissionShort) 本次提交记录

响应 400
- **err** (number) 错误代码
	- **1** 代码超过长度限制
	- **2** 代码包含不合法字符

响应 403
- 空对象 {}
- 未报名比赛或比赛未开始 —— 前端检查严格时不应出现此项

### :construction: 排行榜 GET /contest/{cid}/ranklist

响应
- 一个数组，按排名从高到低排序，每个元素如下
	- **participant** (UserShort) 参赛者
	- **win_count** (number) 胜场数
	- **lose_count** (number) 败场数
	- **rating** (number) 匹配积分

### :construction: 对局列表 GET /contest/{cid}/matches

响应
- 若干 MatchShort 组成的数组，从最新到最旧排序

### :construction: 对局详情 GET /contest/{cid}/match/{mid}

响应
- 一个 Match
