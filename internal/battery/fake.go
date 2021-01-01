package battery

type FakeBattery struct {
    next int
    low int
}

func (bat *FakeBattery) Poll() (bool, error) {
    bat.next--
    return bat.next >= bat.low, nil
}

func (bat *FakeBattery) Level() int {
    return bat.next
}

func (bat *FakeBattery) Status() Status {
    return Any
}

func (bat *FakeBattery) String() string {
    return "FAKE"
}

func NewFakeBattery(high, low int) *FakeBattery {
    return &FakeBattery{next: high + 1, low: low}
}
