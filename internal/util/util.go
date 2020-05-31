package util

import "log"

var Verbose bool

func Check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
