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
/dev/urandom
/dev/zero
EOF
./addtojail /usr/local/bin/busybox
./addtojail /usr/lib/dyld
./addtojail /bin/sh
./addtojail /etc/resolv.conf
sudo ./mkjail 1

# Inside chroot shell
/usr/local/bin/busybox mkdir /usr/bin
/usr/local/bin/busybox --install -s /usr/bin

cd /usr/local
mkdir homebrew
# Not yet working...
wget -O - https://github.com/Homebrew/brew/tarball/master | tar xz --strip 1 -C homebrew
```
