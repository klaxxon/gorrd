package main

import (
	"../rrd"
	"log"
	"os"
	"time"
)

func main() {
	resolution := 1

	defer func() {
		log.Printf("Purging test.rrd")
		_ = os.Remove("test.rrd")
	}()

	log.Printf("Creating test.rrd")
	err := rrd.Create("test.rrd", uint64(resolution), time.Now(), []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	for i := 1; i <= 5; i++ {
		log.Printf("Waiting %d sec so that we aren't reporting at creation time", resolution)
		time.Sleep(time.Duration(resolution) * time.Second)
		log.Printf("Updating test.rrd, iter %d", i)
		value := rrd.RrdValue{
			Time:  time.Now(),
			Value: int64(15 + i),
		}
		err = rrd.UpdateValues("test.rrd", "ok", []rrd.RrdValue{
			value,
		})
		if err != nil {
			log.Printf("Error: %s", err)
			return
		}
	}

	lastTime := rrd.Last("test.rrd")
	log.Printf("Last time = %s", lastTime.String())

	dsCount, dsNames, data, err := rrd.Fetch("test.rrd", rrd.CF_AVERAGE, uint64(time.Now().Unix()-30), uint64(time.Now().Unix()), uint64(resolution))
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	// (dsCount uint64, dsNames []string, data [][]float64, err error) {
	log.Printf("dsCount = %d, dsNames.len = %d, data.len = %d", dsCount, len(dsNames), len(data))

	err = rrd.Dump("test.rrd", "/dev/stderr")

	log.Printf("Everything is OK")
}
