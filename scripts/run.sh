#!/usr/bin/env bash

export GO111MODULE=on
export GOPROXY=https://goproxy.io
export GOROOT=/usr/lib/go
export GOPATH=/data/go


git pull origin master

cd ../
go install

supervisorctl restart antPush
