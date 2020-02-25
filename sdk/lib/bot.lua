return {
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
        io.stderr:flush()
    end
}
