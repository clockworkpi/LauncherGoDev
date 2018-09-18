#!/bin/bash

go build main.go

##HelloWorld
cd Menu/GameShell/HelloWorld/
go build -o HelloWorld.so -buildmode=plugin
cd -

