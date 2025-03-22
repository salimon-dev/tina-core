package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"salimon/nexus/types"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthClaims struct {
	Type   string
	UserID uuid.UUID
}

func getSecretKey() string {
	return os.Getenv("JWT_SECRET")
}

func generateJWTString(claims jwt.Claims, secretKey string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	return tokenString, err
}

func generateAccessToken(userId string, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userId,
		"tokenType": "access",
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Hour * 6)),
	}
	return generateJWTString(claims, secretKey)
}

func generateRefreshToken(userId string, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"sub":       userId,
		"tokenType": "refresh",
		"exp":       jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 14)),
	}
	return generateJWTString(claims, secretKey)
}

func GenerateNexusJwt(user *types.User) (*string, *string, error) {
	secretKey := getSecretKey()
	return GenerateJWT(user, string(secretKey))
}

func GenerateJWT(user *types.User, secretKey string) (*string, *string, error) {
	accessToken, err := generateAccessToken(user.Id.String(), secretKey)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := generateRefreshToken(user.Id.String(), secretKey)
	if err != nil {
		return nil, nil, err
	}
	return &accessToken, &refreshToken, nil
}

func GenerateJWTFromEntity(entityId uuid.UUID, secretKey string) (*string, *string, error) {
	accessToken, err := generateAccessToken(entityId.String(), secretKey)
	if err != nil {
		return nil, nil, err
	}
	refreshToken, err := generateRefreshToken(entityId.String(), secretKey)
	if err != nil {
		return nil, nil, err
	}
	return &accessToken, &refreshToken, nil
}

func VerifyNexusJWT(token string) (*AuthClaims, error) {
	secretKey := getSecretKey()
	return VerifyJWT(token, secretKey)
}

func VerifyJWT(token string, secretKey string) (*AuthClaims, error) {
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

func ExtractSubFromJwt(token string) (*uuid.UUID, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid jwt token")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims jwt.Claims
	err = json.Unmarshal(payload, &claims)
	if err != nil {
		return nil, err
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, err
	}
	if sub == "" {
		return nil, errors.New("no subject found in jwt token")
	}
	subId, err := uuid.FromBytes([]byte(sub))
	if err != nil {
		return nil, err
	}
	return &subId, nil
}
