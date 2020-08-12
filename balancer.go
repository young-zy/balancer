package main

import (
	"sync"
)

type balancer struct {
	availId       int // available id
	workers       map[int]*worker
	notifyChannel chan int //notifies which worker is available
	waitChannel   chan func(int)
}

func newBalancer(capacities ...int) *balancer {
	notifyChan := make(chan int)                      // init notify channel
	workers := make(map[int]*worker, len(capacities)) // init workers slice
	bl := &balancer{
		availId:       0,
		workers:       workers,
		notifyChannel: notifyChan,
		waitChannel:   make(chan func(int), 10),
	}
	for _, capacity := range capacities { // init workers
		bl.add(capacity)
	}
	return bl
}

func (bl *balancer) start() {
	go func() {
		for f := range bl.waitChannel { // consume wait channel
			worker := bl.workers[<-bl.notifyChannel] // consume available worker
			work := f                                // important: f could change when being executed concurrently
			go func() {
				if worker.submit(work) { // submit work to worker
					bl.notifyChannel <- worker.id // notify worker is available after done
				} else {
					bl.waitChannel <- work // put work back to wait channel when worker refuses to do the work
				}
			}()
		}
	}()
}

func (bl *balancer) end() {
	go func() {
		for len(bl.waitChannel) > 0 {
		} //wait for all woks to be consumed
		close(bl.waitChannel)
	}()
}

func (bl *balancer) add(capacity int) {
	var mutex sync.RWMutex
	mutex.Lock()
	newId := bl.availId
	bl.availId++
	mutex.Unlock()
	worker := newWorker(newId, capacity)
	bl.workers[newId] = worker
	for i := 0; i < capacity; i++ {
		go func() {
			bl.notifyChannel <- newId
		}()
	}
}

func (bl *balancer) remove(id int) {
	worker := bl.workers[id].channel
	bl.workers[id].channel = nil
	go func() {
		for len(worker) > 0 {
		}
		worker = nil
	}()
}

func (bl *balancer) submit(f func(int)) {
	go func() {
		bl.waitChannel <- f //adds work to wait channel
	}()
}
