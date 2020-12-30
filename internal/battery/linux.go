package battery

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strconv"
    "time"

    "github.com/st3iny/batteryhook/internal/util"
)

const batteryBaseDir string = "/sys/class/power_supply"

type LinuxBattery struct {
    name string
    lastLevel int
    path string
}

func (bat *LinuxBattery) Level() (int, error) {
    file, err := os.Open(path.Join(bat.path, "capacity"))
    if err != nil {
        return 0, err
    }

    defer func() {
        util.Check(file.Close())
    }()

    capacityRaw, err := ioutil.ReadAll(file)
    if err != nil {
        return 0, err
    }

    capacity, err := strconv.ParseInt(string(capacityRaw[:len(capacityRaw) - 1]), 10, 16)
    if err != nil {
        return 0, err
    }

    level := int(capacity)
    return int(level), nil
}

func (bat *LinuxBattery) Status() (Status, error) {
    file, err := os.Open(path.Join(bat.path, "status"))
    if err != nil {
        return 0, err
    }

    defer func() {
        util.Check(file.Close())
    }()

    statusRaw, err := ioutil.ReadAll(file)
    if err != nil {
        return 0, err
    }

    return parseStatus(string(statusRaw[:len(statusRaw) - 1]))
}

func (bat *LinuxBattery) String() string {
    return bat.name
}

func (bat *LinuxBattery) Watch(events chan *Event, interval time.Duration) {
    for {
        level, err := bat.Level()
        if err == nil && bat.lastLevel == -1 {
            events <- &Event{Battery: bat, Level: level}
        } else if err == nil && level != bat.lastLevel {
            step := 1
            if level < bat.lastLevel {
                step = -1
            }

            start := bat.lastLevel + step
            for i := start;; i += step {
                events <- &Event{Battery: bat, Level: i}
                if i == level {
                    break
                }
            }
        }

        bat.lastLevel = level
        time.Sleep(interval)
    }
}

func NewLinuxBattery(name string) (*LinuxBattery, error) {
    batteryDir := path.Join(batteryBaseDir, name)
    if s, err := os.Stat(batteryDir); os.IsNotExist(err) || !s.IsDir() {
        return nil, fmt.Errorf("No such battery %s", name)
    }

    if util.Verbose {
        log.Println("Found battery", name, "at", batteryDir)
    }

    bat := &LinuxBattery{
        name: name,
        lastLevel: -1,
        path: batteryDir,
    }
    return bat, nil
}

func parseStatus(rawStatus string) (Status, error) {
    var status Status
    switch rawStatus {
    case "Unknown":
        status = Unknown
    case "Discharging":
        status = Discharging
    case "Charging":
        status = Charging
    case "Not charging":
        status = NotCharging
    case "Full":
        status = Full
    default:
        return 0, fmt.Errorf("unknown battery status")
    }

    return status, nil
}
