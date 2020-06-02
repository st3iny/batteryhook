# batteryhook
Run commands on certain battery levels.
Written in go.

## Usage
```
Usage: batteryhook [-h|--help] [-v] [-i INTERVAL] [-b BATTTERY] HOOK [HOOK ...]
  -b string
        Select battery to watch (default "BAT0")
  -i uint
        Battery refresh interval in ms (default 5000)
  -v    Increase verbosity

Hooks have the format STATUS,LEVEL,CMD
STATUS is one of <c|d|cd> (charging, discharging, both)
LEVEL  is an int between 0 and 100
CMD    command to be executed (through sh -c) when triggered
```

## Build
Run `make build` to build batteryhook.

## Install
Run `make install` to build and install batteryhook to `/usr/local/bin`.

A custom target directory can be set via the `PREFIX` environment variable.
Run `PREFIX=~/.local make install` to install batteryhook to `~/.local/bin`.

## Examples
Notify about low battery and hibernate on very low battery (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify)):
```sh
batteryhook 'd,25,notify-send "Battery is running low"' 'd,10,systemctl hibernate'
```

Run multiple commands per hook using POSIX shell syntax (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify) and [brightnessctl](https://github.com/Hummer12007/brightnessctl)):
```sh
batteryhook 'd,20,notify-send "Low battery" && brightnessctl s 10%'
```
