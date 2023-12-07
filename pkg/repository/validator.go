// pkg/repository/validator.go

package repository

import (
	"context"
	
	"gorm.io/gorm"

	"github.com/celestiaorg/validator-da-tracker/pkg/models/dbentities"
)

type ValidatorRepository interface {
	GetAllValidators(ctx context.Context) ([]dbentities.Validator, error)
	GetValidatorByID(ctx context.Context, id string) (*dbentities.Validator, error)
	GetValidatorByEmail(ctx context.Context, email string) (*dbentities.Validator, error)
	GetValidatorByName(ctx context.Context, name string) (*dbentities.Validator, error)
	GetValidatorsByPeerIDs(ctx context.Context, peerIDs []string) ([]dbentities.Validator, error)
}

type validatorRepository struct {
	db *gorm.DB
}

func NewValidatorRepository(db *gorm.DB) ValidatorRepository {
	return &validatorRepository{db: db}
}

func (r *validatorRepository) GetAllValidators(ctx context.Context) ([]dbentities.Validator, error) {
	var validators []dbentities.Validator
	result := r.db.WithContext(ctx).Preload("Emails").Preload("PeerIDs").Find(&validators)
	if result.Error != nil {
		return nil, result.Error
	}
	return validators, nil
}

func (r *validatorRepository) GetValidatorByID(ctx context.Context, id string) (*dbentities.Validator, error) {
	var validator dbentities.Validator

	result := r.db.WithContext(ctx).Preload("Emails").Preload("PeerIDs").First(&validator, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &validator, nil
}

func (r *validatorRepository) GetValidatorByEmail(ctx context.Context, email string) (*dbentities.Validator, error) {
	var validator dbentities.Validator

	// Join with the emails table to find the validator by email
	result := r.db.WithContext(ctx).Preload("Emails").Preload("PeerIDs").Joins("JOIN emails ON emails.validator_id = validators.id").
		Where("emails.address = ?", email).First(&validator)
	if result.Error != nil {
		return nil, result.Error
	}

	return &validator, nil
}

func (r *validatorRepository) GetValidatorByName(ctx context.Context, name string) (*dbentities.Validator, error) {
	var validator dbentities.Validator

	result := r.db.WithContext(ctx).Preload("Emails").Preload("PeerIDs").Where("name = ?", name).First(&validator)
	if result.Error != nil {
		return nil, result.Error
	}
	return &validator, nil
}

func (r *validatorRepository) GetValidatorsByPeerIDs(ctx context.Context, peerIDs []string) ([]dbentities.Validator, error) {
	var validators []dbentities.Validator
	result := r.db.WithContext(ctx).Joins("JOIN peer_ids ON peer_ids.validator_id = validators.id").
		Where("peer_ids.peer_unique_id IN ?", peerIDs).
		Group("validators.id").
		Preload("Emails").
		Preload("PeerIDs").
		Find(&validators)
	if result.Error != nil {
		return nil, result.Error
	}
	return validators, nil
}
