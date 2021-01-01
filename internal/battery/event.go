package battery

import "fmt"

type Event struct {
    Level int
    Status Status
}

func (e *Event) String() string {
    return fmt.Sprintf("{Level: %d, Status: %s}", e.Level, e.Status)
}
