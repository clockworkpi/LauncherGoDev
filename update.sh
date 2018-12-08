#!/bin/bash

sleep 10
feh --bg-center /home/cpi/launchergo/sysgo/gameshell/wallpaper/updating.png
cd /home/cpi/launchergo 
git pull
git reset --hard $1
feh --bg-center /home/cpi/launchergo/sysgo/gameshell/wallpaper/loading.png 
./load.sh

