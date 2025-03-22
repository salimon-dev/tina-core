package types

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type EntityStatus int8
type EntityPermission int8

const (
	EntityStatusActive   EntityStatus = 1
	EntityStatusInActive EntityStatus = 2
)

const (
	EntityPermissionPublic      EntityPermission = 1
	EntityPermissionInternal    EntityPermission = 2
	EntityPermissionPrivate     EntityPermission = 3
	EntityPermissionDevelopment EntityPermission = 4
)

func EntityStatusToString(status EntityStatus) string {
	switch status {
	case EntityStatusActive:
		return "active"
	case EntityStatusInActive:
		return "inactive"
	default:
		return "none"
	}
}

type Entity struct {
	Id          uuid.UUID        `json:"id" gorm:"type:uuid;primaryKey"`
	Name        string           `json:"name" gorm:"size:32;unique;not null"`
	Description string           `json:"description" gorm:"size:256"`
	Tags        pq.StringArray   `json:"tags" gorm:"type:text[]"`
	Status      EntityStatus     `json:"status" gorm:"type:numeric"`
	Permission  EntityPermission `json:"permission" gorm:"type:numeric"`
	BaseUrl     string           `json:"base_url" gorm:"size:256"`
	Credit      int32            `json:"credit" gotm:"type:numeric"`
	SecretKey   string           `json:"secret_key" gotm:"size:64;not null"`
	CreatedAt   time.Time        `json:"created_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
	UpdatedAt   time.Time        `json:"updated_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
}
