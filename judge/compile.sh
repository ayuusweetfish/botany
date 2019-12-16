#!/bin/sh

echo "Compiling $1"
SID=$1
CODE=submissions/$SID/code
BIN=submissions/$SID/bin

sleep 1

echo "#!/bin/sh" > $BIN
echo "echo ====== Submission $SID ======" >> $BIN
echo "cat <<EOF" >> $BIN
cat $CODE >> $BIN
echo "" >> $BIN
echo "EOF" >> $BIN
chmod +x $BIN
