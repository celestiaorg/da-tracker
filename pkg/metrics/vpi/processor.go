package vpi

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/celestiaorg/validator-da-tracker/pkg/metrics/agent"
	"github.com/celestiaorg/validator-da-tracker/pkg/repository"
)

var (
	// ValidatorsPeerIDsGauge represents the gauge metrics for validators and their peer IDs.
	ValidatorsPeerIDsGauge *prometheus.GaugeVec
)

func init() {
	// Initialize the ValidatorsPeerIDsGauge
	ValidatorsPeerIDsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "validator_peer_ids",
			Help: "Gauge metric showing validators and their peer IDs",
		},
		[]string{"validator_name", "peer_id"},
	)

	// Register the gauge with Prometheus's default registry
	prometheus.MustRegister(ValidatorsPeerIDsGauge)
}

type Processor struct {
	validatorRepo repository.ValidatorRepository
	peerIDRepo    repository.PeerIDRepository
	metricAgent   agent.MetricAgent
}

func NewProcessor(validatorRepo repository.ValidatorRepository, peerIDRepo repository.PeerIDRepository, agent agent.MetricAgent) *Processor {
	return &Processor{
		validatorRepo: validatorRepo,
		peerIDRepo:    peerIDRepo,
		metricAgent:   agent,
	}
}

func (p *Processor) Process(ctx context.Context) error {
	// Fetch validators and their peer IDs
	validators, err := p.validatorRepo.GetAllValidators(ctx)
	if err != nil {
		return err
	}

	// Update the gauge values for each validator and peer ID
	for _, validator := range validators {
		for _, peerID := range validator.PeerIDs {
			ValidatorsPeerIDsGauge.WithLabelValues(validator.Name, peerID.PeerUniqueID).Set(1) // Or set to another relevant value
		}
	}

	return nil
}
