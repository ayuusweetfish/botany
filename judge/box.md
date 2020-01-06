# Botany: Judge System

## Directory Structure

```
/var/botany
| - compile.sh
| - match.sh
| - submissions/
| | - <sid>/
| | | - code.<lang>
| | | - bin
| | - ...
```

### compile.sh

```
./compile.sh <sid>
```

stdout and stderr are concatenated as the compilation log.

Non-zero return codes denote compilation failures.

This script will be run under resource limits stricter than usual.

### match.sh

```
./match.sh <judge-sid> <mid> <party-sid> <party-sid> ...
```

stdout is the report while stderr is the log.

Non-zero return codes denote internal system errors.

## Steps to build

Set up a chroot jail and log in with `root`.

```sh
mkdir -p /var/botany
cd /var/botany
chown <outside_user> -R .
cp /path/to/compile.sh .
cp /path/to/match.sh .
```

Install: isolate, gcc, lua
