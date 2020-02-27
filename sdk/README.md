# Botany SDK

## 裁判

裁判用 C/C++ 编写，保存为文件 judge.c/cpp，Make 将自动检测并使用对应的编译器。

裁判在 **main()** 函数接收命令行参数，并由此初始化玩家列表：

```c
int main(int argc, char *argv[])
{
    int n;
    bot_player *pl = bot_player_all(argc, argv, &n);
    // ...
```

**bot_player_all()** 的前两个参数一般只需将命令行参数 argc 和 argv 传入即可。第
三个参数为一 **int** 指针，若它不为 **NULL**，则函数将在它所指向的位置保存玩家
的数量。对于固定人数的游戏，可以传入 **NULL**。

函数返回一个 **bot_player** 数组（指针），代表若干玩家。此后需要与玩家交互时，
都需要指定这个数组中的一个。
