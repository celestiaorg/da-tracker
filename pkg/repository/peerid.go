package repository

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/celestiaorg/validator-da-tracker/pkg/models/dbentities"
)

// PeerIDRepository defines the interface for PeerID data operations.
type PeerIDRepository interface {
	FetchPeerUniqueIDs(ctx context.Context) ([]string, error)
	LogPeerID(ctx context.Context, peerID string, status dbentities.PeerIDStatus) error
	GetLastPeerIDStatus(ctx context.Context, peerID string) (dbentities.PeerIDStatus, error)
	// TODO: Add more methods here.
}

type peerIDRepository struct {
	db *gorm.DB
}

// NewPeerIDRepository creates a new instance of PeerIDRepository.
func NewPeerIDRepository(db *gorm.DB) PeerIDRepository {
	return &peerIDRepository{db: db}
}

// FetchPeerUniqueIDs fetches the list of PeerUniqueIDs from the database.
func (r *peerIDRepository) FetchPeerUniqueIDs(ctx context.Context) ([]string, error) {
	var peerIDs []dbentities.PeerID
	if err := r.db.Find(&peerIDs).Error; err != nil {
		return nil, err
	}

	var uniqueIDs []string
	for _, peerID := range peerIDs {
		uniqueIDs = append(uniqueIDs, peerID.PeerUniqueID)
	}

	return uniqueIDs, nil
}

// LogMissingPeerID logs the missing PeerID to the database.
func (r *peerIDRepository) LogPeerID(ctx context.Context, peerID string, status dbentities.PeerIDStatus) error {
	historyEntry := dbentities.PeerIDHistory{
		PeerUniqueID: peerID,
		Status:       status,
		Timestamp:    time.Now(),
	}

	return r.db.WithContext(ctx).Create(&historyEntry).Error
}

// GetLastPeerIDStatus returns the last status of the PeerID from the database.
func (r *peerIDRepository) GetLastPeerIDStatus(ctx context.Context, peerID string) (dbentities.PeerIDStatus, error) {
	var lastEntry dbentities.PeerIDHistory
	result := r.db.WithContext(ctx).Where("peer_unique_id = ?", peerID).Order("timestamp DESC").First(&lastEntry)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", nil // Or return a default status
		}
		return "", result.Error
	}
	return lastEntry.Status, nil
}
