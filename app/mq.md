# Botany: Message Queue Protocol (Internal)

借助 Redis 的 Stream 数据结构实现。

**Backend** Startup
```
XGROUP CREATE compile judge_group 0 MKSTREAM
XGROUP CREATE match judge_group 0 MKSTREAM
```

**Backend** Send
```
XADD compile * sid <sid>
-- or --
XADD match * mid <mid> num_parties <count> party_1 <sid> party_2 <sid> ...
```

**Judge** Claim
```
XREADGROUP GROUP judge_group judge_<cwid> COUNT 1 BLOCK 1000 STREAMS compile match > >

1) 1) "compile"
   2) 1) 1) "1576381626499-0"
         2) 1) "sid"
            2) "1"
-- or --
1) 1) "match"
   2) 1) 1) "1576381668902-0"
         2) 1) "mid"
            2) "1"
            3) "num_parties"
            4) "2"
            5) "party_1"
            6) "10"
            7) "party_2"
            8) "11"
```

**Judge** Update
```
RPUSH compile_result <sid> 1 <message>
-- or --
RPUSH match_result <mid> 1 <report>
```

**Judge** Finish compilation
```
RPUSH compile_result <sid> 9 <message>
XACK compile judge_group <redis-id>
-- or --
RPUSH match_result <mid> 9 <report>
XACK match judge_group <redis-id>
```
