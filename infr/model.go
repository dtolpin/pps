// Package infr contains model and inference for predicting
// pages per session based on page counts.
package infr

// Beliefs is an array of alpha, beta parameters of Beta distribution,
// a tuple per page.
type Beliefs [][2]float64

type Model struct {
	beliefs         Beliefs
	PBounce, PChurn float64
}

// Method Init initializes the model.
func (m *Model) Init(total int) {
	m.beliefs = make(Beliefs, total)
	// half the visitors bounce off the first page
	m.PBounce = 0.5

	// an average visitor views half the pages
	m.PChurn = 2. / float64(total)
	if m.PChurn > 1. {
		m.PChurn = 1.
	}
}

// Method Prior resets the model to the prior beliefs.
func (m *Model) Prior() {
	m.beliefs[0][0], m.beliefs[0][1] = m.PBounce, 1.-m.PBounce
	for i := 1; i != len(m.beliefs); i++ {
		m.beliefs[i][0], m.beliefs[i][1] = m.PChurn, 1.-m.PChurn
	}
}

// Method Update updates the model with evidence.
func (m *Model) Update(bandwidth float64, count int) {
	for i := 0; i != len(m.beliefs); i++ {
		var j int // selects either alpha or beta
		if i < count-1 {
			j = 1
		} else {
			j = 0
		}
        m.beliefs[i][j] ++
		// if the evidence exceeds the bandwidth, scale down
		evidence := m.beliefs[i][0] + m.beliefs[i][1]
		if evidence > bandwidth {
			scale := bandwidth / evidence
			m.beliefs[i][0] *= scale
			m.beliefs[i][1] *= scale
		}
		if j == 0 { // reached the last page of the session
			break
		}
	}
}
