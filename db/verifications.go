package db

import (
	"math/rand"
	"salimon/nexus/types"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateRandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := range result {
		letter := letters[rand.Int()%len(letters)]
		result[i] = letter
	}
	return string(result)
}

func VerificationsModel() *gorm.DB {
	return DB.Model(types.Verification{})
}

func FindVerification(query interface{}, args ...interface{}) (*types.Verification, error) {
	var user types.Verification

	result := DB.Model(types.Verification{}).Where(query, args...).Find(&user)

	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func InsertVerification(verification *types.Verification) error {
	result := DB.Model(types.Verification{}).Create(verification)
	return result.Error
}

// gets active verification record based on token and expire time
func GetVerificationRecord(token string) (*types.Verification, error) {
	var verification types.Verification
	result := DB.Model(types.Verification{}).Where("token = ?", token).Find(&verification)

	if result.RowsAffected == 0 {
		return nil, result.Error
	} else {
		return &verification, nil
	}
}

func InsertRegisterEmailVerification(userId uuid.UUID) (*types.Verification, error) {
	result := VerificationsModel().Where("user_id = ? AND domain = ? AND type = ?", userId, types.VerificationDomainRegister, types.VerificationTypeEmail).Delete(nil)
	if result.Error != nil {
		return nil, result.Error
	}
	verification := types.Verification{
		Id:        uuid.New(),
		UserId:    userId,
		Type:      types.VerificationTypeEmail,
		Domain:    types.VerificationDomainRegister,
		Token:     generateRandomString(16),
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	err := InsertVerification(&verification)
	return &verification, err
}

func InsertPasswordResetVerification(userId uuid.UUID) (*types.Verification, error) {
	result := VerificationsModel().Where("user_id = ? AND domain = ? AND type = ?", userId, types.VerificationDomainPasswordReset, types.VerificationTypeEmail).Delete(nil)
	if result.Error != nil {
		return nil, result.Error
	}
	verification := types.Verification{
		Id:        uuid.New(),
		UserId:    userId,
		Type:      types.VerificationTypeEmail,
		Domain:    types.VerificationDomainPasswordReset,
		Token:     generateRandomString(16),
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	err := InsertVerification(&verification)
	return &verification, err
}

func InsertEmailUpdateVerification(userId uuid.UUID) (*types.Verification, error) {
	result := VerificationsModel().Where("user_id = ? AND domain = ? AND type = ?", userId, types.VerificationDomainEmailUpdate, types.VerificationTypeEmail).Delete(nil)
	if result.Error != nil {
		return nil, result.Error
	}
	verification := types.Verification{
		Id:        uuid.New(),
		UserId:    userId,
		Type:      types.VerificationTypeEmail,
		Domain:    types.VerificationDomainEmailUpdate,
		Token:     generateRandomString(16),
		ExpiresAt: time.Now().Add(2 * time.Hour),
	}
	err := InsertVerification(&verification)
	return &verification, err
}

func DeleteVerification(verification *types.Verification) error {
	result := DB.Delete("id = ?", verification.Id)
	return result.Error
}
