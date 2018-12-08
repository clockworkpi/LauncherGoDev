#!/bin/bash

go build main.go mainscreen.go

##HelloWorld
cd Menu/GameShell/HelloWorld/
go build -ldflags="-s -w" -o HelloWorld.so -buildmode=plugin
cd -



