package types

import (
	"time"

	"github.com/google/uuid"
)

type InvitationStatus int8

const (
	InvitationStatusActive   InvitationStatus = 1
	InvitationStatusInActive InvitationStatus = 2
)

type Invitation struct {
	Id             uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	CreatedBy      uuid.UUID        `json:"created_by" gorm:"type:uuid;not null"`
	Code           string           `json:"code" gorm:"type:string;size:16;not null;unique"`
	UsageRemaining int16            `json:"usage_remaining" gorm:"type:numeric;not null"`
	ExpiresAt      time.Time        `json:"expires_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
	Status         InvitationStatus `json:"status" gorm:"type:numeric"`
	CreatedAt      time.Time        `json:"created_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
	UpdatedAt      time.Time        `json:"updated_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
}
