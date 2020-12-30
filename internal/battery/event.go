package battery

import "fmt"

type Event struct {
    Battery Battery
    Level int
}

func (e *Event) String() string {
    return fmt.Sprintf("{battery: %s, level: %d}", e.Battery, e.Level)
}
