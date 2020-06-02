package battery

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strconv"
    "strings"

    "github.com/st3iny/batteryhook/internal/util"
)

type Hook struct {
    status int
    level int
    cmd string
}

func (h* Hook) String() string {
    var status string
    switch h.status {
    case Both:
        status = "Both"
    case Charging:
        status = "Charging"
    case Discharging:
        status = "Discharging"
    }
    return fmt.Sprintf("{status: %s, level: %d, cmd: \"%s\"}", status, h.level, h.cmd)
}

func (h *Hook) ProcessEvent(event *Event) {
    status, err := event.battery.Status()
    if err == nil && h.level == event.level && (h.status == Both || h.status == status) {
        if util.Verbose {
            log.Println("Trigger battery event", event)
        }

        cmd := exec.Command("/bin/sh", "-c", h.cmd)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        cmd.Run()
    }
}

func ParseHooks(args []string) ([]*Hook, error) {
    hooks := make([]*Hook, 0, len(args))

    for _, arg := range args {
        h, err := parseHook(arg)
        if err != nil {
            return nil, err
        }

        if util.Verbose {
            log.Println("Parsed hook", h)
        }
        hooks = append(hooks, h)
    }

    return hooks, nil
}

func parseHook(raw string) (*Hook, error) {
    parts := strings.Split(raw, ",")
    if len(parts) < 3 {
        return nil, fmt.Errorf("Too few parts in hook %s", raw)
    }
    statusStr := parts[0]
    levelStr := parts[1]
    command := strings.Join(parts[2:], ",")

    var status int
    switch (statusStr) {
    case "cd":
        status = Both
    case "c":
        status = Charging
    case "d":
        status = Discharging
    default:
        return nil, fmt.Errorf("Invalid status in hook %s", raw)
    }

    levelRaw, err := strconv.ParseInt(levelStr, 10, 64)
    if err != nil {
        return nil, err
    }

    level := int(levelRaw)
    if level > 100 || level < 0 {
        return nil, fmt.Errorf("Invalid level in hook %s", raw)
    }

    h := &Hook{
        status: status,
        level: level,
        cmd: command,
    }
    return h, nil
}
