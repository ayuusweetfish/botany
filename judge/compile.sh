#!/bin/sh

echo "Compiling $1"
SID=$1
LANG=$2
CODE=submissions/$SID/code.$LANG
BIN=submissions/$SID/bin

if [[ "$LANG" == "c" ]]; then
    gcc $CODE -O2 -I/var/botany/lib /var/botany/lib/ipc.c -o $BIN 2>&1
elif [[ "$LANG" == "cpp" ]]; then
    g++ $CODE -O2 -I/var/botany/lib /var/botany/lib/ipc.c -o $BIN 2>&1
elif [[ "$LANG" == "lua" ]]; then
    echo "#!/bin/sh" > $BIN
    echo "echo \"====== Submission $SID ($LANG) ======\"" >> $BIN
    echo "echo \"This is log from $SID\" >&2" >> $BIN
    echo "cat <<EOF" >> $BIN
    cat $CODE >> $BIN
    echo "" >> $BIN
    echo "EOF" >> $BIN
    chmod +x $BIN
fi
