package infr

import "testing"

func TestInit(t *testing.T) {
    var m Model
    for _, c := range []struct {
        total int
        beliefs Beliefs
    } {
        {1, Beliefs {{0, 0}}},
        {3, Beliefs {{0, 0}, {0, 0}, {0, 0}}},
    } {
        m.Init(c.total)
        switch {
        case len(c.beliefs) != c.total:
            t.Errorf("wrong test: total=%d, len(beliefs)=%d",
                     c.total, len(c.beliefs))
        case len(m.beliefs) != c.total:
            t.Errorf("wrong length: total=%d, lem(m.beliefs)=%d",
                     c.total, len(m.beliefs))
        default:
            for i := 0; i != c.total; i ++ {
                for j := 0; j != 2; j ++ {
                    if m.beliefs[i][j] != c.beliefs[i][j] {
                        t.Errorf("wrong belief (%d, %d): got %6g, want %6g",
                                 i, j, m.beliefs[i][j], c.beliefs[i][j])
                    }
                }
            }
        }
    }
}
