package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// 1. Cambiamos el parámetro a string
func GenerateAccessToken(userID string) (string, error) {

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "clave_secreta_super_segura_para_el_mvp_123"
	}

	claims := CustomClaims{
		UserID: userID, // 2. Como ya es string, le quitamos el .String()
		RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "track-my-money-api",
		},
	}

	// 3. Crear el token con el algoritmo de firma HMAC SHA256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Firmar el token usando nuestra clave secreta
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
