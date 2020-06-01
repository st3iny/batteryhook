package main

import (
    "flag"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/st3iny/batteryhook/internal/battery"
    "github.com/st3iny/batteryhook/internal/util"
)

func main() {
    var interval uint
    var batteryName string
    flag.BoolVar(&util.Verbose, "v", false, "Increase verbosity")
    flag.UintVar(&interval, "i", 5000, "Battery refresh interval in ms")
    flag.StringVar(&batteryName, "b", "BAT0", "Select battery to watch")
    flag.Usage = func() {
        fmt.Fprintln(os.Stderr, "Usage: batteryhook [-h|--help] [-v] [-i INTERVAL] [-b BATTTERY] HOOK [HOOK ...]")
        flag.PrintDefaults()
        fmt.Fprintln(os.Stderr, "\nHooks have the format STATUS,LEVEL,CMD")
        fmt.Fprintln(os.Stderr, "STATUS is one of <c|d|cd> (charging, discharging, both)")
        fmt.Fprintln(os.Stderr, "LEVEL  is an int between 0 and 100")
        fmt.Fprintln(os.Stderr, "CMD    command to be executed (through sh -c) when triggered")
    }
    flag.Parse()

    bat, err := battery.Get(batteryName)
    util.Check(err)

    hooks, err := battery.ParseHooks(flag.Args())
    util.Check(err)

    if len(hooks) == 0 {
        flag.Usage()
        os.Exit(1)
    }

    events := make(chan *battery.Event, 1)
    go bat.Watch(events, time.Duration(interval) * time.Millisecond)

    for _, h := range hooks {
        go h.ProcessEvents(events)
    }

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    for {
        <-sigc
        break
    }
}
