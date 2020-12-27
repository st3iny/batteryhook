package conf

import (
    "fmt"
    "log"
    "os"
    "os/exec"

    "github.com/st3iny/batteryhook/internal/battery"
    "github.com/st3iny/batteryhook/internal/util"
)

type Hook struct {
    Charging bool `yaml:"charging"`
    Discharging bool `yaml:"discharging"`
    Level int `yaml:"level"`
    Command string `yaml:"command"`
}

func (h Hook) String() string {
    return fmt.Sprintf(
        "{charging: %t, discharging: %t, level: %d, command: \"%s\"}",
        h.Charging, h.Discharging, h.Level, h.Command,
    )
}

func (h *Hook) ProcessEvent(event *battery.Event) error {
    status, err := event.Battery.Status()
    if err != nil {
        return err
    }

    trigger := false
    if status == battery.Charging && h.Charging {
        trigger = true
    } else if status == battery.Discharging && h.Discharging {
        trigger = true
    } else if status == battery.Both && (h.Charging || h.Discharging) {
        trigger = true
    }

    if trigger && h.Level == event.Level {
        if util.Verbose {
            log.Println("Trigger battery event", event)
        }

        go func() {
            cmd := exec.Command("/bin/sh", "-c", h.Command)
            cmd.Stdout = os.Stdout
            cmd.Stderr = os.Stderr
            cmd.Run()
        }()
    }

    return nil
}
