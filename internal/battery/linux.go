package battery

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "path"
    "strconv"

    "github.com/st3iny/batteryhook/internal/util"
)

const batteryBaseDir string = "/sys/class/power_supply"

type LinuxBattery struct {
    name string
    path string
    level int
    status Status
}

func (bat *LinuxBattery) Poll() (bool, error) {
    var err error

    bat.level, err = bat.getLevel()
    if err != nil {
        return true, err
    }

    bat.status, err = bat.getStatus()
    if err != nil {
        return true, err
    }

    return true, nil
}

func (bat *LinuxBattery) Level() int {
    return bat.level
}

func (bat *LinuxBattery) Status() Status {
    return bat.status
}

func (bat *LinuxBattery) String() string {
    return bat.name
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
        path: batteryDir,
    }
    return bat, nil
}

func (bat *LinuxBattery) getLevel() (int, error) {
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

    return int(capacity), nil
}

func (bat *LinuxBattery) getStatus() (Status, error) {
    file, err := os.Open(path.Join(bat.path, "status"))
    if err != nil {
        return 0, err
    }

    defer file.Close()

    statusRaw, err := ioutil.ReadAll(file)
    if err != nil {
        return 0, err
    }

    var status Status
    switch string(statusRaw[:len(statusRaw) - 1]) {
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
        return 0, fmt.Errorf("can't parse battery status")
    }

    return status, nil
}
