package battery

import (
    "log"
    "time"

    "github.com/st3iny/batteryhook/internal/util"
)

type Status int

const (
    Unknown     Status = 0
    Discharging Status = 1
    Charging    Status = 2
    NotCharging Status = 3
    Full        Status = 4
    Any         Status = 5
)

func (status Status) String() string {
    statusString := ""
    switch status {
    case Unknown:
        statusString = "Unknown"
    case Discharging:
        statusString = "Discharging"
    case Charging:
        statusString = "Charging"
    case NotCharging:
        statusString = "NotCharging"
    case Full:
        statusString = "Full"
    case Any:
        statusString = "Any"
    }
    return statusString
}

type Battery interface {
    Poll() (bool, error)
    Level() int
    Status() Status
    String() string
}

func Watch(bat Battery, events chan *Event, interval time.Duration) {
    defer close(events)

    lastLevel := -1
    for {
        hasNext, err := bat.Poll()
        if err != nil && util.Verbose {
            log.Println("Error while polling battery", bat)
            log.Println(err)
        }

        if err != nil {
            continue
        }

        if !hasNext {
            return
        }

        level := bat.Level()
        if err == nil && lastLevel == -1 {
            events <- &Event{Level: level, Status: bat.Status()}
        } else if err == nil && level != lastLevel {
            step := 1
            if level < lastLevel {
                step = -1
            }

            start := lastLevel + step
            for i := start;; i += step {
                events <- &Event{Level: i, Status: bat.Status()}
                if i == level {
                    break
                }
            }
        }

        lastLevel = level
        time.Sleep(interval)
    }
}
