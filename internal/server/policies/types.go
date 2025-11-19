package policies

var DefaultPercentagesPolicies = []float32{0.5, 0.9, 1.0}

type Policy struct {
	ID         string
	Message    string
	Active     bool
	Percentage []float32
	Actions    []Action
}

type Action struct {
	Callback func()
}

func NewPolicy(name string, percentages []float32) {

}
