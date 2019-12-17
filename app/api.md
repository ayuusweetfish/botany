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
| 发布比赛   | Publish     | —                          | Publicize          |

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

### 退出登录 POST /logout

响应 200
- 空对象
- 会清除对应 Cookie，以后的请求不必特殊处理

### 当前帐号 GET /whoami

响应 200
- 一个 UserShort

响应 400
- 空对象 {}
- 未登录

### 个人主页 GET /user/{handle}/profile

请求

- **page** (optional number) 请求的页数
- **count** (optional number) 每页的个数

响应

- **user** (User) 帐号信息
- **contests** ([ContestShort]) 参与的比赛列表（不需要分页）
- **matches** ([MatchShort]) 最近对局列表（需要分页）
- **total_matches** (number) 对局列表总数

### 修改个人信息 POST /user/{handle}/profile/edit

请求
- **email** (string) 电子邮箱
- **nickname** (string) 昵称
- **bio** (string) 个性签名

响应 200
- 空对象 {}

响应 403 —— 前端处理正确时不应出现此项
- 空对象 {}
- 除站长外，不能修改其他人的个人信息

响应 400
- **err** ([number]) 错误列表
	- **1** 邮箱重复
	- **-1** 任何一项信息过长、过短或包含不合法字符(例如邮箱格式错误) —— 前端检查严格时不应出现此项

### 修改密码 POST /user/{handle}/password

请求
- **old** (string) 原密码
- **new** (string) 新密码

响应 200
- 空对象 {}

响应 400
- 空对象 {}
- 原密码错误

响应 403 —— 前端处理正确时不应出现此项
- 空对象 {}
- 除站长外，不能修改其他人的密码

### 赋予或撤回主办权限 POST /user/{handle}/promote

请求
- **set** (boolean) true 表示赋予权限，false 表示取消权限

响应 200
- 空对象 {}
- 无论前后权限是否变化，都正常返回

响应 403
- 空对象 {}
- 无站长权限 —— 前端检查严格时不应出现此项

### 按登录名搜索 GET /user_search/{handle}

响应
- 若干 UserShort 组成的数组
- 登录名包含给定字符串的用户，至多返回 5 个


## 比赛

### 比赛数据结构 Contest

- **id** (number) ID
- **title** (string) 标题
- **banner** (string) Banner 图片链接（暂不使用）
- **start_time** (number) 开始时刻的 Unix 时间戳，单位为秒
- **end_time** (number) 结束时刻的 Unix 时间戳，单位为秒
- **desc** (string) 简要描述
- **details** (string) 长篇详细说明
- **is_visible** (boolean) 是否公开显示
- **is_reg_open** (boolean) 是否公开接受报名
- **script** (string) 赛制脚本
- **owner** (UserShort) 创建者
- **moderators** ([number]) 管理员的 ID 列表，用逗号分隔
- **my_role** (number) 自己的参加情况
	- **-1** 未登录或未报名
	- **0** 拥有管理权限（管理员或创建者）
	- **1** 作为选手参加

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
- **lang** (string) 程序语言
- **contents** (string) 代码

### 提交记录数据结构 SubmissionShort

不包含 Submission 的 **msg** 和 **contents**

### 对局数据结构 Match

- **id** (number) ID
- **parties** ([SubmissionShort]) 参与对局的各方，每个元素为一个提交记录
- **status** (number) 状态
	- **0** 等待处理
	- **1** 正在运行
	- **9** 完成
	- **-9** 系统错误（请联系管理员）
- **report** (object) 对局报告，交给动画播放器

### 对局数据结构 MatchShort

不包含 Match 的 **report**

### 创建比赛 POST /contest/create

请求
- **title** (string) 标题
- **banner** (string) Banner 图片链接（暂不使用）
- **start_time** (number) 开始时刻的 Unix 时间戳，单位为秒
- **end_time** (number) 结束时刻的 Unix 时间戳，单位为秒
- **desc** (string) 简要描述
- **details** (string) 长篇详细说明
- **is_visible** (boolean) 是否公开显示
- **is_reg_open** (boolean) 是否公开接受报名
- **moderators** ([number]) 管理员的 ID 列表，用逗号分隔
- :construction: **script** (string) 赛制脚本

响应 200
- **id** (number) 新比赛的 ID

响应 400
- 空对象 {}
- 任何一项信息过长、过短或不合法 —— 前端检查严格时不应出现此项

响应 403
- 空对象 {}
- 无主办权限 —— 前端检查严格时不应出现此项

### 修改比赛 POST /contest/{cid}/edit

请求
- 同 /contest/create

响应 200
- 空对象 {}

响应 400
- 空对象 {}
- 任何一项信息过长、过短或不合法 —— 前端检查严格时不应出现此项

