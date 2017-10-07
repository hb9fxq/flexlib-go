#!/bin/bash
PATH="/usr/local/go/bin:$PATH"
export GOPATH=$(pwd):~/devLibs/gopath

rm -rf bin/*

export GOARCH=386
export GOOS=windows
go build --o ./bin/iq-transfer_win32.exe github.com/krippendorf/cmd/iq-transfer

export GOARCH=amd64
export GOOS=windows
go build --o ./bin/iq-transfer_win64.exe github.com/krippendorf/cmd/iq-transfer

export GOARCH=amd64
export GOOS=linux
go build --o ./bin/iq-transfer_linux64 github.com/krippendorf/cmd/iq-transfer

export GOARCH=386
export GOOS=linux
go build --o ./bin/iq-transfer_linux32 github.com/krippendorf/cmd/iq-transfer

export GOARCH=arm
export GOOS=linux
export GOARM=5
go build --o ./bin/iq-transfer_arm5_raspi github.com/krippendorf/cmd/iq-transfer

