# Botany: Contest Scripting

Test the script with [the contest script tool](../tools/script_test).

## 触发事件，创建对局

- 在每次提交时触发
- 定时触发，约 10 秒执行一次
- 手动触发，可输入字符串参数

参赛者以用户 ID 表示，所有评测以最后一次成功的提交为准

提供 Lua 的所有模块（包括 io 和 os）

`create_match()` 接收任意多个用户 ID 作为参数，创建一场对局。若某个 ID 不存在，则忽略整场对局。  
`get_id()` 接收一个登录名，返回对应的用户 ID。若登录名不存在，则返回 0。  
`get_handle()` 接收一个用户 ID，返回对应的登录名。若不存在，返回空字符串。

- 以下函数的 **all** (table) 参数包含若干元素，每个元素包含如下字段
	- **id** (number)
	- **handle** (string)
	- **rating** (number)
	- **performance** (string)
- **on_submission(all, from)**: 每次提交编译通过后触发，**from** (number) 为提交者的 ID
- **on_timer(all)**: 每隔 2 秒触发
- **on_manual(all, arg)**: 在网页界面手动触发，**arg** (string) 为触发时输入的参数

## 在一场对局结束后结算 rating 和 performance

接收裁判程序产生的 report (string) 以及一个数组，包含若干由 rating (number) 和 performance (string) 组成的选手信息，直接修改之即可，不必返回。由于传入的是引用，调用者 Golang 可以获取结果。
