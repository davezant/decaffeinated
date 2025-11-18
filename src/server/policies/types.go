package policies

type Policy struct {
	ID      string
	Message string
	Active  bool
}

type Log struct {
	Policy
	Level string
}

type Warning struct {
	Policy
	Severity int
}

type Killer struct {
	Policy
	TargetProcess string
}

type Restriction struct {
	Policy
	BlockNetwork bool
}

type Action struct {
	Percentage float32
	Callback   func()
}
