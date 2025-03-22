package types

import (
	"time"

	"github.com/google/uuid"
)

type UserStatus int8

const (
	UserStatusActive   UserStatus = 1
	UserStatusInActive UserStatus = 2
)

type UserRole int8

const (
	UserRoleKeyMaker  UserRole = 1
	UserRoleAdmin     UserRole = 2
	UserRoleDeveloper UserRole = 3
	UserRoleMember    UserRole = 4
)

type User struct {
	Id           uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey"`
	Username     string     `json:"username" gorm:"size:32;unique;not null"`
	Password     string     `json:"password" gorm:"size:32"`
	InvitationId uuid.UUID  `json:"invitation_id" gorm:"type:uuid"`
	Credit       int32      `json:"credit" gorm:"type:numeric"`
	Role         UserRole   `json:"role" gorm:"type:numeric"`
	Status       UserStatus `json:"status" gorm:"type:numeric"`
	SecretKey    string     `json:"secret_key" gotm:"size:64;not null"`
	RegisteredAt time.Time  `json:"registered_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
}

type PublicUser struct {
	Id           string    `json:"id"`
	Username     string    `json:"username"`
	Credit       int32     `json:"credit"`
	Usage        int32     `json:"usage"`
	Role         string    `json:"role"`
	Status       string    `json:"status"`
	RegisteredAt time.Time `json:"registered_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
