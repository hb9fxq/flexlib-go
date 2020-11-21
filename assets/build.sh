#!/bin/bash




cd ..

# Linux
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flexlib-go/linux64/iq-transfer github.com/hb9fxq/flexlib-go/cmd/iq-transfer
env GOOS=linux GOARCH=386 go build -o ../../../../bin/flexlib-go/linux32/iq-transfer github.com/hb9fxq/flexlib-go/cmd/iq-transfer
env GOOS=linux GOARCH=amd64 go build -o ../../../../bin/flexlib-go/linux64/MqttAdapter github.com/hb9fxq/flexlib-go/cmd/MqttAdapter
env GOOS=linux GOARCH=386 go build -o ../../../../bin/flexlib-go/linux32/MqttAdapter github.com/hb9fxq/flexlib-go/cmd/MqttAdapter

# Raspi
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flexlib-go/raspberryPi/iq-transfer github.com/hb9fxq/flexlib-go/cmd/iq-transfer
env GOOS=linux GOARCH=arm GOARM=5 go build -o ../../../../bin/flexlib-go/raspberryPi/MqttAdapter github.com/hb9fxq/flexlib-go/cmd/MqttAdapter

# Windows
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flexlib-go/Win64/MqttAdapter.exe github.com/hb9fxq/flexlib-go/cmd/iq-transfer
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flexlib-go/Win32/MqttAdapter.exe github.com/hb9fxq/flexlib-go/cmd/iq-transfer
env GOOS=windows GOARCH=amd64 go build -o ../../../../bin/flexlib-go/Win64/MqttAdapter.exe github.com/hb9fxq/flexlib-go/cmd/MqttAdapter
env GOOS=windows GOARCH=386 go build -o ../../../../bin/flexlib-go/Win32/MqttAdapter.exe github.com/hb9fxq/flexlib-go/cmd/MqttAdapter


# pfsense
#env GOOS=freebsd GOARCH=amd64 go build -o ../../../../bin/flexlib-go/pfSense64/iqtransfer github.com/hb9fxq/flexlib-go/cmd/iqtransfer
#env GOOS=freebsd GOARCH=386 go build -o ../../../../bin/flexlib-go/pfSense32/iqtransfer github.com/hb9fxq/flexlib-go/cmd/iqtransfer