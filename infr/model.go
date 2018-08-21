// Package infr contains model and inference for predicting
// pages per session based on page counts.
package infr

// Beliefs is an array of alpha, beta parameters of Beta distribution,
// a tuple per page.
type Beliefs [][2]float64

// Model is the model of a session. The fields are the current beliefs,
// and the prior probability of bouncing off the first page and churning
// at any subsequent page.
type Model struct {
	Beliefs         Beliefs
	pBounce, pChurn float64
}

// Method Init initializes the model.
func (m *Model) Init(total int) {
	m.Beliefs = make(Beliefs, total)

    // We set prior probabilities here so that they
    // can be re-used for resetting the beliefs:
	//   * half the visitors bounce off the first page,
	m.pBounce = 0.5
	//   * an average visitor views half the pages.
	m.pChurn = 2. / float64(total)
	if m.pChurn > 1. {
		m.pChurn = 1.
	}
}

// Method Prior resets the model to the prior beliefs.
func (m *Model) Prior() {
	m.Beliefs[0][0], m.Beliefs[0][1] = m.pBounce, 1.-m.pBounce
	for i := 1; i != len(m.Beliefs); i++ {
		m.Beliefs[i][0], m.Beliefs[i][1] = m.pChurn, 1.-m.pChurn
	}
}

// Method Update updates the model with evidence.
func (m *Model) Update(bandwidth float64, count int) {
	for i := 0; i != len(m.Beliefs); i++ {
		var j int // selects either alpha or beta
		if i < count-1 {
			j = 1
		} else {
			j = 0
		}
        m.Beliefs[i][j] ++
		// if the evidence exceeds the bandwidth, scale down
		evidence := m.Beliefs[i][0] + m.Beliefs[i][1]
		if evidence > bandwidth {
			scale := bandwidth / evidence
			m.Beliefs[i][0] *= scale
			m.Beliefs[i][1] *= scale
		}
		if j == 0 { // reached the last page of the session
			break
		}
	}
}
