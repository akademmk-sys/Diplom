package db

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = []byte("chocalte_starfish_and_hotdog_flavored_water")

func CreateToken(pass string) (string, error) {
	res := sha256.Sum256([]byte(pass))
	hash := hex.EncodeToString(res[:])
	claims := jwt.MapClaims{
		"hash": hash,
		"exp":  time.Now().Add(8 * time.Hour).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidToken(token string) (bool, error) {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
		return secret, nil
	})
	if err != nil {
		return false, err
	}
	if !jwtToken.Valid {
		return false, nil
	}
	claims, ok := jwtToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("failed to type assertion of jwt.MapClaims")
	}
	envPass := os.Getenv("TODO_PASSWORD")
	res := sha256.Sum256([]byte(envPass))
	hashEnvPass := hex.EncodeToString(res[:])
	hashRaw := claims["hash"]
	hash, ok := hashRaw.(string)
	if !ok {
		return false, errors.New("invalid token")
	}
	if hashEnvPass != hash {
		return false, errors.New("invalid token")
	}
	return true, nil
}
