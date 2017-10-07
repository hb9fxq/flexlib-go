#!/bin/bash
PATH="/usr/local/go/bin:$PATH"
export GOPATH=$(pwd):~/devLibs/gopath

rm -rf bin/*

export GOARCH=386
export GOOS=windows
go install github.com/krippendorf/cmd/iq-transfer

export GOARCH=amd64
export GOOS=windows
go install github.com/krippendorf/cmd/iq-transfer

export GOARCH=amd64
export GOOS=linux
go install github.com/krippendorf/cmd/iq-transfer

export GOARCH=amd64
export GOOS=freebsd
go install github.com/krippendorf/cmd/iq-transfer

export GOARCH=386
export GOOS=linux
go install github.com/krippendorf/cmd/iq-transfer

export GOARCH=arm
export GOOS=linux
export GOARM=5
go install github.com/krippendorf/cmd/iq-transfer

