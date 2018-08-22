package model

import (
	"math"
	"testing"
)

// comparison accuracy
const epsilon = 1E-6

func TestNewModel(t *testing.T) {
	for _, c := range []struct {
		total   int
		beliefs Beliefs
	}{
		{1, Beliefs{{0, 0}}},
		{3, Beliefs{{0, 0}, {0, 0}, {0, 0}}},
		{5, Beliefs{{0, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}},
	} {
		m := NewModel(c.total)

		// check that we have the belief vector of the right size
		switch {
		case len(c.beliefs) != c.total:
			t.Errorf("wrong test: total=%d, len(beliefs)=%d",
				c.total, len(c.beliefs))
		case len(m.Beliefs) != c.total:
			t.Errorf("wrong length: total=%d, lem(m.Beliefs)=%d",
				c.total, len(m.Beliefs))
		default:
			for i := 0; i != c.total; i++ {
				for j := 0; j != 2; j++ {
					if m.Beliefs[i][j] != c.beliefs[i][j] {
						t.Errorf("wrong belief (%d, %d): got %.6g, want %.6g",
							i, j, m.Beliefs[i][j], c.beliefs[i][j])
					}
				}
			}
		}

		// check that the prior parameters are set properly
		if m.pChurn > 1 {
			t.Errorf("pChurn must be at most 1, got pChurn=%g", m.pChurn)
		}

		if m.pBounce > 1 {
			t.Errorf("pBounce must be at most 1, got pBounce=%g", m.pBounce)
		}

		// check that the average length is half total
		if c.total > 1 && math.Abs(2./m.pChurn-float64(c.total)) > 0.5 {
			t.Errorf("average length must be half total, but got: "+
				"total=%d, average=%.6g",
				c.total, 1./m.pChurn)
		}
	}
}

func TestPrior(t *testing.T) {
	for _, total := range []int{1, 2, 5} {
		m := NewModel(total)
		m.Prior()
		pBounce := m.Beliefs[0][0] / (m.Beliefs[0][0] + m.Beliefs[0][1])
		if math.Abs(pBounce-m.pBounce) > epsilon {
			t.Errorf("wrong prior bounce probability: got %.6g, wanted %.6g",
				pBounce, m.pBounce)
		}
		for _, belief := range m.Beliefs[1:] {
			pChurn := belief[0] / (belief[0] + belief[1])
			if math.Abs(pChurn-m.pChurn) > epsilon {
				t.Errorf("wrong prior churn probability: got %.6g, wanted %.6g",
					pChurn, m.pChurn)
			}
		}
	}
}

func TestUpdate(t *testing.T) {
	// check with bandwidth high enough to keep all evidence
	m := NewModel(5)
	bandwidth := 1000.
	for k, c := range []struct {
		pps     int
		beliefs Beliefs
	}{
		{1, Beliefs{{1, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}},
		{3, Beliefs{{1, 1}, {0, 1}, {1, 0}, {0, 0}, {0, 0}}},
		{2, Beliefs{{1, 2}, {1, 1}, {1, 0}, {0, 0}, {0, 0}}},
		{5, Beliefs{{1, 3}, {1, 2}, {1, 1}, {0, 1}, {1, 0}}},
		{8, Beliefs{{1, 4}, {1, 3}, {1, 2}, {0, 2}, {1, 1}}},
	} {
		m.Update(bandwidth, c.pps)
		for i := 0; i != len(m.Beliefs); i++ {
			for j := 0; j != len(m.Beliefs[i]); j++ {
				if math.Abs(m.Beliefs[i][j]-c.beliefs[i][j]) > epsilon {
					t.Errorf("%d (bandwidth=%g): wrong belief [%d, %d]: wanted %4g, got %g",
						k, bandwidth, i, j, c.beliefs[i][j], m.Beliefs[i][j])
				}
			}
		}
	}

	// check with low bandwidth
	m = NewModel(5)
	bandwidth = 2.
	for k, c := range []struct {
		pps     int
		beliefs Beliefs
	}{
		{1, Beliefs{{1, 0}, {0, 0}, {0, 0}, {0, 0}, {0, 0}}},
		{3, Beliefs{{1, 1}, {0, 1}, {1, 0}, {0, 0}, {0, 0}}},
		{2, Beliefs{{1 / 1.5, 2 / 1.5}, {1, 1}, {1, 0}, {0, 0}, {0, 0}}},
		{5, Beliefs{{1 / 2.25, 3.5 / 2.25}, {1 / 1.5, 2 / 1.5}, {1, 1}, {0, 1}, {1, 0}}},
		{8, Beliefs{{1 / 3.375, 5.75 / 3.375}, {1 / 2.25, 3.5 / 2.25}, {1 / 1.5, 2 / 1.5}, {0, 2}, {1, 1}}},
	} {
		m.Update(bandwidth, c.pps)
		for i := 0; i != len(m.Beliefs); i++ {
			for j := 0; j != len(m.Beliefs[i]); j++ {
				if math.Abs(m.Beliefs[i][j]-c.beliefs[i][j]) > epsilon {
					t.Errorf("%d (bandwidth=%g): wrong belief [%d, %d]: wanted %4g, got %g",
						k, bandwidth, i, j, c.beliefs[i][j], m.Beliefs[i][j])
				}
			}
		}
	}
}
