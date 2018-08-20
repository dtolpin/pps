// Package infr contains model and inference for predicting
// pages per session based on page counts.
package infr

type Beliefs = [][2]float64

type Model struct {
    beliefs Beliefs
}

func (m *Model) Init(total int) {
    m.beliefs = make(Beliefs, total)
}

func (m *Model) Update(bandwidth int, count int) {
    for i := 0 ;; i ++ {
        if i < len(m.beliefs) {
            m.beliefs[i][1] ++
        } else {
            m.beliefs[i][0] ++
            break
        }
    }
}
