## Building BusyBox under macOS

```sh
git apply busybox-macos.diff
LC_ALL=C make allnoconfig
LC_ALL=C make HOSTLDFLAGS=-lcurses menuconfig
LC_ALL=C make CROSS_COMPILE=llvm- SKIP_STRIP=y
```

## Building chroot filesystem under macOS

```sh
git apply mkjail.diff
cat > mkjail.files <<EOF
/dev/null
/dev/random
/dev/urandom
/dev/zero
EOF
./addtojail /usr/local/bin/busybox
./addtojail /usr/lib/dyld
sudo ./mkjail 1 # Will fail to enter shell
sudo mkdir jail/1/bin
sudo mv jail/1/usr/local/bin/busybox jail/1/bin/busybox
sudo chroot jail/1 busybox sh

# Inside chroot shell
/bin/busybox --install -s /bin
```

## With networking :construction:

```sh
git apply mkjail.diff
cat > mkjail.files <<EOF
/dev/null
/dev/random
/dev/urandom
/dev/zero
EOF
./addtojail /usr/local/bin/busybox
./addtojail /usr/lib/dyld
./addtojail /bin/sh
./addtojail /bin/bash
./addtojail /etc/resolv.conf
./addtojail /etc/hosts
./addtojail /opt/local/bin/curl
./addtojail /usr/bin/sw_vers
./addtojail /usr/bin/ruby
sudo ./mkjail 1

# Inside chroot shell
/usr/local/bin/busybox mkdir /usr/bin
/usr/local/bin/busybox --install -s /usr/bin

cd /usr/local
mkdir homebrew
curl -k -L https://github.com/Homebrew/brew/tarball/master | tar xz --strip 1 -C homebrew

# TODO: Install Ruby
curl -k https://cache.ruby-lang.org/pub/ruby/2.6/ruby-2.6.5.tar.gz

export PATH=$PATH:/usr/local/homebrew/bin
```
