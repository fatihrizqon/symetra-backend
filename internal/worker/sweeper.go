package worker

import (
	"context"
	"time"

	"github.com/fatihrizqon/gofiber-microservice/internal/repository"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
)

type JobSweeper struct {
	repo        repository.IRedisJobRepository
	asynqClient *asynq.Client
	log         *logrus.Logger
}

func NewJobSweeper(repo repository.IRedisJobRepository, asynqClient *asynq.Client, log *logrus.Logger) *JobSweeper {
	return &JobSweeper{
		repo:        repo,
		asynqClient: asynqClient,
		log:         log,
	}
}

func (s *JobSweeper) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	s.log.Infof("Starting RedisJob Sweeper with interval %v", interval)

	for {
		select {
		case <-ctx.Done():
			s.log.Info("Stopping RedisJob Sweeper")
			return
		case <-ticker.C:
			s.sweep()
		}
	}
}

func (s *JobSweeper) sweep() {
	jobs, err := s.repo.GetPendingJobs(100) // fetch batch
	if err != nil {
		s.log.Errorf("JobSweeper: failed to fetch pending jobs: %v", err)
		return
	}

	if len(jobs) == 0 {
		return
	}

	for _, job := range jobs {
		task := asynq.NewTask(job.Type, job.Payload)
		info, err := s.asynqClient.Enqueue(task)
		if err != nil {
			s.log.Errorf("JobSweeper: failed to enqueue job %s: %v", job.ID, err)
			_ = s.repo.UpdateStatus(job.ID, "FAILED", err.Error())
			continue
		}

		s.log.Infof("JobSweeper: successfully enqueued job %s to asynq (queue=%s id=%s)", job.ID, info.Queue, info.ID)
		_ = s.repo.UpdateStatus(job.ID, "COMPLETED", "")
	}
}
