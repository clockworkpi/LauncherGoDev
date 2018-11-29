#!/bin/bash

go build main.go

##HelloWorld
cd Menu/GameShell/HelloWorld/
go build -ldflags="-s -w" -o HelloWorld.so -buildmode=plugin
cd -


cd Menu/GameShell/10_Settings
go build -o  Settings.so -buildmode=plugin
cd -

cd Menu/GameShell/10_Settings/About
go build -o  about.so -buildmode=plugin
cd -

cd Menu/GameShell/10_Settings/Wifi
go build -o  wifi.so -buildmode=plugin
cd -

