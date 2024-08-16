package worker

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
)

var ErrQueueFull = errors.New("queue is full")

type Config struct {
	Concurrency int `json:"concurrency" conf:"default:10"`
	QueueSize   int `json:"queue_size"  conf:"default:100"`
}

var _ tasks.Queue = (*Service)(nil)

type Service struct {
	cfg   Config
	sem   chan struct{}
	queue chan tasks.Task
	db    *db.QueriesExt
}

func New(config Config, db *db.QueriesExt) *Service {
	return &Service{
		cfg:   config,
		sem:   make(chan struct{}, config.Concurrency),
		queue: make(chan tasks.Task, config.QueueSize),
		db:    db,
	}
}

func (w *Service) Start(ctx context.Context) {
	wg := sync.WaitGroup{}
	wg.Add(w.cfg.Concurrency)

	for range w.cfg.Concurrency {
		go func() {
			defer wg.Done()
			worker := newWorker(w.db)
			worker.run(ctx, w.queue)
		}()
	}

	wg.Wait()
}

// Enqueue implements tasks.Queue.
func (w *Service) Enqueue(task tasks.Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	select {
	case <-ctx.Done():
		return ErrQueueFull
	case w.queue <- task:
		return nil
	}
}
