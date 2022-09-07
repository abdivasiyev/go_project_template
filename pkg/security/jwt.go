package security

import (
	"time"

	"github.com/abdivasiyev/project_template/internal/models"
	"github.com/abdivasiyev/project_template/pkg/security/jwt"
)

func (p *handler) GenerateToken(user models.GetUserResponse) (string, string, error) {
	j, err := jwt.NewJwt(p.jwtSecret)

	if err != nil {
		return "", "", err
	}

	accessToken, err := j.CreateToken(user, 24*time.Hour, false)

	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.CreateToken(user, 2*24*time.Hour, true)

	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (p *handler) VerifyToken(token string, isRefreshToken bool) (models.GetUserResponse, error) {
	j, err := jwt.NewJwt(p.jwtSecret)

	if err != nil {
		return models.GetUserResponse{}, err
	}

	payload, err := j.VerifyToken(token)

	if err != nil {
		return models.GetUserResponse{}, err
	}

	if payload.IsRefreshToken != isRefreshToken {
		return models.GetUserResponse{}, jwt.ErrInvalidToken
	}

	return payload.User, nil
}
