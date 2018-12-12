# GameShell launcher
This is the launcher for GameShell based on 320x240 resolution and D-Pad layout.
![Screenshot](https://github.com/clockworkpi/GameShellDocs/blob/master/screenshot.png)

# Create the necessary user and group
* User name: cpi
* Password: cpi
* Group ID: 31415 with group name: cpifav

```
sudo adduser cpi  
sudo groupadd cpifav -g 31415  
sudo adduser cpi cpifav  
```

# Directory structure
```
/home/cpi/
├── launchergo <- Here we are
├── apps
│   ├── emulators
│   └── Menu
├── games
│   ├── FreeDM
│   ├── MAME
│   └── nxengine
└── music
```
## Create the necessary directories
```
mkdir -p /home/cpi/apps/emulators  
mkdir -p /home/cpi/games  
mkdir -p /home/cpi/music  
```

## Install dependent packages
```
go get -u -v github.com/cuu/LauncherGoDev
```

## Relevant opensource packages 
* https://github.com/godbus/dbus
* https://github.com/veandco/go-sdl2
* https://github.com/itchyny/volume-go
* https://github.com/go-ini/ini
* github.com/cavaliercoder/grab
* github.com/yookoala/realpath

# Create “.mpd_cpi.conf” config

vim ~/.mpd_cpi.conf

```
music_directory    "/home/cpi/music"
playlist_directory    "/home/cpi/music/playlists"
db_file    "/home/cpi/music/tag_cache"
log_file    "/tmp/mpd.log"
pid_file    "/tmp/mpd.pid"
state_file    "/home/cpi/music/mpd_state"
sticker_file    "/home/cpi/music/sticker.sql"
user    "cpi"
bind_to_address    "/tmp/mpd.socket"
auto_update    "yes"
auto_update_depth    "3" 
input {
    plugin "curl"
}

audio_output {
    type    "alsa"
    name    "My ALSA Device"
}

audio_output {
    type    "fifo"
    name    "my_fifo"
    path    "/tmp/mpd.fifo"
    format    "44100:16:2"
}

filesystem_charset    "UTF-8"
```
