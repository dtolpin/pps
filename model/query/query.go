// Probabilistic query around campaign model. The query provides the Observe
// method which computes unnormalized log likelihood of parameters given the
// data. With this Query at hand, inference can be performed on the parameters.
package query

// Structure Query encapsulates the model, the data and the prior on the
// parameters. An alternative approach would be to define function NewQuery
// which returns the observe function.
type Query struct {
	*Model            // model
	bandwidth float64 // initial bandwidth
	counts    []int   // page counts to tune the bandwidth
}

// Function NewQuery initializes a new query from total page count,
// initial bandwidth, and counts to tune the bandwidth from.
func NewQuery(total int, bandwidth float64, counts []int) *Query {
	return &Query{NewModel(total), bandwidth, counts}
}

// Method Observe computes log probability of the data given the bandwidth.
func (q *Query) Observe(bandwidth float64) (logp float64) {
	logp = -bandwidth / q.bandwidth // Exponential prior on the bandwidth

	q.Prior()
	for _, count := range q.counts {
		logp += q.Model.Observe(count)
		q.Update(bandwidth, count)
	}

	return logp
}
