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
    Status HookStatus `yaml:"status,omitempty"`
    Level int `yaml:"level"`
    Command string `yaml:"command"`
}

type HookStatus struct {
    Unknown bool `yaml:"unknown,omitempty"`
    Charging bool `yaml:"Charging,omitempty"`
    Discharging bool `yaml:"discharging,omitempty"`
    NotCharging bool `yaml:"not_charging,omitempty"`
    Full bool `yaml:"full,omitempty"`
}

func (h Hook) String() string {
    return fmt.Sprintf(
        "{Status: {%s: %t, %s: %t, %s: %t, %s: %t, %s: %t}, %s: %d, %s: \"%s\"}",
        "Unknown", h.Status.Unknown,
        "Discharging", h.Status.Discharging || h.Status.IsZero(),
        "Charging", h.Status.Charging,
        "NotCharging", h.Status.NotCharging,
        "Full", h.Status.Full,
        "Level", h.Level,
        "Command", h.Command,
    )
}

func (status HookStatus) IsZero() bool {
    return status == (HookStatus{})
}

func (h *Hook) ProcessEvent(event *battery.Event) error {
    trigger := false
    if event.Status == battery.Any {
        trigger = true
    } else if event.Status == battery.Unknown && h.Status.Unknown {
        trigger = true
    } else if event.Status == battery.Discharging && (h.Status.Discharging || h.Status.IsZero()) {
        trigger = true
    } else if event.Status == battery.Charging && h.Status.Charging {
        trigger = true
    } else if event.Status == battery.NotCharging && h.Status.NotCharging {
        trigger = true
    } else if event.Status == battery.Full && h.Status.Full {
        trigger = true
    }

    if trigger && h.Level == event.Level {
        if util.Verbose {
            log.Println("Trigger battery event", event)
        }

        cmd := exec.Command("/bin/sh", "-c", h.Command)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        if util.Verbose {
            log.Println("Running", cmd.Args)
        }

        err := cmd.Run()
        if err != nil {
            return err
        }
    }

    return nil
}
