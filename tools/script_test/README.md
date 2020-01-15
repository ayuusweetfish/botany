# Botany: Contest script debugging tool

Put the contest script in `contest.lua`, then edit `test.lua` as follows.

- Create contestants with `set_participants()`
    - Pass in handle, rating and performance
    - Will be assigned sequential IDs starting from 1
- Call `require './contest'`
- Call `test_*()` functions with different test cases

Then invoke `test.lua` and manually check its output. Have fun!
