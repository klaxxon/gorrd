package rrd

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestCreateDs(t *testing.T) {
	cleanup()

	values := []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	}
	err := Create("test.rrd", 5, time.Now().Add(-10*time.Second), values)

	if err != nil {
		t.Errorf("Error: %s", err)
	}
}

func TestCreateError(t *testing.T) {
	cleanup()

	values := []string{
		"DS:ok:GAUGE:600:0:U",
	}
	err := Create("test.rrd", 5, time.Now().Add(-10*time.Second), values)
	if err == nil {
		t.Fatalf("Expected error: you must define at least one Round Robin Archive")
	}
}

func TestUpdate(t *testing.T) {
	return // TODO: FIXME: FIX TEST AT SOME POINT

	cleanup()

	values := []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	}
	err := Create("test.rrd", 15, time.Now().Add(-10*time.Second), values)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Log("Waiting 2 seconds for sanity")
	time.Sleep(time.Duration(2) * time.Second)
	err = Update("test.rrd", "ok", []string{fmt.Sprintf("%d:%d", time.Now(), 15)})
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestUpdateInvalidDs(t *testing.T) {
	cleanup()

	values := []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	}
	err := Create("test.rrd", 15, time.Now().Add(-10*time.Second), values)
	if err != nil {
		t.Fatalf("Error: %s", err)
	}
	t.Log("Waiting 2 seconds for sanity")
	time.Sleep(time.Duration(2) * time.Second)
	err = Update("test.rrd", "fail", []string{fmt.Sprintf("%d:%d", time.Now(), 15)})
	if err == nil {
		t.Fatalf("Expected error: unknown DS name 'fail'", err)
	}
}

func TestFetch(t *testing.T) {
	dsCount, dsNames, data, err := Fetch("test.rrd", CF_AVERAGE, time.Now().Unix()-(30*3600*24), time.Now().Unix(), uint64(60))
	if err != nil {
		t.Errorf("Error: %s", err)
		return
	}
	t.Logf("dsCount = %d\n", dsCount)
	for n := 0; n < int(dsCount); n++ {
		t.Logf("dsName[%d] = %s\n", n, dsNames[n])
		for k, v := range data[n] {
			//t.Logf("\t%d = %d\n", k, v)
		}
	}

}

func cleanup() {
	os.Remove("test.rrd")
	clearError()
}
