package util

import (
    "log"
    "path"

    "github.com/adrg/xdg"
)

var Verbose bool

// log and quit if an error occurred
func Check(err error) {
    if err != nil {
        log.Fatalln(err)
    }
}

// build path to given config file using the xdg standard
func BuildConfigPath(file string) (string, error) {
    path, err := xdg.ConfigFile(path.Join("batteryhook", file))
    if err != nil {
        return "", err
    }

    return path, nil
}

