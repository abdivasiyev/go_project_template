package job

import (
	"context"
	"fmt"
	"github.com/abdivasiyev/project_template/config"
	v1 "github.com/abdivasiyev/project_template/internal/services/v1"
	"github.com/abdivasiyev/project_template/pkg/logger"
	"github.com/abdivasiyev/project_template/pkg/sentry"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

var Module = fx.Options(fx.Invoke(New))

type Func func(ctx context.Context) error

type Job struct {
	Name     string
	Interval time.Duration
	Fn       Func
}

type Provider interface {
	Add(jobs ...Job)
}

type Params struct {
	fx.In
	Lifecycle  fx.Lifecycle
	Logger     logger.Logger
	Sentry     sentry.Handler
	JobService v1.JobServiceV1
}

type jobProvider struct {
	log    logger.Logger
	sentry sentry.Handler

	jobService v1.JobServiceV1

	jobs chan Job
	stop chan struct{}
}

func New(params Params) Provider {
	provider := &jobProvider{
		log:        params.Logger,
		sentry:     params.Sentry,
		jobService: params.JobService,
		jobs:       make(chan Job),
		stop:       make(chan struct{}),
	}

	params.Lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go provider.start(context.Background())

			provider.registerJobs()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			provider.stop <- struct{}{}
			return nil
		},
	})

	return provider
}

func (p *jobProvider) registerJobs() {
	p.Add(Job{
		Name:     "Example Job",
		Interval: 5 * time.Second,
		Fn:       p.jobService.ExampleJob,
	})
}

func (p *jobProvider) Add(jobs ...Job) {
	for _, job := range jobs {
		p.log.Info("job added", zap.String("job", job.Name), zap.Duration("job", job.Interval))
		p.jobs <- job
	}
}

func (p *jobProvider) start(ctx context.Context) {
	defer close(p.jobs)
	p.log.Info("job provider started")
	for {
		select {
		case <-ctx.Done():
			p.log.Info("job provider stopped via context")
			return
		case <-p.stop:
			p.log.Info("job provider stopped via stop channel")
			return
		case job := <-p.jobs:
			go p.startJob(ctx, job)
		}
	}
}

func (p *jobProvider) startJob(ctx context.Context, job Job) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-p.stop:
			p.log.Info("job stopped via stop channel")
			return
		case <-time.After(job.Interval):
			jobName := fmt.Sprintf("[%s] ====> ", job.Name)
			startedAt := time.Now().UTC()
			p.log.Info(jobName+"[RUNNING]", zap.Any("startedAt", startedAt.Format(config.DateTimeFormat)))
			err := job.Fn(ctx)
			finishedAt := time.Now().UTC()
			if err != nil {
				p.log.Error(jobName+"[FAILED]", zap.Error(err))
				p.sentry.HandleError(err)
				continue
			}
			p.log.Info(
				jobName+"[FINISHED]",
				zap.Any("finishedAt", finishedAt.Format(config.DateTimeFormat)),
				zap.Any("duration", finishedAt.Sub(startedAt).String()),
			)
		}
	}
}
