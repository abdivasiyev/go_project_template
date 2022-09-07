package security

import (
	"github.com/abdivasiyev/project_template/config"
	"github.com/abdivasiyev/project_template/internal/models"
	"go.uber.org/fx"
)

var Module = fx.Provide(New)

type Handler interface {
	GenerateToken(user models.GetUserResponse) (accessToken string, refreshToken string, err error)
	VerifyToken(token string, isRefreshToken bool) (user models.GetUserResponse, err error)
	GenerateHash(plainText string) (hashedText string, err error)
	CompareHash(plainText string, hashedText string) (valid bool, err error)
	Md5Sum(value any) (string, error)
}

type handler struct {
	jwtSecret   string
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

type Params struct {
	fx.In
	Config config.Config
}

func New(params Params) Handler {
	return &handler{
		jwtSecret:   params.Config.GetString(config.JwtSecretKey),
		memory:      params.Config.GetUInt32(config.SecurityMemoryKey),
		iterations:  params.Config.GetUInt32(config.SecurityIterationsKey),
		parallelism: params.Config.GetUInt8(config.SecurityParallelismKey),
		saltLength:  params.Config.GetUInt32(config.SecuritySaltLengthKey),
		keyLength:   params.Config.GetUInt32(config.SecurityKeyLengthKey),
	}
}
