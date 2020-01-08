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
| - matches/
| | - <mid>/
| | | - 0.log
| | | - 1.log
| | | - ...
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
./match.sh <mid> <judge-sid> <party-sid> <party-sid> ...
```

stdout is the report while stderr is the log.

Non-zero return codes denote internal system errors.

## Steps to build

Set up a chroot jail and log in with `root`.

```sh
mkdir -p /var/botany/submissions
chown 1000 /var/botany/submissions
chgrp 1000 /var/botany/submissions
mkdir -p /var/botany/matches
chown 1000 /var/botany/matches
chgrp 1000 /var/botany/matches
cd /var/botany
chown <outside_user> -R .
cp /path/to/compile.sh .
cp /path/to/match.sh .
```

Install: isolate, gcc, lua

Requires the host environment to have the cURL binary in $PATH.

### Alpine mini-rootfs

```sh
# Move projects inside
mkdir -p alpine/var/botany
cp botany/judge/compile.sh botany/judge/match.sh alpine/var/botany
git clone https://github.com/ioi/isolate.git alpine/home/isolate

# Create mount points
sudo mount --bind alpine alpine
sudo mount --bind /proc alpine/proc # TODO: Do not mount /proc entirely

# Set up network
sudo cp /etc/resolv.conf alpine/etc

# Enter chroot
sudo chroot alpine sh

# Install packages
sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
apk add gcc make libc-dev libcap-dev

# Build isolate
cd /home/isolate
make
make install

# Test isolate
isolate --init
isolate --run -- /bin/echo hi
isolate --cleanup
```

In the host environment, run Botany's judge side:

```sh
sudo apt-get install libhiredis-dev libb2-dev

cd botany/judge
sudo ./build.sh
./a.out -i 1 -d /path/to/alpine
```
