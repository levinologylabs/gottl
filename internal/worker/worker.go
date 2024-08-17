// Package worker provides the worker that processes the tasks.
package worker

import (
	"context"

	"github.com/jalevin/gottl/internal/core/mailer"
	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
)

type worker struct {
	sender mailer.Sender
	db     *db.QueriesExt
}

func newWorker(services *db.QueriesExt, sender mailer.Sender) *worker {
	return &worker{
		db:     services,
		sender: sender,
	}
}

func (w *worker) run(ctx context.Context, ch <-chan tasks.Task) {
	defer func() {
		if r := recover(); r != nil {
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
	data, ok := task.Payload.(mailer.Message)
	if !ok {
		println("Task integrity error")
		panic("task integrity error")
	}

	err := w.sender.Send(data)
	if err != nil {
		println(err.Error())
		panic(err)
	}
}

func (w *worker) deleteExpiredData(task tasks.Task) {
	panic("not implemented")
}
