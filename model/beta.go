package model

import (
    "fmt"
)

// Beta distribution parameters
type Beta struct {Alpha, Beta float64}

func (d Beta) String () string {
    return fmt.Sprintf("Beta(%v, %v)", d.Alpha, d.Beta)
}

// Mean of the Beta distribution
func (d Beta) Mean () float64 {
    return d.Alpha / (d.Alpha + d.Beta)
}

// Variance of the Beta distribution
func (d Beta) Variance () float64 {
    v := d.Alpha + d.Beta
    return d.Alpha * d.Beta / (v * v  * (v + 1.))
}
