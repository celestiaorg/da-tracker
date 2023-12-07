package metrics

import (
	"context"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/celestiaorg/validator-da-tracker/pkg/metrics/agent"
	"github.com/celestiaorg/validator-da-tracker/pkg/metrics/buildinfo"
	"github.com/celestiaorg/validator-da-tracker/pkg/metrics/vpi"
	"github.com/celestiaorg/validator-da-tracker/pkg/repository"
)

const (
	timeout = time.Second * 120
)

func NewPrometheusClient() *http.Client {
	return &http.Client{
		Timeout: timeout,
	}
}

// StartIntervalTask starts a task that runs at specified intervals.
func StartIntervalTask(ctx context.Context, interval time.Duration, task func(context.Context) error) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := task(ctx); err != nil {
				return err
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// StartMetricsScraper starts the metric scraping in a concurrent fashion and writes to a DB if peers are missing.
func StartMetricsScraper(ctx context.Context, db *gorm.DB, client *http.Client, endpoint string, token string, interval time.Duration) error {
	peerIDRepo := repository.NewPeerIDRepository(db)
	scaleWayAgent := agent.NewPromScaleWayAgent(client, endpoint, token)
	buildInfoCache := buildinfo.NewCache()
	buildInfoProcessor := buildinfo.NewProcessor(scaleWayAgent, buildInfoCache, peerIDRepo)

	return StartIntervalTask(ctx, interval, func(ctx context.Context) error {
		return buildInfoProcessor.ProcessBuildInfo(ctx)
	})
}

// StartMetricsProcessor starts the metric processing and pushing in a concurrent fashion.
func StartMetricsProcessor(ctx context.Context, db *gorm.DB, client *http.Client, endpoint string, token string, interval time.Duration) error {
	validatorRepo := repository.NewValidatorRepository(db)
	peerIDRepo := repository.NewPeerIDRepository(db)
	scaleWayAgent := agent.NewPromScaleWayAgent(client, endpoint, token)
	vpiproc := vpi.NewProcessor(validatorRepo, peerIDRepo, scaleWayAgent)

	return StartIntervalTask(ctx, interval, vpiproc.Process)
}
