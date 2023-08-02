package jtime

import (
	"testing"
	"time"
)

func TestParseLen(t *testing.T) {
	tNow := []time.Time{
		time.Now(),
		time.Now().AddDate(0, 0, 2),
	}
	for i, s := range timeFormats0 {
		t.Run(s, func(t *testing.T) {
			for _, t1 := range tNow {
				s1 := t1.Format(s)
				if len(s1) != len(s) {
					t.Errorf("[%d] length not equal: %s => %s", i, s, s1)
				}
			}
		})
	}

}
