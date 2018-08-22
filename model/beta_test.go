package model

import (
    "testing"
	"fmt"
    "math"
)

func TestBeta(t *testing.T) {
	for _, c := range []struct {
        alpha, beta, mean, variance float64
	} {
        {0, 1, 0., 0.},
        {1, 0, 1., 0.},
        {2, 3, 0.4, 0.04},
        {3, 2, 0.6, 0.04},
    } {
        dist := Beta{c.alpha, c.beta}
		mean := dist.Mean()
		variance := dist.Variance()
        if math.Abs(mean - c.mean) > epsilon {
            t.Errorf("wrong mean of %v: got %.6g, want %v",
				&dist, mean, c.mean)
		}
        if math.Abs(variance - c.variance) > epsilon {
            t.Errorf("wrong variance of %v: got %.6g, want %v",
				&dist, variance, c.variance)
		}
	}
}

func TestBetaString(t *testing.T) {
	dist := Beta{1, 2}
	if fmt.Sprint(dist) != "Beta(1, 2)" {
		t.Errorf("%#v must print as Beta(%v, %v)",
			dist, dist.Alpha, dist.Beta)
	}
	if fmt.Sprint(&dist) != "Beta(1, 2)" {
		t.Errorf("%#v must print as Beta(%v, %v)",
			&dist, dist.Alpha, dist.Beta)
	}
}
