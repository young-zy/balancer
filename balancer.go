package main

type balancer struct{
	workers []*worker
	notifyChannel chan int			//notifies which worker is available
	waitChannel chan func(int)
}

func newBalancer(capacities ...int) *balancer{
	notifyChan := make(chan int)					// init notify channel
	workers := make([]*worker, 0, 100)				// init workers slice
	for i, capacity := range capacities{			// init workers
		worker := newWorker(i, capacity)
		workers = append(workers, worker)
		id := i
		for j := 0 ; j< capacity; j++{
			go func() {
				notifyChan <- id						// add available worker to notify
			}()
		}
	}
	return &balancer{
		workers: workers,
		notifyChannel: notifyChan,
		waitChannel: make(chan func(int), 10),
	}
}

func (bl *balancer) start(){
	go func() {
		for f := range bl.waitChannel{						// consume wait channel
			worker := bl.workers[<-bl.notifyChannel]		// consume available worker
			work := f										// important: f could change when being executed concurrently
			go func(){
				if worker.submit(work){						// submit work to worker
					bl.notifyChannel<-worker.id				// notify worker is available after done
				} else {
					bl.waitChannel <- work						// put work back to wait channel when worker refuses to do the work
				}
			}()
		}
	}()
}

func (bl *balancer) end(){
	go func() {
		for len(bl.waitChannel) > 0 {}		//wait for all woks to be consumed
		close(bl.waitChannel)
	}()
}

func (bl *balancer) submit(f func(int)){
	go func() {
		bl.waitChannel <- f										//adds work to wait channel
	}()
}

