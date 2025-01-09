#!/bin/bash
DIR=`cd $(dirname $0);pwd`
cd $DIR

## 拷贝协议到项目跟目录 
# ln -sf ../../goprotos ./goprotos

rm -f ../../src/protocol/ss/service_register.pb.go
rm -f ../../src/protocol/cs/client_register.pb.go

go${exe} mod tidy

go${exe} build -o register${exe} scripts/main.go

chmod +x register${exe}

