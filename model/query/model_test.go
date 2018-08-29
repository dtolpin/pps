package query

import (
    "testing"
    "math"
)

func TestNewModel(t *testing.T) {
    // Test that the underlying the model.Model is initialized.
    total := 3
    m := NewModel(total)
    if len(m.Beliefs) != total {
        t.Errorf("Underlying model.Model is not initialized")
    }
}

func TestObserve(t *testing.T) {
    for _, c := range []struct {
        total, count int
        logp float64
    } {{1, 1, math.Log(0.5)},
       {2, 1, math.Log(0.5)},
       {2, 2, math.Log(0.5)},
       {4, 2, math.Log(0.5) + math.Log(0.5)}} {
        m := NewModel(c.total)
        m.Prior()
        logp := m.Observe(c.count)
        if  math.Abs(logp - c.logp) > epsilon {
            t.Errorf("wrong logp for total=%d, count=%d: got %.4g want %.4g",
                c.total, c.count, logp, c.logp)
        }
    }
}
