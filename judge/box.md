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

### Alpine mini-rootfs

```sh
# Move projects inside
git clone https://github.com/kawa-yoiko/isolate.git alpine/home/isolate
git -C alpine/home/isolate checkout botany

# Create mount points
sudo mount -B alpine alpine
sudo mount -t proc none alpine/proc
# NOTE: In case where /dev is needed, use
# sudo mount -B /dev alpine/dev

# Set up network
sudo cp /etc/resolv.conf alpine/etc

# Enter chroot
sudo chroot alpine sh

# Install packages
sed -i 's/dl-cdn.alpinelinux.org/mirrors.tuna.tsinghua.edu.cn/g' /etc/apk/repositories
apk add gcc make libc-dev libcap-dev bash g++ lua python3

# Build isolate
cd /home/isolate
make isolate
make install

# Test isolate
isolate --init
isolate --run --dir=box=./box --dir=tmp= --dir=proc= -- /bin/echo hi
isolate --run --dir=box=./box --dir=tmp= --dir=proc= -- /usr/bin/lua -e "print('hi')"
isolate --run --dir=box=./box --dir=tmp= --dir=proc= -- /usr/bin/lua -e "io.open('1.txt', 'w'):write('hi')" # Should fail
isolate --cleanup
```

In the host environment, run Botany's judge side:

```sh
sudo apt-get install libhiredis-dev libb2-dev

cd botany/judge
sudo ./build.sh
./a.out -i 1 -d /path/to/alpine
```

#### Cleanup

Use `sudo` if necessary.

```sh
# umount -l alpine/proc alpine/dev alpine
umount -l alpine
rm -rf alpine
```
