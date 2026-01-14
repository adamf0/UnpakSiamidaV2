package helper

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret")

type JWTService struct {
	Secret []byte
}

func GenerateToken(sid string) (string, string, error) {
	// Access token (15 menit)
	accessClaims := jwt.MapClaims{
		"sid": sid,
		"exp": time.Now().Add(1 * 24 * time.Hour).Unix(),
	}
	accessJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenStr, err := accessJWT.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	// Refresh token (7 hari)
	refreshClaims := jwt.MapClaims{
		"sid": sid,
		"exp": time.Now().Add(14 * 24 * time.Hour).Unix(),
	}
	refreshJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenStr, err := refreshJWT.SignedString(jwtSecret)
	if err != nil {
		return "", "", err
	}

	return accessTokenStr, refreshTokenStr, nil
}

// func NewJWTService() *JWTService {
// 	return &JWTService{
// 		Secret: jwtSecret,
// 	}
// }

// func (s *JWTService) ValidateToken(authHeader string) (jwt.MapClaims, error) {
// 	if authHeader == "" {
// 		return nil, errors.New("authorization header missing")
// 	}

// 	parts := strings.Split(authHeader, " ")
// 	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
// 		return nil, errors.New("authorization header format must be Bearer {token}")
// 	}

// 	tokenStr := parts[1]

// 	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
// 		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("invalid signing method")
// 		}
// 		return s.Secret, nil
// 	})

// 	if err != nil || !token.Valid {
// 		if err != nil {
// 			return nil, errors.New("invalid token: " + err.Error())
// 		}
// 		return nil, errors.New("invalid token")
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok {
// 		return nil, errors.New("invalid token claims")
// 	}

// 	if exp, ok := claims["exp"].(float64); ok && int64(exp) < time.Now().Unix() {
// 		return nil, errors.New("token expired")
// 	}

// 	return claims, nil
// }