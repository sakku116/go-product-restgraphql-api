package jwt_util

import (
	"backend/domain/dto"
	"backend/domain/enum"
	"backend/domain/model"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func GenerateJwtToken(user *model.UserModel, secretKey string, expMins int, tokenId *string) (string, error) {
	JWT_SIGNATURE_KEY := []byte(secretKey)

	claims := jwt.MapClaims{
		"sub":      user.UUID,
		"username": user.Username,
		"role":     user.Role.String(),
		"email":    user.Email,
		"exp":      time.Now().Add(time.Minute * time.Duration(expMins)).Unix(),
	}

	if tokenId != nil {
		claims["token_id"] = *tokenId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(JWT_SIGNATURE_KEY)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func ValidateJWT(tokenString string, secretKey string) (*dto.CurrentUser, error) {
	var JWT_SIGNATURE_KEY = []byte(secretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("signing method invalid")
		}

		return JWT_SIGNATURE_KEY, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	sub, _ := claims["sub"].(string)
	username, _ := claims["username"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)
	normalizedRole := enum.UserRole(role)
	isRoleValid := normalizedRole.IsValid()
	if !isRoleValid {
		return nil, errors.New("invalid role")
	}

	return &dto.CurrentUser{
		UUID:     sub,
		Username: username,
		Role:     normalizedRole,
		Email:    email,
	}, nil
}
