package utils

import (
	"testing"
	"time"
)

func TestCalc(t *testing.T) {
	startTime := time.Now()
	t.Log(startTime.UnixMilli())
	time.Sleep(1 * time.Second)
	endTime := time.Now()
	t.Log(endTime.UnixMilli())
	gap := endTime.Sub(startTime)
	t.Log(gap)
	tps := 100 / gap.Seconds()
	t.Log(tps)
}
