# Botany: Contest script debugging tool

Put the contest script in `contest.lua`, then edit `my.lua` as follows.

- Create contestants with `set_participants()`
    - Pass in handle, rating and performance
    - Will be assigned sequential IDs starting from 1
- Call `require './contest'`
- Call `test_*()` functions with different test cases and check log output
