package db

import (
	"salimon/nexus/types"

	"gorm.io/gorm"
)

func InvitationsModel() *gorm.DB {
	return DB.Model(types.Invitation{})
}

func FindInvitation(query interface{}, args ...interface{}) (*types.Invitation, error) {
	var invitation types.Invitation

	result := DB.Model(types.Invitation{}).Where(query, args...).Find(&invitation)

	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &invitation, nil
	}
}

func InsertInvitation(invitation *types.Invitation) error {
	result := DB.Model(types.Invitation{}).Create(invitation)
	return result.Error
}
func UpdateInvitation(invitation *types.Invitation) error {
	result := DB.Model(types.Invitation{}).Where("id = ?", invitation.Id).Updates(invitation)
	return result.Error
}

func FindInvitations(query interface{}, offset int, limit int, args ...interface{}) ([]types.Invitation, error) {
	var invitations []types.Invitation
	result := DB.Model(types.Invitation{}).Select("*").Where(query, args...).Offset(offset).Limit(limit).Find(invitations)
	return invitations, result.Error
}
