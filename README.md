# Botany

Botany is a customizable contest platform for duels among programs.

## Running the server

Tested Go version: 1.13.1

```sh
cd app
go get -d .
go run .
```

## 规范与约定

提交前使用 Gofmt 格式化代码。

```sh
gofmt -w .
```

Vim 用户可以在 `~/.vimrc` 中加入下列命令，这样在 `app/` 目录下打开文件时会使用目录的缩进设置（4 格宽制表符）。

```vimrc
:set exrc
:set secure
```
