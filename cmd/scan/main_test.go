package main

import (
	"bitbucket.org/dtolpin/pps/model"
	"reflect"
	"testing"
)

func TestMakeHeader(t *testing.T) {
	for _, c := range []struct {
		total  int
		header []string
	}{
		{0, []string{"iline", "mean", "variance"}},
		{1, []string{"iline", "mean", "variance", "a1", "b1"}},
		{2, []string{"iline", "mean", "variance",
			"a1", "b1",
			"a2", "b2"}}} {

		m := model.NewModel(c.total)
		header := makeHeader(m)
		if !reflect.DeepEqual(header, c.header) {
			t.Errorf("wrong header: total=%v, got %#v, want %#v",
				c.total, header, c.header)
		}
	}
}

func TestMakeRecord(t *testing.T) {
	for _, c := range []struct {
		total  int
        iline int
		record []string
	}{
        {0, 10, []string{"10", "1.0", "0.0"}},
		{1, 20, []string{"20", "1.0", "0.0", "1.0", "0.0"}},
		{2, 30, []string{"30", "1.0", "0.0",
            "1.0", "0.0",
            "1.0", "0.0"}}} {

        // avoid zero evidence
		m := model.NewModel(c.total)
        for i := 0; i != c.total; i++ {
            m.Beliefs[i][0] = 1.
        }

		record := makeRecord(c.iline, m, "%.1f")
		if !reflect.DeepEqual(record, c.record) {
			t.Errorf("wrong record: total=%v, got %#v, want %#v",
				c.total, record, c.record)
		}
	}
}
