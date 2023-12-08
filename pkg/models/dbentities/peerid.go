package dbentities

import (
	"time"

	"gorm.io/gorm"
)

type PeerIDStatus string

const (
	PeerIDStatusMissing PeerIDStatus = "Missing"
	PeerIDStatusActive  PeerIDStatus = "Active"
)

type PeerID struct {
	gorm.Model
	ValidatorID  uint
	PeerUniqueID string `gorm:"type:varchar(100);uniqueIndex;not null"`
}

type PeerIDHistory struct {
	gorm.Model
	PeerUniqueID string       `gorm:"type:varchar(100);index;not null"`
	Status       PeerIDStatus `gorm:"type:varchar(100);not null"`
	Timestamp    time.Time    `gorm:"type:timestamp;not null"`
}
