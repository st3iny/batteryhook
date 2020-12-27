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
  -v    Increase verbosity

Hooks are defined in the file $XDG_CONFIG_HOME/batteryhook/config.yaml.
This path defaults to ~/.config/batteryhook/config.yaml on most machines.

Example config.yaml:
hooks:
  - status:
      discharging: true
    level: 10
    command: systemctl hibernate

This will hibernate your machine when it falls below 10% while discharging.

Valid statuses are: unknown, discharging, charging, not_charging and full
Status defaults to (if omitted): discharging
Refer to the linux documentation for more information on battery status (power_supply.h).
```

## Build
Run `make build` to build batteryhook.

## Install
Run `make install` to build and install batteryhook to `/usr/local/bin`.

A custom target directory can be set via the `PREFIX` variable.
Run `make PREFIX=~/.local install` to install batteryhook to `~/.local/bin`.

## Examples for `config.yaml`
Notify about low battery and hibernate on very low battery (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify)):
```yaml
hooks:
  - status:
      discharging: true
    level: 25
    command: notify-send "Battery is running low"
  - status:
      discharging: true
    discharging: true
    level: 10
    command: systemctl hibernate
```

Run multiple commands per hook using POSIX shell syntax (requires [libnotify](https://gitlab.gnome.org/GNOME/libnotify) and [brightnessctl](https://github.com/Hummer12007/brightnessctl)):
```yaml
hooks:
  - status:
      discharging: true
    level: 20
    command: notify-send "Low battery" && brightnessctl s 10%
```

Omitting a status will cause it to default to discharging.
