package main

import (
	"testing"
)

func TestPagePChurnErr(t *testing.T) {
	for _, c := range []struct {
		x, m, s float64
	}{{3., 0., 2.},
		{1., 2., 0.},
		{5., 2., 1.}} {
		page := pagePChurnErr(c.x, c.m, c.s)
		if len(page) < 4 {
			t.Errorf("wrong page representation length: got %v,  want >= 4",
				len(page))
		}

		// all points have the same x coordinate
		for i, p := range page {
			if p.X != c.x {
				t.Errorf("wrong x coordinate at point %d: got %v, want %v",
					i, p.X, c.x)
			}
		}

		// extreme points are set to the mean for proper connection
		// between points
		for _, i := range []int{0, len(page) - 1} {
			if page[i].Y != c.m {
				t.Errorf("page[%d].Y must be set to the mean: got %v, want %v",
					i, page[i].Y, c.m)
			}
		}

		// some intermediate points deviate from the mean
		if c.s > 0 {
			offmean := false
			for i := 1; i < len(page)-1; i++ {
				if page[i].Y != c.m {
					offmean = true
					break
				}
			}
			if !offmean {
				t.Errorf("some points must deviate from the mean, but got "+
					"mean = %.4g, page = %v",
					c.m, page)
			}
		}
	}
}

func TestPpsMeanErr(t *testing.T) {
	for _, c := range []struct {
		m, s float64
	}{{0., 1.}, {2., 0.}, {2., 3.}} {
		marker := ppsMeanErr(c.m, c.s)
        if marker[len(marker) / 2].X != c.m {
            t.Errorf("the midpoint must be at the mean %.4g, but got %.4g",
                c.m, marker[len(marker) / 2].X)
        }
        for i := 1; i < len(marker); i++ {
            if marker[i].X < marker[i - 1].X {
                t.Errorf("X coordinates must be ordered, but go %#v", 
                    marker)
            }
        }
    }
}
