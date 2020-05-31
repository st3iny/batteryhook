package battery

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strconv"
    "strings"
    "time"

    "github.com/st3iny/batteryhook/internal/util"
)

const (
    Both int = 1
    Charging int = 2
    Discharging int = 3
)

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
        if err == nil && level != bat.lastLevel {
            start := bat.lastLevel + 1
            if bat.lastLevel == -1 {
                start = level
            }

            step := 1
            if level < bat.lastLevel {
                step = -1
            }

            for i := start;; i += step {
                events <- &Event{battery: bat, level: i}
                if i == level {
                    break
                }
            }
            bat.lastLevel = level
        }
        time.Sleep(interval * time.Millisecond)
        level += 3
    }
}

func GetAll() []*Battery {
    const batteryDir string = "/sys/class/power_supply"
    files, err := ioutil.ReadDir(batteryDir)
    util.Check(err)

    batteries := []*Battery{}
    for _, file := range files {
        if strings.HasPrefix(file.Name(), "BAT") {
            if util.Verbose {
                log.Println("Found battery", file.Name())
            }

            bat := &Battery{
                name: file.Name(),
                lastLevel: -1,
                path: path.Join(batteryDir, file.Name()),
            }
            batteries = append(batteries, bat)
        }
    }
    return batteries
}
