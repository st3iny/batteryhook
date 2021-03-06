package main

import (
    "flag"
    "fmt"
    "log"
    "os"
    "os/signal"
    "strconv"
    "strings"
    "syscall"
    "time"

    "github.com/st3iny/batteryhook/internal/battery"
    "github.com/st3iny/batteryhook/internal/conf"
    "github.com/st3iny/batteryhook/internal/util"
)

var interval uint
var batteryName string
var testInterval string

func main() {
    flag.BoolVar(&util.Verbose, "v", false, "Increase verbosity")
    flag.UintVar(&interval, "i", 5000, "Battery refresh interval in ms")
    flag.StringVar(&batteryName, "b", "BAT0", "Select battery to watch")
    flag.StringVar(&testInterval, "t", "", "Test hooks in level interval (format \"BEGIN[-END]\")")
    flag.Usage = usage
    flag.Parse()

    conf, err := conf.Load()
    util.Check(err)

    if util.Verbose {
        for _, h := range conf.Hooks {
            log.Println("Parsed hook", h)
        }
    }

    if len(conf.Hooks) == 0 {
        fmt.Fprintln(os.Stderr, "No hooks found (try --help)")
        os.Exit(1)
    }

    events := make(chan *battery.Event, 1)

    // Create real or fake testing battery
    var bat battery.Battery
    var pollInterval time.Duration
    if testInterval == "" {
        bat, err = battery.NewLinuxBattery(batteryName)
        util.Check(err)

        pollInterval = time.Duration(interval) * time.Millisecond
    } else {
        testInterval := strings.Split(testInterval, "-")

        testBegin, err := strconv.Atoi(testInterval[0])
        util.Check(err)

        testEnd := testBegin
        if len(testInterval) > 1 {
            testEnd, err = strconv.Atoi(testInterval[1])
            util.Check(err)
        }

        bat = battery.NewFakeBattery(testBegin, testEnd)
        pollInterval = 0
    }

    // Watch battery level and status
    go battery.Watch(bat, events, pollInterval)

    // Terminate on signal
    go func() {
        sigc := make(chan os.Signal, 1)
        signal.Notify(sigc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
        for {
            <-sigc
            os.Exit(0)
        }
    }()

    // Listen for events
    for event := range events {
        for _, h := range conf.Hooks {
            err := h.ProcessEvent(event)
            if util.Verbose && err != nil {
                log.Println("Error while handling battery event")
                log.Println(err)
            }
        }
    }
}

func usage() {
    fmt.Fprintln(os.Stderr, "Usage: batteryhook [-h|--help] [-v] [-i INTERVAL] [-b BATTTERY]")
    flag.PrintDefaults()
    help := []string{
        "",
        "Hooks are defined in the file $XDG_CONFIG_HOME/batteryhook/config.yaml.",
        "This path defaults to ~/.config/batteryhook/config.yaml on most machines.",
        "",
        "Example config.yaml:",
        "hooks:",
        "  - status:",
        "      discharging: true",
        "    level: 10",
        "    command: systemctl hibernate",
        "",
        "This will hibernate your machine when it falls below 10% while discharging.",
        "",
        "Valid statuses are: unknown, discharging, charging, not_charging and full",
        "Status defaults to (if omitted): discharging",
        "Refer to the linux documentation for more information on battery status (power_supply.h).",
    }
    fmt.Fprintln(os.Stderr, strings.Join(help, "\n"))
}
