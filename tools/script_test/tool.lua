inspect = require('./inspect/inspect')

local ids = {}
local map_to_id = {}
local map_to_handle = {}

local stats = {}

get_id = function (handle) return map_to_id[handle] or 0 end
get_handle = function (id) return map_to_handle[id] or '' end

set_participants = function (...)
    local a = { ... }
    for i = 1, #a - 2, 3 do
        local id, handle, rating, perf = #ids + 1, a[i], a[i + 1], a[i + 2]
        ids[id] = id
        map_to_id[handle] = id
        map_to_handle[id] = handle
        stats[id] = { rating = rating, performance = perf }
    end
end

local matches = {}

create_match = function (...)
    local a = { ... }
    matches[#matches + 1] = a
end

--------

local test = function (fn_name, fn, arg)
    print('==== Running on_' .. fn_name .. '(' .. (arg and tostring(arg) or '') .. ') ====')
    local all = {}
    for _, id in pairs(ids) do
        local st = stats[id]
        all[#all + 1] = {
            id = id,
            handle = get_handle(id),
            rating = st.rating,
            performance = st.performance
        }
    end
    fn(all, arg)
    print('==== ' .. tostring(#matches) .. ' match' ..
        (#matches == 1 and '' or 'es') .. ' created ====')
    for i, m in ipairs(matches) do
        print(string.format('#%2d: %s', i, inspect(m)))
    end
    matches = {}
    print('')
end

test_submission = function (from) test('submission', on_submission, from) end
test_timer = function () test('timer', on_timer, nil) end
test_manual = function (arg) test('manual', on_manual, arg) end

test_update_stats = function (report, par_ids)
    print('==== Running update_stats(<report>, ' .. inspect(par_ids) .. ') ====')
    local orig_par, par = {}, {}
    for i, id in ipairs(par_ids) do
        orig_par[i] = stats[id]
        par[i] = {
            id = id,
            rating = stats[id].rating,
            performance = stats[id].performance
        }
    end
    update_stats(report, par)
    print('==== Report ====')
    print(report)
    print('==== Results ====')
    for i = 1, #par do
        print(string.format('%10s: %d -> %d | %s -> %s',
            map_to_handle[par_ids[i]],
            orig_par[i].rating, par[i].rating,
            orig_par[i].performance, par[i].performance))
        stats[par_ids[i]] = par[i]
    end
    print('')
end
