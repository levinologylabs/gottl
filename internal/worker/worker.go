// Package worker provides the worker that processes the tasks.
package worker

import (
	"context"

	"github.com/jalevin/gottl/internal/core/mailer"
	"github.com/jalevin/gottl/internal/core/tasks"
	"github.com/jalevin/gottl/internal/data/db"
	"github.com/rs/zerolog"
)

type worker struct {
	l      zerolog.Logger
	sender mailer.Sender
	db     *db.QueriesExt
}

func newWorker(l zerolog.Logger, services *db.QueriesExt, sender mailer.Sender) *worker {
	return &worker{
		l:      l,
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
		w.l.Error().Msg("invalid payload for email task")
	}

	err := w.sender.Send(data)
	if err != nil {
		w.l.Error().Err(err).
			Str("email", data.To).
			Str("subject", data.Subject).
			Msg("failed to send email")
	}
}

func (w *worker) deleteExpiredData(task tasks.Task) {
	panic("not implemented")
}
