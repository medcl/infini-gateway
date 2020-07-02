#!/bin/sh

cd /go/src/infini.sh/

echo "INFINI Elasticsearch PROXY READY TO ROCK!"

cd proxy
make build

cd /go/src/infini.sh/proxy && ./bin/proxy