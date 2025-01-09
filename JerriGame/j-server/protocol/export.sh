#!/bin/bash

TargetPath=../src/protocol

if [ -e $TargetPath ]; then
    echo "Error: Directory '$TargetPath' does not exist."
    rm -rf $TargetPath
else
    echo "Create directory: $TargetPath"
    mkdir -p $TargetPath
fi

EXPORT_PATH=(./cs ./ss)
for path in ${EXPORT_PATH[@]}
do
    if [ ! -d $path ]; then
        echo "Error: Directory '$path' does not exist."
        exit 1
    else
        echo "Export for .proto files in directory: $path"
        cd $path
        bash export.sh
        cd ..
    fi
done

cd ../src/protocol

go mod init protocol
go mod tidy

cd ../../tools/register

rm -f ../../src/protocol/ss/service_register.pb.go
rm -f ../../src/protocol/cs/client_register.pb.go

go mod init protocol
go mod tidy

go build -o register${exe} scripts/main.go
chmod +x register${exe}
./register${exe}
