package model

// Beta distribution parameters
type beta struct {alpha, beta float64}

// Mean of the Beta distribution
func (d *beta) mean () float64 {
    return d.alpha / (d.alpha + d.beta)
}

// Variance of the beta distribution
func (d *beta) variance () float64 {
    v := d.alpha + d.beta
    return d.alpha * d.beta / (v * v  * (v + 1.))
}
