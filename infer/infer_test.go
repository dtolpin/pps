package infer

import (
	"testing"
)

// In the tests below I mix pointer and value receivers intentionally
// to test both ways.

// Identity query, the value is the logp
type IdentityQuery struct{}

func (q IdentityQuery) Observe(x float64) float64 {
	return x
}

// Constant query, same logp for any value
type ConstantQuery float64

func (q *ConstantQuery) Observe(x float64) float64 {
	return float64(*q)
}

// Identity proposal, always proposes last value.
type IdentityProposal struct{}

func (p *IdentityProposal) Propose(x float64) float64 {
	return x
}

// Increment proposal, increments the last value by a constant.
type IncProposal float64

func (p IncProposal) Propose(x float64) float64 {
	return x + float64(p)
}

func TestMHAccepts(t *testing.T) {
	for _, c := range []struct {
		query    Query
		proposal Proposal
		zofi     func(int) float64
	}{
		{IdentityQuery{}, IncProposal(1.), func(i int) float64 { return float64(i) }},
		{IdentityQuery{}, &IdentityProposal{}, func(_ int) float64 { return 0. }},
		{func() *ConstantQuery {
			x := ConstantQuery(1.)
			return &x
		}(), IncProposal(3.), func(i int) float64 { return float64(3 * i) }},
	} {
		samples := make(chan float64)
		go MH(c.query, c.proposal, 0., samples)
		for i := 0; i != 10; i++ {
			z := c.zofi(i)
			x := <-samples
			if x != z {
				t.Errorf("Invalid sample from %T with %T: got %.2f, want %.2f",
					c.query, c.proposal, x, z)
			}
		}
	}
}

func TestMHRejects(t *testing.T) {
    samples := make(chan float64)
    query := IdentityQuery{}
    proposal := IncProposal(-1.)
    const N = 100
    go MH(query, proposal, float64(N), samples)
    var x float64
    for i := 0; i != N; i++ {
        x = <- samples
    }
    if x < 0. {
        t.Errorf("Last sample must be at least 0, but got %v", x)
    }
    if x < 1. {
        t.Errorf("At least one sample must be rejected for %T, %T(%v)",
            query, proposal, float64(proposal))
    }
}
