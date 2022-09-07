package jwt

import (
	"time"

	"github.com/abdivasiyev/project_template/internal/models"
	jwtGo "github.com/golang-jwt/jwt"

	"github.com/google/uuid"
)

type Payload struct {
	jwtGo.StandardClaims
	IsRefreshToken bool                   `json:"is_refresh_token"`
	User           models.GetUserResponse `json:"user"`
}

func NewPayload(user models.GetUserResponse, duration time.Duration, isRefreshToken bool) (*Payload, error) {
	issuedAt := time.Now().UTC()
	expiresAt := issuedAt.Add(duration).Unix()

	return &Payload{
		User: models.GetUserResponse{
			ID:       user.ID,
			Username: user.Username,
			Company: models.GetCompanyResponse{
				ID: user.Company.ID,
			},
		},
		IsRefreshToken: isRefreshToken,
		StandardClaims: jwtGo.StandardClaims{
			ExpiresAt: expiresAt,
			Id:        uuid.New().String(),
			IssuedAt:  issuedAt.Unix(),
		},
	}, nil
}

func (p *Payload) Valid() error {
	if time.Unix(p.ExpiresAt, 0).Before(time.Now().UTC()) {
		return ErrExpiredToken
	}

	return nil
}
