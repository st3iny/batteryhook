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

const (
    Both int = 1
    Charging int = 2
    Discharging int = 3
)

const batteryBaseDir string = "/sys/class/power_supply"

type Battery struct {
    name string
    lastLevel int
    path string
}

func (bat *Battery) Level() (int, error) {
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

func (bat *Battery) Status() (int, error) {
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

    var status int
    switch string(statusRaw[:len(statusRaw) - 1]) {
    case "Charging":
        status = Charging
    case "Discharging":
        status = Discharging
    default:
        return 0, fmt.Errorf("unknown battery status")
    }

    return status, nil
}

func (bat *Battery) Watch(events chan *Event, interval time.Duration) {
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

func Get(name string) (*Battery, error) {
    batteryDir := path.Join(batteryBaseDir, name)
    if s, err := os.Stat(batteryDir); os.IsNotExist(err) || !s.IsDir() {
        return nil, fmt.Errorf("No such battery %s", name)
    }

    if util.Verbose {
        log.Println("Found battery", name, "at", batteryDir)
    }

    bat := &Battery{
        name: name,
        lastLevel: -1,
        path: batteryDir,
    }
    return bat, nil
}
