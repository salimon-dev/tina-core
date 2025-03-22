package types

import (
	"time"

	"github.com/google/uuid"
)

type VerificationType int8

const (
	VerificationTypeEmail VerificationType = 0
	VerificationTypeSMS   VerificationType = 1
)

type VerificationDomain int8

const (
	VerificationDomainRegister       VerificationDomain = 0
	VerificationDomainPasswordReset  VerificationDomain = 1
	VerificationDomainEmailUpdate    VerificationDomain = 2
	VerificationDomainUsernameUpdate VerificationDomain = 3
)

type Verification struct {
	Id        uuid.UUID          `json:"id" gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID          `json:"user_id" gorm:"type:uuid"`
	Type      VerificationType   `json:"type" gorm:"type:numeric"`
	Domain    VerificationDomain `json:"domain" gorm:"type:numeric"`
	Token     string             `json:"token" gorm:"size:16"`
	ExpiresAt time.Time          `json:"expires_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
}
