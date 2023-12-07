package buildinfo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/celestiaorg/validator-da-tracker/pkg/metrics/agent"
	"github.com/celestiaorg/validator-da-tracker/pkg/models/dbentities"
	"github.com/celestiaorg/validator-da-tracker/pkg/models/metricprom"
	"github.com/celestiaorg/validator-da-tracker/pkg/repository"
)

type Processor struct {
	MetricReader   agent.MetricAgent
	buildInfoCache *Cache
	peerIDRepo     repository.PeerIDRepository
}

func NewProcessor(MetricAgent agent.MetricAgent, cache *Cache, peerIDRepo repository.PeerIDRepository) *Processor {
	return &Processor{
		MetricReader:   MetricAgent,
		buildInfoCache: cache,
		peerIDRepo:     peerIDRepo,
	}
}

func (bip *Processor) ProcessBuildInfo(ctx context.Context) error {
	if bip.buildInfoCache.IsStale() {
		resp, err := bip.MetricReader.FetchMetrics("build_info")
		if err != nil {
			return fmt.Errorf("error fetching build info metrics: %w", err)
		}

		var response metricprom.BuildInfoResponse
		if err := json.Unmarshal(resp, &response); err != nil {
			return fmt.Errorf("error unmarshalling build info metrics: %w", err)
		}

		var buildInfoMetrics []metricprom.BuildInfo
		buildInfoMetrics = append(buildInfoMetrics, response.Data.Result...)

		bip.buildInfoCache.Update(buildInfoMetrics)
	}

	return bip.compareMetricsWithDB(ctx)
}

func (bip *Processor) compareMetricsWithDB(ctx context.Context) error {
	categorizedIDs, err := bip.categorizePeerIDs(ctx)
	if err != nil {
		return err
	}

	for status, peerIDs := range categorizedIDs {
		if err := bip.logPeerIDs(ctx, peerIDs, status); err != nil {
			return err
		}
	}

	return nil
}

func (bip *Processor) categorizePeerIDs(ctx context.Context) (map[dbentities.PeerIDStatus][]string, error) {
	categorizedIDs := make(map[dbentities.PeerIDStatus][]string)
	buildInfos := bip.buildInfoCache.GetMetrics()
	metricMap := bip.convertMetricsToMap(buildInfos)

	peerUniqueIDs, err := bip.peerIDRepo.FetchPeerUniqueIDs(ctx)
	if err != nil {
		return nil, err
	}

	for _, peerID := range peerUniqueIDs {
		status, err := bip.determineStatus(ctx, peerID, metricMap)
		if err != nil {
			return nil, err
		}
		if status != "" {
			categorizedIDs[status] = append(categorizedIDs[status], peerID)
		}
	}

	return categorizedIDs, nil
}

// determineStatus decides the current status of a peer ID based on metrics and history.
func (bip *Processor) determineStatus(ctx context.Context, peerID string, metricMap map[string]metricprom.BuildInfoDetails) (dbentities.PeerIDStatus, error) {
	_, existsInMetrics := metricMap[peerID]
	lastStatus, err := bip.peerIDRepo.GetLastPeerIDStatus(ctx, peerID)
	if err != nil {
		return "", err
	}

	switch {
	case !existsInMetrics && lastStatus != dbentities.PeerIDStatusMissing:
		// Peer ID is missing from metrics and is not already marked as missing
		return dbentities.PeerIDStatusMissing, nil

	case existsInMetrics && lastStatus != dbentities.PeerIDStatusActive:
		// Peer ID is present in metrics and is either new or previously marked as missing
		return dbentities.PeerIDStatusActive, nil
	}

	return "", nil // Return an empty string if no status change is needed
}

func (bip *Processor) logPeerIDs(ctx context.Context, peerIDs []string, status dbentities.PeerIDStatus) error {
	for _, peerID := range peerIDs {
		if err := bip.peerIDRepo.LogPeerID(ctx, peerID, status); err != nil {
			return err
		}
	}
	return nil
}

func (bip *Processor) convertMetricsToMap(metrics []metricprom.BuildInfo) map[string]metricprom.BuildInfoDetails {
	metricMap := make(map[string]metricprom.BuildInfoDetails)
	for _, metric := range metrics {
		details := metricprom.BuildInfoDetails{
			SemanticVersion: metric.Metric.SemanticVersion,
			Job:             metric.Metric.Job,
			SystemVersion:   metric.Metric.SystemVersion,
		}
		metricMap[metric.Metric.Instance] = details
	}
	return metricMap
}
