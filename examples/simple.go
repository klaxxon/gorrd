package main

import (
	"../rrd"
	"fmt"
	"log"
	"time"
)

func main() {
	err := rrd.Create("test.rrd", int64(10), time.Now().Unix(), []string{
		"DS:ok:GAUGE:600:0:U",
		"RRA:AVERAGE:0.5:1:25920",
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	err = rrd.Update("test.rrd", "ok", []string{
		fmt.Sprintf("%d:%d", time.Now(), 15),
	})
	if err != nil {
		log.Printf("Error: %s", err)
		return
	}

	log.Printf("Everything is OK")
}
