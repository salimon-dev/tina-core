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
	Network      string     `json:"network" gorm:"size:32;not null"`
	NetworkId    uuid.UUID  `json:"network_id" gorm:"type:uuid"`
	Username     string     `json:"username" gorm:"size:32;unique;not null"`
	Usage        int32      `json:"usage" gorm:"type:numeric"`
	Status       UserStatus `json:"status" gorm:"type:numeric"`
	Role         UserRole   `json:"role"`
	RegisteredAt time.Time  `json:"registered_at"`
	CreateAt     time.Time  `json:"created_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"type:TIMESTAMP WITH TIME ZONE"`
}

type UserData struct {
	Id           uuid.UUID  `json:"id"`
	Username     string     `json:"username"`
	Status       UserStatus `json:"status"`
	Role         UserRole   `json:"role"`
	RegisteredAt time.Time  `json:"registered_at"`
}
