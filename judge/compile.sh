#!/bin/sh

echo "Compiling $1"
SID=$1
LANG=$2
CODE=submissions/$SID/code.$LANG
BIN=submissions/$SID/bin

sleep 1

if [[ "$LANG" == "c" ]]; then
    gcc $CODE -O2 -o $BIN 2>&1
elif [[ "$LANG" == "lua" ]]; then
    echo "#!/bin/sh" > $BIN
    echo "echo \"====== Submission $SID ($LANG) ======\"" >> $BIN
    echo "cat <<EOF" >> $BIN
    cat $CODE >> $BIN
    echo "" >> $BIN
    echo "EOF" >> $BIN
    chmod +x $BIN
fi
