package jwt

import (
	"fmt"
	"gateway/internal/jwtutil"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWT(tokenStr string) (string, error) {
	pubKey, err := jwtutil.GetPublicKey()
	if err != nil {
		return "", fmt.Errorf("public key not loaded: %w", err)
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return pubKey, nil
	})
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid claims")
	}

	userID, ok := claims["user_id"]
	if !ok {
		return "", fmt.Errorf("user_id not found in token")
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return "", fmt.Errorf("user_id is not a string")
	}

	return userIDStr, nil
}