响应 403
- 空对象 {}
- 不是比赛拥有者 —— 前端检查严格时不应出现此项

### 发布或隐藏比赛 POST /contest/{cid}/publish

请求
- **set** (boolean) true 表示设为公开显示，false 表示设为隐藏

响应 200
- 空对象 {}
- 无论前后可见性是否变化，都正常返回

响应 403
- 空对象 {}
- 无站长权限 —— 前端检查严格时不应出现此项

### 比赛列表 GET /contest/list

响应
- 若干 ContestShort 组成的数组
- 只返回对当前用户可见的比赛

### 比赛信息 GET /contest/{cid}/info

响应
- 一个 Contest

### 比赛报名 POST /contest/{cid}/join

响应 200
- 空对象 {}
- 报名成功或此前已报名

响应 403
- 空对象 {}
- 比赛不开放报名 —— 前端处理正确时不应出现此项

### 自己的提交历史 GET /contest/{cid}/my

响应 200
- 若干 SubmissionShort 组成的数组，从最新到最旧排序

响应 403
- 空数组 []
- 未报名比赛或比赛未开始 —— 前端检查严格时不应出现此项

### 所有提交历史 GET /contest/{cid}/submission/list

请求
- :construction: **page** (optional number) 请求的页数
- :construction: **count** (optional number) 每页的个数

响应 200
- **total** (number) 提交历史总数
- **submissions** ([SubmissionShort])若干 SubmissionShort 组成的数组，从最新到最旧排序

响应 403
- 空数组 []
- 比赛未开始 —— 前端检查严格时不应出现此项

### 提交详情 GET /contest/{cid}/submission/{sid}

响应 200
- 一个 Submission

响应 403
- 空对象 {}
- 普通参赛者不能在比赛期间查看不属于自己的提交 —— 前端处理正确时不应出现此项

### 提交代码 POST /contest/{cid}/submit

请求
- **lang** (string) 程序语言
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

### 选择出战提交 POST /contest/{cid}/delegate

请求
- **submission** (number) 提交 ID
	- 必须是编译通过的提交
	- 如果等于 -1 那么视为暂时退赛

响应 200
- 空对象 {}

响应 400
- 空对象 {}
- 格式不正确

响应 403
- 空对象 {}
- 未报名比赛或比赛未开始，或不是自己的编译通过的提交 —— 前端检查严格时不应出现此项

### 排行榜 GET /contest/{cid}/ranklist

请求
- :construction: **page** (optional number) 请求的页数
- :construction: **count** (optional number) 每页的个数

响应 200
- :construction: **total** 当前玩家总数
- :construction: **participants** 一个数组，按排名从高到低排序，如遇并列则按选手登录名字典序升序排列。每个元素如下
	- **participant** (UserShort) 参赛者
	- **rating** (number) 匹配积分
	- **performance** (string) 额外战绩数据

### 对局列表 GET /contest/{cid}/matches

请求
- :construction: **page** (optional number) 请求的页数
- :construction: **count** (optional number) 每页的个数

响应
- :construction: **total** (number) 对局总数
- :construction: **matches** ([MatchShort]) 若干 MatchShort 组成的数组，从最新到最旧排序

### 对局详情 GET /contest/{cid}/match/{mid}

响应
- 一个 Match

注：比赛不存在、对局不存在或对局不属于比赛均认为 404

### 手动发起对局 POST /contest/{cid}/match/manual

请求
- **submissions** (string) 参与对局的各提交记录编号，用逗号分隔，如 `1,2,3`

响应 200
- 一个 MatchShort

响应 400
- 空对象 {}
- 格式不正确

响应 403
- 空对象 {}
- 非管理员不能手动发起对局

### :construction: 赛制脚本日志 GET /contest/{cid}/match/script_log

响应 200
- 大量纯文本

响应 403
- 空响应 Content-Length: 0
- 非管理员不能查看日志

### :construction: 手动执行赛制脚本 POST /contest/{cid}/match/manual_script

请求
- **arg** (string) 需要传递的参数

响应 200
- 空对象 {}

响应 403
- 空对象 {}
- 非管理员不能手动执行脚本


## 调试

### 复位数据库并填充测试数据 POST /fake

响应 200

测试用数据包括：
- 用户
	- 站长权限：登录名 su
	- 主办权限：登录名 o1 … o5
	- 普通权限：登录名 p1 … p20
	- 密码均为 qwq
- 5 场比赛，其中 1 不公开显示，5 不开放公开报名
	- 其中第 i 场比赛的拥有者是 o\<i\>，初始没有其他管理员
	- 每一场比赛都含有一些伪随机的参赛者、提交和对局
	- 尚未开始的比赛确实是不应该已有对局的，但是对逻辑没有实质影响
