#!/bin/bash
set -x
go build -o launchergo main.go mainscreen.go
#go build -ldflags="-s -w" -o main main.go mainscreen.go




