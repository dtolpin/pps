// Probabilistic query around campaign model
package query

type Query struct {
    m *Model          // model
    bandwidth float64 // initial bandwidth
    counts []int      // page counts to tune the bandwidth
}

// Function NewQuery initializes a new query from total page count,
// initial bandwidth, and counts to tune the bandwidth from.
func NewQuery(total int, bandwidth float64, counts []int) *Query {
    return &Query{NewModel(total), bandwidth, counts}
}

// Method Observe computes log probability of the data given the bandwidth.
func (q Query) Observe(bandwidth float64) (logp float64) {
    logp = - bandwidth / q.bandwidth // Exponential prior on the bandwidth

    q.m.Prior()
    for _, count := range q.counts {
        logp += q.m.Observe(count)
        q.m.Update(bandwidth, count)
    }

    return logp
}
