package db

import (
	"salimon/nexus/types"

	"gorm.io/gorm"
)

func EntityModel() *gorm.DB {
	return DB.Model(types.Entity{})
}

func FindEntity(query interface{}, args ...interface{}) (*types.Entity, error) {
	var entity types.Entity

	result := DB.Model(types.Entity{}).Where(query, args...).Find(&entity)

	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &entity, nil
	}
}

func InsertEntity(entity *types.Entity) error {
	result := EntityModel().Create(entity)
	return result.Error
}
