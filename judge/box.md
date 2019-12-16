# Botany: Directory Structure for Judge

```
/var/botany
| - compile.sh: Call ./compile.sh <sid> to compile
| - match.sh: Call ./match.sh <mid> <sid> <sid> ... to run a match
| - submissions/
| | - 1/
| | | - lang
| | | - code
| | | - bin
| | - ...
```

### Steps to build

Set up a chroot jail and log in with `root`.

```sh
mkdir -p /var/botany
cd /var/botany
chown <outside_user> -R .`
cp /path/to/compile.sh .
cp /path/to/match.sh .
```
