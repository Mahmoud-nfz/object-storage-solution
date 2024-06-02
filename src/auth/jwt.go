package auth

import (
	"data-storage/src/config"

	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTPayload struct {
	ID     string `json:"id"`
	UserID string `json:"userId"`
	Hash   string `json:"hash"`
	Size   int64  `json:"size,string"`
	Path   string `json:"path"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

func verify(tokenString string) (*JWTPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTPayload{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Unexpected signing method: ", token.Header["alg"])
			return nil, errors.New("unexpected signing method")
		}

		return config.Env.JWTSecret, nil
	})
	if err != nil {
		log.Println("Error parsing JSON: ", err)
		return nil, err
	}

	claims, ok := token.Claims.(*JWTPayload)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func sign(userID, hash string) (string, error) {
	expirationTime := time.Now().Add(config.Env.JWTMaximumAge)

	claims := &JWTPayload{
		UserID: userID,
		Hash:   hash,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(config.Env.JWTSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
