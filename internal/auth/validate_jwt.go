package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    token, err := jwt.ParseWithClaims(
        tokenString,
        &jwt.RegisteredClaims{},
        func(token *jwt.Token) (interface{}, error) {
            return []byte(tokenSecret), nil
        },
    )
    if err != nil {
        return uuid.Nil, err
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok || !token.Valid {
        return uuid.Nil, errors.New("invalid token")
    }

    userID, err := uuid.Parse(claims.Subject)
    if err != nil {
        return uuid.Nil, err
    }

    return userID, nil
}
