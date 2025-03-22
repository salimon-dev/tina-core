package helpers

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthClaims struct {
	Type   string
	UserID uuid.UUID
}

func getSecretKey() string {
	return os.Getenv("SECRET_KEY")
}

func GenerateJwtString(claims jwt.Claims, secretKey string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

func GenerateNexusAccessToken() (string, error) {
	secretKey := getSecretKey()
	claims := jwt.MapClaims{
		"sub":       os.Getenv("ENTITY_ID"),
		"tokenType": "access",
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
	}
	return GenerateJwtString(claims, secretKey)
}

func VerifyJWT(token string) (*AuthClaims, error) {
	secretKey := getSecretKey()
	result, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("uxpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := result.Claims.(jwt.MapClaims); ok && result.Valid {
		result := AuthClaims{}
		sub, ok := claims["sub"].(string)
		if !ok {
			return nil, nil
		}
		uuid, err := uuid.Parse(sub)
		if err != nil {
			return nil, err
		}
		result.UserID = uuid

		tokenType, ok := claims["tokenType"].(string)
		if !ok {
			return nil, nil
		}

		result.Type = tokenType
		return &result, nil
	} else {
		return nil, nil
	}
}
