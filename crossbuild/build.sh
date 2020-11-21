#!/bin/bash

cd ..
#Mac
env GOOS=darwin GOARCH=amd64 go build -o ../../../../bin/flexlib-go/osx/smartsdr-daxclient github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
env GOOS=darwin GOARCH=amd64 go build -o ../../../../bin/flexlib-go/osx/smartsdr-iqtransfer github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
env GOOS=darwin GOARCH=amd64 go build -o ../../../../bin/flexlib-go/osx/smartsdr-mqttadapter github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter

# Linux
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flexlib-go/linux64/smartsdr-daxclient github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flexlib-go/linux64/smartsdr-iqtransfer github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flexlib-go/linux64/smartsdr-mqttadapter github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter

# Raspi
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flexlib-go/raspberryPi/smartsdr-daxclient github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flexlib-go/raspberryPi/smartsdr-iqtransfer github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flexlib-go/raspberryPi/smartsdr-mqttadapter github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter

# Windows
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flexlib-go/Win64/smartsdr-daxclient.exe github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flexlib-go/Win64/smartsdr-iqtransfer.exe github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flexlib-go/Win64/smartsdr-mqttadapter github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter

env GOOS=windows GOARCH=386 go build -o ../../../../bin/flexlib-go/Win32/smartsdr-daxclient github.com/hb9fxq/flexlib-go/cmd/smartsdr-daxclient
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flexlib-go/Win32/smartsdr-iqtransfer.exe github.com/hb9fxq/flexlib-go/cmd/smartsdr-iqtransfer
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flexlib-go/Win32/smartsdr-mqttadapter.exe github.com/hb9fxq/flexlib-go/cmd/smartsdr-mqttadapter