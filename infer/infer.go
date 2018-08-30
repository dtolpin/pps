// Package infer implements basic probabilistic inference.
package infer

import (
	"math"
	"math/rand"
)

// Interface Query has Observe method which computes the log-likelihood
// given the parameters.
type Query interface {
	Observe(x float64) (logp float64)
}

// MH sampler needs a proposal. Interface Proposal has Propose method
// which takes a value and proposes a new one.
type Proposal interface {
	Propose(x float64) float64
}

// Function MH takes a query and a proposal and writes samples to the
// samples channel.
func MH(query Query, proposal Proposal, x float64, samples chan<- float64) {
	logp := query.Observe(x)
	for {
		samples <- x
		x0, logp0 := x, logp
		x = proposal.Propose(x)
		logp = query.Observe(x)
		if logp-logp0 < math.Log(1.-rand.Float64()) {
			x, logp = x0, logp0
		}
	}
}

// Type RandomWalk implements the proposal interface and provides a random walk
// proposal.
type RandomWalk float64

func (std RandomWalk) Propose(x float64) float64 {
	return x + rand.NormFloat64()*float64(std)
}
