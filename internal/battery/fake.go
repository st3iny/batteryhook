package battery

type FakeBattery struct {
    level int
}

func (bat *FakeBattery) Level() (int, error) {
    return bat.level, nil
}

func (bat *FakeBattery) Status() (Status, error) {
    return Any, nil
}

func (bat *FakeBattery) String() string {
    return "FAKE"
}

func NewFakeBattery(level int) *FakeBattery {
    return &FakeBattery{level: level}
}
