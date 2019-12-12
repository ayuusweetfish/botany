# Botany: Message Queue Protocol (Internal)

借助 Redis 的 Stream 数据结构实现。

## 编译 (key: compile)

**Backend** Startup
```
XGROUP CREATE compile compile_group 0 MKSTREAM
```

**Backend** Send compilation
```
XADD compile * sid <submission-id> contents <code>
```

**Judge** Claim compilation
```
XREADGROUP GROUP compile_group compile_worker_<cwid> COUNT 1 BLOCK 1000 STREAMS compile >
```

**Judge** Update compilation
```
RPUSH compile_result <submission-id> 1 <message>
```

**Judge** Finish compilation
```
RPUSH compile_result <submission-id> 9 <message>
XACK compile compile_group <redis-id>
```
