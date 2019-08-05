package gotictok

import (
	"fmt"
	"sync"
	"time"
)

//Job - Job structure
type Job struct {
	lbl       string
	done      chan bool
	callfuncs []interface{}
	intrvl    time.Duration
	jblck     *sync.Mutex
}

//ScheduleJob - invoke new Schedulable(Job)
//label - job label
//interval - interval in nano seconds
//callfunc - (optional) ...interface{} list of functions that will be actioned
func ScheduleJob(label string, interval time.Duration, callfunc ...interface{}) (jb *Job) {
	jb = &Job{jblck: &sync.Mutex{}, lbl: label, done: make(chan bool), callfuncs: callfunc[:], intrvl: interval}
	return
}

//Run - run or activate Job
func (jb *Job) Run() *Job {
	var tck = time.NewTicker(jb.intrvl)
	var jobaction = func(t time.Time) {
		jb.jblck.Lock()
		defer jb.jblck.Unlock()
		jb.performActions(t)

	}
	go scheduler(false, tck, jb.done, jobaction)
	return jb
}

func (jb *Job) performActions(t time.Time) {
	fmt.Println(jb.lbl, "perform action", t)
}

func scheduler(doNow bool, tick *time.Ticker, done chan bool, jobaction func(t time.Time)) {
	if doNow {
		jobaction(time.Now())
	}
	var t time.Time
	for {
		select {
		case t = <-tick.C:
			jobaction(t)
		case isDone := <-done:
			if isDone {
				tick.Stop()
				select {
				case <-tick.C:
				default:
				}
				tick = nil
				done <- true
				done = nil
				jobaction = nil
				return
			}
		}
	}
}

//Done - job done and can not be used after this
func (jb *Job) Done() {
	jb.done <- true
	<-jb.done
	jb.callfuncs = nil
	close(jb.done)
	jb.done = nil
	if jb.jblck != nil {
		jb.jblck = nil
	}
	jb = nil
}
