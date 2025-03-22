package db

import (
	"salimon/tina-core/types"

	"gorm.io/gorm"
)

func UsersModel() *gorm.DB {
	return DB.Model(types.User{})
}

func FindUser(query interface{}, args ...interface{}) (*types.User, error) {
	var user types.User

	result := DB.Model(types.User{}).Where(query, args...).Find(&user)

	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func InsertUser(user *types.User) error {
	result := DB.Model(types.User{}).Create(user)
	return result.Error
}
func UpdateUser(user *types.User) error {
	result := DB.Model(types.User{}).Where("id = ?", user.Id).Updates(user)
	return result.Error
}

func FindUsers(query interface{}, offset int, limit int, args ...interface{}) ([]types.User, error) {
	var users []types.User
	result := DB.Model(types.User{}).Select("*").Where(query, args...).Offset(offset).Limit(limit).Find(users)
	return users, result.Error
}
