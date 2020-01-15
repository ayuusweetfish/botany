# Botany

master [![Build Status](https://travis-ci.com/kawa-yoiko/botany.svg?token=GcJo7cdxZitoWy9qXz8p&branch=master)](https://travis-ci.com/kawa-yoiko/botany) | dev [![Build Status](https://travis-ci.com/kawa-yoiko/botany.svg?token=GcJo7cdxZitoWy9qXz8p&branch=dev)](https://travis-ci.com/kawa-yoiko/botany)

Botany is a customizable contest platform for games among programs.

## Development tools for jury and participants

See [README in sdk/ directory](sdk/README.md).

## Running the server

Tested Go version: 1.13.1

Clone repository into `$GOPATH/src/github.com/kawa-yoiko/botany`, or create a symlink, whichever works.

Set up PostgresQL and Redis server (refer to documentations for the target platform).

```sh
cd app
go get -d .

cp config_example.json config.json
vim config.json     # Edit in any convenient way

go run .
```

In case of unstable connection, use [goproxy.io](https://goproxy.io/) or try:

```sh
mkdir -p $GOPATH/src/golang.org/x
cd $GOPATH/src/golang.org/x/
git clone https://github.com/golang/crypto
git clone https://github.com/golang/image
git clone https://github.com/golang/sys
```

## Running the judge

The judge runs on Linux only.

See [judge/box.md](judge/box.md) for details.

## Standards and conventions

Run the code through Gofmt before committing. Not doing so will result in CI complaining.

```sh
gofmt -w .
```

Vim users may consider adding the following to `~/.vimrc` so that indentation rules (`noexpandtab ts=4 sw=4`) apply inside `app/` directory.

```vimrc
:set exrc
:set secure
```
