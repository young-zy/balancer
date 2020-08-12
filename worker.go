package main

type worker struct {
	id int
	channel chan func(int)
}

func newWorker(id int, capacity int) *worker{
	return &worker{
		id:      id,
		channel: make(chan func(int), capacity),
	}
}

func (w *worker) submit(f func(int)) bool {
	select {
	case w.channel <- f:
		(<- w.channel)(w.id)
		return true
	default:
		return false
	}
}