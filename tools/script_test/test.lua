require './tool'

set_participants(
    'su', 1000, '',
    'quq', 1500, '',
    'qvq', 1600, '',
    'qwq', 1700, '',
    'qxq', 1800, ''
)

require './contest'

test_submission(3)
test_timer()
test_timer()
test_timer()
test_manual('1, 2, 3, 4')
test_update_stats('hahahaha', {1, 2})
test_update_stats('hahahaha', {1, 3})
test_update_stats('hahahaha', {2, 4})
test_update_stats('hahahaha', {1, 2})
