package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
    Email string `json:"email"`
    jwt.StandardClaims
}

func GenerateToken(email string) (string, error) {
    expirationTime := time.Now().Add(72 * time.Hour)
    claims := &Claims{
        Email: email,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateToken(tokens string) (*Claims, error) {
    claims := &Claims{}

    token, err := jwt.ParseWithClaims(tokens, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, err
        }
        return nil, err
    }
    if !token.Valid {
        return nil, err
    }

    return claims, nil
}