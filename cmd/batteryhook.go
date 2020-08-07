package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strings"
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
    flag.Usage = usage
    flag.Parse()

    bat, err := battery.Get(batteryName)
    util.Check(err)

    hooks, err := battery.LoadHooks()
    util.Check(err)

    if util.Verbose {
        for _, h := range hooks {
            log.Println("Parsed hook", h)
        }
    }

    if len(hooks) == 0 {
        fmt.Fprintln(os.Stderr, "No hooks found (try --help)")
        os.Exit(1)
    }

    events := make(chan *battery.Event, 1)
    go bat.Watch(events, time.Duration(interval) * time.Millisecond)
    go func() {
        for {
            event := <-events
            for _, h := range hooks {
                err := h.ProcessEvent(event)
                if util.Verbose && err != nil {
                    log.Println("Error while handling battery event")
                    log.Println(err)
                }
            }
        }
    }()

    sigc := make(chan os.Signal, 1)
    signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
    for {
        <-sigc
        break
    }
}

func usage() {
    fmt.Fprintln(os.Stderr, "Usage: batteryhook [-h|--help] [-v] [-i INTERVAL] [-b BATTTERY]")
    flag.PrintDefaults()
    help := []string{
        "",
        "Hooks are defined in the file $XDG_CONFIG_HOME/batteryhook/hooks.yaml.",
        "This path defaults to ~/.config/batteryhook/hooks.yaml on most machines.",
        "",
        "Example hooks.yaml:",
        "- charging: false",
        "  discharging: true",
        "  level: 10",
        "  command: systemctl hibernate",
        "",
        "This will hibernate your machine when it falls below 10% while not being connected to external power.",
    }
    fmt.Fprintln(os.Stderr, strings.Join(help, "\n"))
}
