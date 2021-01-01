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
            events <- &Event{Battery: bat, Level: level}
        } else if err == nil && level != lastLevel {
            step := 1
            if level < lastLevel {
                step = -1
            }

            start := lastLevel + step
            for i := start;; i += step {
                events <- &Event{Battery: bat, Level: i}
                if i == level {
                    break
                }
            }
        }

        lastLevel = level
        time.Sleep(interval)
    }
}
