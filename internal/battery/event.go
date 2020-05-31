package battery

import "fmt"

type Event struct {
    battery *Battery
    level int
}

func (e *Event) String() string {
    return fmt.Sprintf("{battery: %s, level: %d}", e.battery.name, e.level)
}
