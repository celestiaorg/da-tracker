package dbentities

import "gorm.io/gorm"

type Validator struct {
	gorm.Model
	Name    string `gorm:"type:varchar(100);uniqueIndex;not null"`
	Emails  []Email
	PeerIDs []PeerID
}

type Email struct {
	gorm.Model
	ValidatorID uint
	Address     string `gorm:"type:varchar(100);uniqueIndex;not null"`
}
