# Botany

master [![Build Status](https://travis-ci.com/kawa-yoiko/botany.svg?token=GcJo7cdxZitoWy9qXz8p&branch=master)](https://travis-ci.com/kawa-yoiko/botany) | dev [![Build Status](https://travis-ci.com/kawa-yoiko/botany.svg?token=GcJo7cdxZitoWy9qXz8p&branch=dev)](https://travis-ci.com/kawa-yoiko/botany)

Botany is a customizable contest platform for duels among programs.

## Starting the database server

Install PostgreSQL

```sh
initdb -D <directory for storage>
pg_ctl -D <directory for storage> start
createdb <database name> -U <database user name>
```

## Running the server

Tested Go version: 1.13.1

Clone repository into `$GOPATH/src/github.com/kawa-yoiko/botany`, or create a symlink, whichever works.

```sh
cd app
go get -d .

cp config_example.json config.json
vim config.json     # Edit in any convenient way

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
