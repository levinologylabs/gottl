package tasks

// NoopQueue is a no-op implementation of the Queue interface.
var NoopQueue Queue = noopqueue{}

type noopqueue struct{}

func (q noopqueue) Enqueue(task Task) error {
	return nil
}
