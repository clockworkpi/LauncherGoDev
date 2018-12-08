#!/bin/bash
set -x
DST=/home/cpi/launchergo
cp -f main $DST

cp -f Menu/GameShell/HelloWorld/HelloWorld.so $DST/Menu/GameShell/HelloWorld/

cp -f Menu/GameShell/10_Settings/Settings.so $DST/Menu/GameShell/10_Settings/

cp -f Menu/GameShell/10_Settings/About/about.so $DST/Menu/GameShell/10_Settings/About/

cp -f Menu/GameShell/10_Settings/Wifi/wifi.so $DST/Menu/GameShell/10_Settings/Wifi/

cp -f Menu/GameShell/10_Settings/Sound/sound.so $DST/Menu/GameShell/10_Settings/Sound/

cp -f Menu/GameShell/10_Settings/Brightness/brightness.so $DST/Menu/GameShell/10_Settings/Brightness/

cp -f Menu/GameShell/10_Settings/Update/update.so $DST/Menu/GameShell/10_Settings/Update/
