package tokenizer

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenGeneration = errors.New("access token generation failed")
	ErrTokenInvalid    = errors.New("access token invalid")
	ErrTokenExpired    = errors.New("access token expired")
)

type JWTTokenizer interface {
	GenerateAccessToken(userID string) (*string, error)
	VerifyAccessToken(tokenString string) (jwt.MapClaims, error)
}

type tokenizer struct {
	tokenIssuer string
	secretKey   []byte
	tokenExpire time.Duration
}

func New(iss, key string, expire time.Duration) JWTTokenizer {
	return &tokenizer{
		tokenIssuer: iss,
		secretKey:   []byte(key),
		tokenExpire: expire,
	}
}

func (t *tokenizer) GenerateAccessToken(userGUID string) (*string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userGUID,
		"iss": t.tokenIssuer,
		"exp": time.Now().Add(t.tokenExpire).Unix(),
		"iat": time.Now().Unix(),
	})

	accessToken, err := claims.SignedString(t.secretKey)
	if err != nil {
		return nil, ErrTokenGeneration
	}

	return &accessToken, nil
}

func (t *tokenizer) VerifyAccessToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrTokenInvalid
		}
		return t.secretKey, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
