#!/bin/bash

DIR=`cd $(dirname $0);pwd`
cd $DIR

if [[ $# -eq 0 ]]; then
    BuildTags="internal_test"
else
    BuildTags=$1
fi

echo "building gamesvr, tags: ${BuildTags}"
PROJECT_NAME=`basename ${DIR}`
PROJECT_BIN=${PROJECT_NAME}${exe}
PROJECT_HOME=${DIR}/../../run/${PROJECT_NAME}/bin

go${exe} mod tidy

#go${exe} build -race -gcflags="all=-N -l" -tags="${BuildTags}" -o ${PROJECT_BIN} scripts/main.go
go${exe} build -gcflags="all=-N -l" -tags="${BuildTags}" -o ${PROJECT_BIN} main.go

let retCode=$?

if [ -e ${PROJECT_BIN} ]
then
    if [ "x${exe}" == "x" ] 
    then
        chmod +x ${PROJECT_BIN};
    fi

    mv ${PROJECT_BIN} ${PROJECT_HOME};
fi

# go build -o gamesvr main.go 