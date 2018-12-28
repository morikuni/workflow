package workflow

import (
	"context"
	"sync"
)

type WorkerStarter interface {
	Start(ctx context.Context, f func(ctx context.Context) error)
}

type Supervisor struct {
	wg     sync.WaitGroup
	logger Logger
}

var _ interface {
	WorkerStarter
} = (*Supervisor)(nil)

func NewSupervisor(logger Logger) *Supervisor {
	return &Supervisor{
		logger: logger,
	}
}

func (s *Supervisor) Start(ctx context.Context, f func(ctx context.Context) error) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		err := f(ctx)
		if err != nil {
			s.logger.Error(err.Error())
		}
	}()
}

func (s *Supervisor) Wait() error {
	s.wg.Wait()
	return nil
}
