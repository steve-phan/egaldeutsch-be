package auth

import (
	"egaldeutsch-be/internal/config"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserId string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func CreateAccessToken(userID string, role string, jwtConfig config.JwtConfig) (string, error) {
	now := time.Now()
	expiresAt := now.Add(time.Duration(jwtConfig.ExpirationHours) * time.Hour)
	claims := Claims{
		UserId: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    jwtConfig.Issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtConfig.SecretKey)) // Note: convert to []byte for HMAC (SigningMethodHS256)

}

func ParseToken(tokenString string, jwtConfig config.JwtConfig) (*Claims, error) {
	if jwtConfig.SecretKey == "" {
		return nil, errors.New("jwt secret key is empty")
	}

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtConfig.SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if cliams, ok := token.Claims.(*Claims); ok && token.Valid {
		return cliams, nil
	} else {
		return nil, errors.New("invalid token claims")
	}
}
