package authentication

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"films-api.rdelgado.es/src/internals/models"
	"github.com/golang-jwt/jwt/v5"
)

type JwtToken struct {
	SecretJwt []byte
}

func (t *JwtToken) ExtractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", models.ErrInvalidAuthHeader
	}

	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", models.ErrInvalidAuthHeader
	}

	return parts[1], nil
}

func (t *JwtToken) VerifyToken(tokenString string) (int, error) {

	// TODO: Pasar el logger aqui
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return t.SecretJwt, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, models.ErrTokenExpired
		} else {
			return 0, err
		}
	}

	// Check if token has key "user_id"
	id, ok := token.Claims.(jwt.MapClaims)["user_id"]
	if !ok {
		return 0, models.ErrInvalidToken
	}

	// Check user_id value is a number
	user_id, ok := id.(float64)
	if !ok {
		return 0, models.ErrInvalidToken
	}

	return int(user_id), nil
}

func (t *JwtToken) CreateToken(id int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	// Crear un nuevo token con los claims y el método de firma
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Firmar y obtener el token codificado como string
	tokenString, err := token.SignedString(t.SecretJwt)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
