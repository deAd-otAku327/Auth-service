package tokenizer

import (
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const RefreshCookieName = "refresh_session"

var (
	ErrTokenInvalid = errors.New("access token invalid")
	ErrTokenExpired = errors.New("access token expired")
)

type Tokenizer interface {
	GenerateAccessTokenJWT(userID string) (*string, error)
	VerifyAccessTokenJWT(tokenString string) (jwt.MapClaims, error)
	GenerateRefreshTokenCookie() *http.Cookie
}

type tokenizer struct {
	tokenIssuer        string
	jwtSecret          []byte
	accessTokenExpire  time.Duration
	refreshTokenExpire time.Duration
}

func New(iss, jwtKey string, accessExpire, refreshExpire time.Duration) Tokenizer {
	return &tokenizer{
		tokenIssuer:        iss,
		jwtSecret:          []byte(jwtKey),
		accessTokenExpire:  accessExpire,
		refreshTokenExpire: refreshExpire,
	}
}

func (t *tokenizer) GenerateAccessTokenJWT(userGUID string) (*string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"sub": userGUID,
		"iss": t.tokenIssuer,
		"exp": time.Now().Add(t.accessTokenExpire).Unix(),
		"iat": time.Now().Unix(),
	})

	accessToken, err := claims.SignedString(t.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (t *tokenizer) VerifyAccessTokenJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS512.Alg() {
			return nil, ErrTokenInvalid
		}
		return t.jwtSecret, nil
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

func (t *tokenizer) GenerateRefreshTokenCookie() *http.Cookie {
	refreshToken := base64.StdEncoding.EncodeToString([]byte(uuid.NewString()))
	return &http.Cookie{
		Name:    RefreshCookieName,
		Value:   refreshToken,
		Expires: time.Now().Add(t.refreshTokenExpire),
		Path:    "/api/auth/refresh",
	}
}
