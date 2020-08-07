# batteryhook
Run commands on certain battery levels.
Written in go.

## Usage
```
Usage: batteryhook [-h|--help] [-v] [-i INTERVAL] [-b BATTTERY]
  -b string
    	Select battery to watch (default "BAT0")
  -i uint
    	Battery refresh interval in ms (default 5000)
  -v	Increase verbosity

Hooks are defined in the file $XDG_CONFIG_HOME/batteryhook/hooks.yaml.
This path defaults to ~/.config/batteryhook/hooks.yaml on most machines.

Example hooks.yaml:
- charging: false
  discharging: true
  level: 10
  command: systemctl hibernate

This will hibernate your machine when it falls below 10% while not being connected to external power.
```

## Build
Run `make build` to build batteryhook.

## Install
Run `make install` to build and install batteryhook to `/usr/local/bin`.

A custom target directory can be set via the `PREFIX` variable.
Run `make PREFIX=~/.local install` to install batteryhook to `~/.local/bin`.

## Examples for `hooks.yaml`
Notify about low battery and hibernate on very low battery (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify)):
```yaml
- charging: false
  discharging: true
  level: 25
  command: notify-send "Battery is running low"
- charging: false
  discharging: true
  level: 10
  command: systemctl hibernate
```

Run multiple commands per hook using POSIX shell syntax (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify) and [brightnessctl](https://github.com/Hummer12007/brightnessctl)):
```yaml
- charging: false
  discharging: true
  level: 20
  command: notify-send "Low battery" && brightnessctl s 10%
```
