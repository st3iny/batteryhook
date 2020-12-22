package battery

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "os/exec"

    "github.com/st3iny/batteryhook/internal/util"

    "gopkg.in/yaml.v2"
)

type Hook struct {
    Charging bool `yaml:"charging"`
    Discharging bool `yaml:"discharging"`
    Level int `yaml:"level"`
    Command string `yaml:"command"`
}

func LoadHooks() ([]Hook, error) {
    hookPath, err := util.BuildConfigPath("hooks.yaml")
    if err != nil {
        return nil, err
    }

    _, err = os.Stat(hookPath)
    if os.IsNotExist(err) {
        return nil, nil
    } else if err != nil {
        return nil, err
    }

    hooksBlob, err := ioutil.ReadFile(hookPath)
    if err != nil {
        return nil, err
    }

    var hooks []Hook
    err = yaml.Unmarshal(hooksBlob, &hooks)
    if err != nil {
        return nil, err
    }

    return hooks, nil
}

func (h Hook) String() string {
    return fmt.Sprintf(
        "{charging: %t, discharging: %t, level: %d, command: \"%s\"}",
        h.Charging, h.Discharging, h.Level, h.Command,
    )
}

func (h *Hook) ProcessEvent(event *Event) error {
    status, err := event.battery.Status()
    if err != nil {
        return err
    }

    trigger := false
    if status == Charging && h.Charging {
        trigger = true
    } else if status == Discharging && h.Discharging {
        trigger = true
    } else if status == Both && (h.Charging || h.Discharging) {
        trigger = true
    }

    if trigger && h.Level == event.level {
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
