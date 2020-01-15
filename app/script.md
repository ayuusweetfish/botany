# Botany: Contest Scripting

Test the script with [the contest script tool](../tools/script_test).

需求：
- 在每次提交时触发
- 定时触发，约 10 秒执行一次
- 手动触发，可输入字符串参数

参赛者以用户 ID 表示，所有评测以最后一次成功的提交为准

假设比赛管理员可信，提供 Lua 的所有模块（包括 io 和 os）

`create_match()` 接收任意多个用户 ID 作为参数，创建一场对局。若某个 ID 不存在，则忽略整场对局。  
`get_id()` 接收一个登录名，返回对应的用户 ID。若登录名不存在，则返回 0。

另：需要在一场对局结束后结算 rating 和 performance

接收裁判程序产生的 report (string) 以及一个数组，包含若干由 rating (number) 和 performance (string) 组成的选手信息，直接修改之即可，不必返回。由于传入的是引用，调用者 Golang 可以获取结果。

```lua
function on_submission(all, from)
	create_match(from, get_id('SampleAI'))
end

function on_timer(all)
	if os.date('*t').min % 10 ~= 0 then return end
	for i = 1, #all do
		for j = i + 1, #all do
			create_match(all[i], all[j])
		end
	end
end

function on_manual(all, arg)
	if arg == 'match with sample' then
		for i = 1, #all do
			create_match(all[i], get_id('SampleAI'))
		end
	end
end

function update_stats(report, par)
	local id = 1
	if report:sub(1, 1) == 'B' then id = 2 end
	par[id].rating = par[id].rating + 1
	par[id].performance = par[id].performance .. ' hahaha'
end
```
