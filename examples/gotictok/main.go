package main

import (
	"fmt"
	"time"

	"github.com/efjoubert/goforit/gotictok"
)

func main() {
	var jb = gotictok.ScheduleJob("jb1", 10*time.Millisecond, func() {
		fmt.Println("bla", time.Now())
	})
	jb.Run()
	time.Sleep(30 * time.Second)
	jb.Done()
}
