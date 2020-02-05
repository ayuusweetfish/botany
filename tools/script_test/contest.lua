local count = 2
local su_id = get_id('su')

function on_submission(all, from)
    print('Submission', from)
end

function on_timer(all)
    count = count + 1
    if count < 3 then return end
    count = 0
    print('Superuser has ID ' .. tostring(su_id))
    print('Creating matches for contest #0')
    print('Number of participants with delegates ' .. tostring(#all))
    for i = 1, #all do
        print(string.format('Contestant %s (%d), rating %d, performance "%s"',
            all[i].handle, all[i].id, all[i].rating, all[i].performance))
        if i > 1 then create_match(all[i].id, all[i - 1].id) end
    end
end

function on_manual(all, arg)
    print('Manual', arg)
end

function update_stats(report, par)
    print('Update with ' .. tostring(#par) .. ' parties')
    print(report)
    for i = 1, #par do
        print(i, par[i].rating, par[i].performance)
        par[i].rating = par[i].rating + 1
        par[i].performance = 'Took part in ' .. tostring(par[i].rating) .. ' match'
        if par[i].rating ~= 1 then par[i].performance = par[i].performance .. 'es' end
    end
end
