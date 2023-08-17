package scheduler

import (
	"time"

	"github.com/go-co-op/gocron"
)

type Scheduler struct {
	ctx                 Context
	cron                *gocron.Scheduler
	operationsRegistrar OperationRegistrar
}

func New(context Context) *Scheduler {
	return &Scheduler{
		ctx:  context,
		cron: gocron.NewScheduler(time.UTC),
	}
}

func (r *Scheduler) SetOperationsRegistrar(registrar OperationRegistrar) {
	r.operationsRegistrar = registrar
}

func (r *Scheduler) StartAsync() {
	if r.operationsRegistrar != nil {
		r.operationsRegistrar(r.ctx, r.cron)
	}

	r.cron.StartAsync()
}

func (r *Scheduler) Stop() {
	r.cron.Stop()
}
