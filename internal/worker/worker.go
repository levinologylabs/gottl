// Package worker provides the worker that processes the tasks.
package worker

import (
	"context"

	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
)

type worker struct {
	db *db.QueriesExt
}

func newWorker(services *db.QueriesExt) *worker {
	return &worker{
		db: services,
	}
}

func (w *worker) run(ctx context.Context, ch <-chan tasks.Task) {
	defer func() {
		if r := recover(); r != nil {
			// log the error
			go w.run(ctx, ch)
		}
	}()

	select {
	case <-ctx.Done():
		return
	case task := <-ch:
		switch task.ID {
		case tasks.TaskIDSendEmail:
			w.sendEmail(task)
		case tasks.TaskIDDeleteExpiredData:
			w.deleteExpiredData(task)
		}
	}
}

func (w *worker) sendEmail(task tasks.Task) {
	panic("not implemented")
}

func (w *worker) deleteExpiredData(task tasks.Task) {
	panic("not implemented")
}
