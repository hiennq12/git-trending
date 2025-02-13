package scheduler

import (
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

type Job interface {
	Run() error
}

type Scheduler struct {
	cron *cron.Cron
}

func NewScheduler() *Scheduler {
	// Sử dụng timezone Asia/Ho_Chi_Minh
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Printf("Error loading location, using default: %v", err)
		loc = time.Local
	}

	c := cron.New(
		cron.WithLocation(loc),
		cron.WithLogger(
			cron.VerbosePrintfLogger(log.New(log.Writer(), "cron: ", log.LstdFlags)),
		),
	)

	return &Scheduler{
		cron: c,
	}
}

func (s *Scheduler) AddJob(spec string, job Job) error {
	_, err := s.cron.AddFunc(spec, func() {
		if err := job.Run(); err != nil {
			log.Printf("Error running job: %v", err)
		}
	})
	return err
}

func (s *Scheduler) Start() {
	s.cron.Start()
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done() // Đợi tất cả jobs complete
}
