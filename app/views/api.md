## RESTful API 要求

以下省略了JSON格式



在任何时候，返回\_\_\_代表登录失效，返回\_\_\_代表权限错误
返回500表示服务器内部错误


### 登录/注册

* 进入登录页面时，发送GET /captcha/login，要求在pic字段下返回验证码图片的base64编码

* 点击“登录”按钮，发送POST /login

  ```
  username: 	string,
  password:		string,
  captcha:		string
  ```

  要求登录成功后返回200，以及

  ```
  uid:          int(用户唯一认证)
  usertype:		string(用户权限)
  ```

  登录失败时返回403，以及 

  ```
  error:			string(错误信息: 用户名不存在或密码错误-'wrong username or password'/验证码错误-'wrong captcha')
  ```

* 进入注册页面时，发送GET /captcha/register，要求在pic字段下返回验证码图片的base64编码

* 点击“注册”按钮，发送POST /register

  ```
  username:		string,
  password:		string,
  email:		string,
  captcha:		string
  ```

  要求注册成功后返回200

  注册失败后返回403，以及

  ```
  error:		string(错误信息: 用户名已被注册-'username already exists'/邮箱已被注册-'email already exists'/验证码错误-'wrong captcha')
  ```



### 比赛列表

* 进入页面时，发送GET /gamelist

  要求返回200，以及

  ```
  total:		string(比赛总数)
  games:[{
  		id:							string(比赛在服务器中的唯一标识符),
  		name:						string,
  		time_start:			string,
  		time_end:				string(与上述两个时间格式都是yyyy-MM-dd),
  		info:						string(简要说明)
  	},...
  ]
  ```



### 一项赛事

以下发送的GET和POST请求都会带以下参数

```
uid             int(用户在服务器中唯一标示符）
id:				string(比赛在服务器中的唯一标识符)
```



#### 首页

* 进入界面时，发送GET /gamemain

  要求返回200，以及

  ```
  name:				string,
  owner:			string(主办组织名),
  banner:			base64 code(条幅图片),
  brief:			string(比赛简介),
  timeline:[{
  		time:		string(yyyy-MM-dd),
  		info:		string(该时间的安排),
  		status:	string(已结束/进行中/待定等)
  	},...
  ]
  ```

* 报名，发送POST /joingame

  要求成功时返回200，失败时返回\_\_\_，以及

  ```
  error:			string(报名失败的原因)
  ```

#### 参赛指南

* 进入界面，发送GET /gamehelp

  要求返回

  ```
  detail:			string(详细的介绍)
  ```

#### 代码

* 进入界面，发送GET /gamecode

  要求返回

  ```
  status:		string(当前代码的状态：提交/处理/通过等)
  code:			string(当前代码)
  history:[{
  		index:	string(代码在服务器中的唯一标识符),
  		time:		string(yyyy-MM-dd),
  		info:		string(状态)
  	},...
  ]
  ```

* 提交代码，POST /gamecode/submit

  参数

  ```
  code:		string(用户代码)
  ```

* 保存草稿，POST /gamecode/draft

  参数同上

* 导出草稿，GET /gamecode/draft

  返回值同上

* 导出历史代码，GET /gamecode/previous

  参数为

  ```
  index:			string(代码在服务器中的唯一标识符)
  ```

  返回值同上

#### 排行榜

* 刷新，GET /rankinglist

  返回值为

  ```
  total:		string(选手总数),
  myrank:		string(用户排名),
  rankings:[{
  		ranking:	string,
  		avatar:		base64 code(头像),
  		name:			string,
  		email:		string,
  		win:			string,
  		loss:			string,
  		winrate:	string(x%),
  		MMR:			string(匹配分)
  	},...
  ]
  ```

#### 对局列表

* 刷新，GET /versuslist

  返回值为

  ```
  total:		string(总数),
  myvs:			string(该用户参与了几场),
  versus:[{
  		id:				string(对局在服务器中的唯一标识符),
  		player1:{
  			name:		string,
  			avatar: base64 code
  		},
  		player2:{
  			name:		string,
  			avatar:	base64 code
  		},
  		winner:		string(胜者用户名)
  	},...
  ]
  ```

####对局详情

* 刷新，GET /versusinfo

  参数除比赛id以外还有

  ```
  versusid:			string(对局的id)
  ```

  要求返回

  ```
  player1:{
  	name:			string,
  	avatar:		base64 code,
  	email:		string
  },
  player2:{
  	name:			string,
  	avatar:		base64 code,
  	email:		string
  },
  winner:			string(1或2),
  replay:			object
  ```



### 用户详情

* 刷新，GET /personal，参数如下

  ```
  username:			string(请求谁的信息)
  ```

  要求返回

  ```
  avatar:				base64 code,
  self:					string(bool,是否是本人)
  email:				string,
  otherinfo:		string(个人签名等),
  highlight:[
  	string, string, ... (成就列表)
  ],
  game:{
  	total:			string,
  	list:[{
        id:				string,
        name:			string,
        win:			string,
        loss:			string,
        winrate:	string,
        MMR:			string
  		},...
  	]
  },
  versus:{
  	total:		string,
  	list:[{
  			id:				string,
  			opponent:{
  				name:		string,
  				avatar:	string
  			},
  			result:		string(胜负)
  		},...
  	]
  }
  ```

  

