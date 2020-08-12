package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

type work struct{
	workId int
}

func (w *work) doWork(workerId int){
	fmt.Println("work", w.workId, "started on worker ", workerId)
	time.Sleep(time.Duration(1) * time.Second)
	fmt.Println("\twork ", w.workId, "ended on worker ", workerId)
	wg.Done()
}

func main() {
	ba := newBalancer(1, 1, 1, 1)		// defines workers' capacity
	ba.start()
	for i := 0; i < 10; i++ {
		w := &work{i}					// create work
		ba.submit(w.doWork)					// submit work to balancer
		wg.Add(1)
	}
	wg.Wait()
	ba.end()
}