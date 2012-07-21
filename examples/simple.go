package main

import (
	"../rrd"
	"log"
	"os"
	"time"
)

func main() {
	defer func() {
		log.Printf("Purging test.rrd")
		_ = os.Remove("test.rrd")
	}()

	log.Printf("Creating test.rrd")
	err := rrd.Create("test.rrd", int64(10), time.Now().Unix(), []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}
	log.Printf("Waiting 1 sec so that we aren't reporting at creation time")
	time.Sleep(1 * time.Second)

	log.Printf("Updating test.rrd")
	value := rrd.RrdValue{
		Time:  time.Now(),
		Value: 15,
	}

	err = rrd.UpdateValues("test.rrd", "ok", []rrd.RrdValue{
		value,
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	lastTime := rrd.Last("test.rrd")
	log.Printf("Last time = %s", lastTime.String())

	log.Printf("Everything is OK")
}
