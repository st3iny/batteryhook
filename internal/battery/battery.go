package battery

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
    Level() (int, error)
    Status() (Status, error)
    String() string
}
