bot = {
    send = function (s)
        local len = #s
        if len >= 65536 * 256 then
            io.stderr('Message too long!')
            return
        end
        local a, b, c = len % 256, math.floor(len / 256) % 256, math.floor(len / 65536)
        local t = string.char(a, b, c)
        io.stdout:write(t)
        io.stdout:write(s)
        io.stdout:flush()
    end,

    recv = function ()
        local s = io.stdin:read(3)
        local a, b, c = string.byte(s, 1, 3)
        local t = io.stdin:read(a + (b * 256) + (c * 65536))
        return t
    end,

    log = function (s)
        io.stderr:write(s)
    end
}

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
