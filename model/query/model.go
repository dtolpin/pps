package query

import (
	"bitbucket.org/dtolpin/pps/model"
	"math"
)

// Model extends model.Model with Observe and prior bandwidth. Observe is used
// in the probabilistic query.
type Model struct {
	model.Model
	bandwidth float64
}

// Function NewModel creates and initializes a model.
func NewModel(total int) *Model {
	m := new(Model)
	m.Init(total)
	return m
}

// Method observe computes log probability of the page count
// given the model.
func (m *Model) Observe(count int) float64 {
	logp := 0.
	for i := 0; i != len(m.Beliefs); i++ {
		var j int // selects either alpha or beta
		if i < count-1 {
			j = 1
		} else {
			j = 0
		}
		logp += math.Log(m.Beliefs[i][j]) -
			math.Log(m.Beliefs[i][0]+m.Beliefs[i][1])
		if j == 0 { // reached the last page of the session
			break
		}
	}
	return logp
}
