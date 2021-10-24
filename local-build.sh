#!/bin/bash
set -x
go build -o app appinstaller/appinstaller.go appinstaller/app_notifier.go

if [ $? -eq 0 ]; then
    cp app ~/github/clockworkpi/launchergo
fi


go build -o launchergo main.go mainscreen.go
#go build -ldflags="-s -w" -o main main.go mainscreen.go

if [ $? -eq 0 ]; then
    cp launchergo ~/github/clockworkpi/launchergo
fi




