package infr

import "testing"
import "math"

func TestInit(t *testing.T) {
	var m Model
	for _, c := range []struct {
		total   int
		beliefs Beliefs
	}{
		{1, Beliefs{{0, 0}}},
		{3, Beliefs{{0, 0}, {0, 0}, {0, 0}}},
		{5, Beliefs{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}},
	} {
		m.Init(c.total)

		// check that we have the belief vector of the right size
		switch {
		case len(c.beliefs) != c.total:
			t.Errorf("wrong test: total=%d, len(beliefs)=%d",
				c.total, len(c.beliefs))
		case len(m.beliefs) != c.total:
			t.Errorf("wrong length: total=%d, lem(m.beliefs)=%d",
				c.total, len(m.beliefs))
		default:
			for i := 0; i != c.total; i++ {
				for j := 0; j != 2; j++ {
					if m.beliefs[i][j] != c.beliefs[i][j] {
						t.Errorf("wrong belief (%d, %d): got %6g, want %6g",
							i, j, m.beliefs[i][j], c.beliefs[i][j])
					}
				}
			}
		}

		// check that the prior parameters are set properly
		if m.PChurn > 1 {
			t.Errorf("PChurn must be at most 1, got PChurn=%g", m.PChurn)
		}

		if m.PBounce > 1 {
			t.Errorf("PBounce must be at most 1, got PBounce=%g", m.PBounce)
		}

		// check that the average length is half total
		if c.total > 1 && math.Abs(2./m.PChurn-float64(c.total)) > 0.5 {
			t.Errorf("average length must be half total, but got: "+
				"total=%d, average=%6g",
				c.total, 1./m.PChurn)
		}
	}
}

func TestPrior(t *testing.T) {
	var m Model
	for _, total := range []int{1, 2, 5} {
		m.Init(total)
		m.Prior()
		p_bounce := m.beliefs[0][0] / (m.beliefs[0][0] + m.beliefs[0][1])
		if math.Abs(p_bounce-m.PBounce) > 1E-6 {
			t.Errorf("wrong prior bounce probability: got %6g, wanted %6g",
				p_bounce, m.PBounce)
		}
		for _, belief := range m.beliefs[1:] {
			p_churn := belief[0] / (belief[0] + belief[1])
			if math.Abs(p_churn-m.PChurn) > 1E-6 {
				t.Errorf("wrong prior churn probability: got %6g, wanted %6g",
					p_churn, m.PChurn)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
}
