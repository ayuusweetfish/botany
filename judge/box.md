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
./match.sh <mid> <judge-sid> <party-sid> <party-sid> ...
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

Requires the host environment to have the cURL binary in $PATH.

### Alpine mini-rootfs

```sh
# Move projects inside
mkdir -p alpine/var/botany
cp botany/judge/compile.sh botany/judge/match.sh alpine/var/botany
git clone https://github.com/ioi/isolate.git alpine/home/isolate

sudo mkdir alpine/proc/self
sudo mount --bind alpine alpine
sudo mount --bind /proc alpine/proc

# Set up network
sudo cp /etc/resolv.conf alpine/etc

sudo chroot alpine sh
# Inside chroot
sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
apk add gcc make libc-dev libcap-dev

cd /home/isolate
make
make install

# Test isolate
isolate --init
isolate --run -- /bin/echo hi
isolate --cleanup
```
