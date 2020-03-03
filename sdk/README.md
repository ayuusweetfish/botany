# Botany SDK

## 使用

当前版本需要在 Unix-like 环境（Unix, Linux, macOS, WSL 等）下运行，依赖 GCC 与 Make 工具。

```sh
# 运行裁判，裁判将运行两个同样的选手程序
make
make run

# 运行 Lua 语言选手程序
make run-lua

# 运行 Python 语言选手程序
make run-py

# 只编译裁判
make judge

# 只编译选手
make player

# 移除生成的文件
make clean
```

## 井字棋

当前版本附带了井字棋游戏的裁判和选手程序，实现较简单，可供参考。

在命令行执行 `make` 即可编译运行。

## 选手

这里只介绍选手的 C/C++ 接口。选手程序保存为文件 player.c/cpp，Make 将自动检测并使用对应的编译器。

函数签名：

```c
void bot_send(const char *s);
char *bot_recv();
```

通过 __bot_send()__ 发送消息，__bot_recv()__ 接收消息。所有消息都是包含除 NUL 字符 __\\0__ 外任何字符的长度严格小于 16 MiB 的字符串。

接收到的消息字符串需要通过 __free()__ 释放。

任意时刻都可以自行向标准错误 stderr 输出日志信息，这些内容可以在对局记录页面下载查看。推荐每次输出内容之后调用 __fflush()__，以避免在程序被裁判停止时丢失部分日志。

```
fprintf(stderr, "Hello world!\n");
fflush(stderr);
```

祝大家玩得开心！

## 裁判

裁判用 C/C++ 编写，保存为文件 judge.c/cpp，Make 将自动检测并使用对应的编译器。

### 初始化

函数签名：

```c
int bot_judge_init(int argc, char *const argv[]);
```

一般只需将命令行参数 __argc__ 和 __argv__ 传入 __bot_judge_all()__，它会返回一个整数 __n__，表示玩家的数量，玩家编号为 __0, 1, …, n-1__。

```c
int main(int argc, char *argv[])
{
    int n = bot_judge_init(argc, argv);
    // ...
```

### 互传信息

函数签名：

```c
void bot_judge_send(int id, const char *str);
char *bot_judge_recv(int id, int *o_len, int timeout);
```

通过 __bot_judge_send()__ 向玩家发送消息。需要指定一个玩家编号与一个字符串，字符串可以包含除了 NUL 字符 __\\0__ 以外的任意字符。

```c
bot_judge_send(0, "0");
bot_judge_send(1, "1");
```

若要发送动态生成的消息，可以利用 __sprintf__ 或 __snprintf__：

```c
snprintf(buf, sizeof buf, "%d %d", row, col);
bot_judge_send(0, buf);
```

各条消息不会自动拼接，且能以队列形式暂存。例如连续两次调用 __bot_judge_send()__ 之后，玩家两次调用 __bot_recv()__ 会分别收到两次发送的消息。

通过 __bot_judge_recv()__ 从玩家接收信息。需要指定一个玩家编号与一个等待时间，等待时间的单位为毫秒。玩家程序会开始运行，直到它发出一条消息后被暂停；如果队列中已经有消息，则玩家程序立刻被暂停。

玩家正常发出消息时，__bot_judge_recv()__ 返回一个字符串，并在参数 __o_len__ 所指向的位置存放字符串的长度；否则它返回空，并在 __o_len__ 所指向的位置存放一个错误代码。错误代码含义如下：

- __BOT_ERR_CLOSED__ (4): 选手程序自行退出或崩溃；
- __BOT_ERR_TIMEOUT__ (5): 选手程序在等待时间内未发送消息；
- __BOT_ERR_FMT__ (1): 不合法数据，可能选手自行向标准输出输出了内容；
- __BOT_ERR_SYSCALL__ (2): 内部错误，正常状态下不会发生。

另外，接收到的消息字符串需要通过 __free()__ 释放。

### 生成对局报告

裁判程序将对局的报告输出到标准输出 stdout，这些内容后续将传给播放器和负责计算选手积分的赛制脚本程序。

将裁判日志输出到标准错误 stderr，这些内容可供选手和管理员看到没有包含在报告中的一些细节。

一般而言，只需使用 __printf()__ 和 __fprintf(stderr, …)__ 即可。

### 结束

函数签名：

```c
void bot_judge_finish();
```

在裁判程序结束前调用 __bot_judge_finish()__，以结束选手程序并将它们的日志写入文件。
