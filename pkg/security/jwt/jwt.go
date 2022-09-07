package jwt

import (
	"fmt"
	"time"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/pkg/errors"

	jwtGo "github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken = errors.New("invalid jwt token")
	ErrExpiredToken = errors.New("expired jwt token")
)

type Jwt struct {
	secretKey string
}

const minSecretKeySize = 32

func NewJwt(secretKey string) (*Jwt, error) {
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}
	return &Jwt{secretKey}, nil
}

func (j *Jwt) CreateToken(user models.GetUserResponse, duration time.Duration, isRefreshToken bool) (string, error) {
	payload, err := NewPayload(user, duration, isRefreshToken)
	if err != nil {
		return "", err
	}

	jwtToken := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(j.secretKey))
}

func (j *Jwt) VerifyToken(token string) (*Payload, error) {
	keyFunc := func(token *jwtGo.Token) (interface{}, error) {
		_, ok := token.Method.(*jwtGo.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(j.secretKey), nil
	}

	jwtToken, err := jwtGo.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwtGo.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	p, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return p, nil
}
