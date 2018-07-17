package require



type Assertions struct {
	t TestingT
}


func New(t TestingT) *Assertions {
	return &Assertions{
		t: t,
	}
}


