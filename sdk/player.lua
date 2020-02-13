bot = require('bot')

math.randomseed(os.time())

local s = bot.recv()
bot.log('Hello from the Moon! Side #' .. s .. '\n')

local used = {}
while true do
    local r, c = bot.recv():match('(%S+) (%S+)')
    local id = tonumber(r) * 3 + tonumber(c) + 1
    used[id] = true

    local sel
    repeat sel = math.random(9) until not used[sel]
    used[sel] = true

    bot.send(string.format('%d %d', math.floor((sel - 1) / 3), (sel - 1) % 3))
end
