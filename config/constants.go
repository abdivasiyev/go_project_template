package config

import (
	"time"
)

// Env keys
const (
	EnvironmentKey         = "environment"
	LogLevelKey            = "log.level"
	JwtSecretKey           = "jwt.secret"
	SecurityMemoryKey      = "security.memory"
	SecurityIterationsKey  = "security.iterations"
	SecurityParallelismKey = "security.parallelism"
	SecuritySaltLengthKey  = "security.salt.length"
	SecurityKeyLengthKey   = "security.key.length"
	BasicAuthUserKey       = "basic.auth.user"
	BasicAuthPasswordKey   = "basic.auth.password"
	PostgresHostKey        = "postgres.host"
	PostgresPortKey        = "postgres.port"
	PostgresUserKey        = "postgres.user"
	PostgresPasswordKey    = "postgres.password"
	PostgresDatabaseKey    = "postgres.database"
	RedisHostKey           = "redis.host"
	RedisPasswordKey       = "redis.password"
	NamespaceKey           = "namespace"
	UploadPathKey          = "upload.path"
	CdnURLKey              = "cdn.url"
	SentryDSNKey           = "sentry.dsn"
	HttpPortKey            = "http.port"
	SpecPath               = "spec.path"
	SpecUrl                = "spec.url"
	SpecTitle              = "spec.title"
	SpecDescription        = "spec.description"
	SmtpHostKey            = "smtp.host"
	SmtpPortKey            = "smtp.port"
	SmtpUsernameKey        = "smtp.username"
	SmtpPasswordKey        = "smtp.password"
)

const (
	Development = "development"
	Test        = "test"
	Staging     = "staging"
	Production  = "production"
)

const (
	DateTimeFormat = "2006-01-02 15:04:05"
	DateFormat     = "2006-01-02"
)

const (
	SuperMinimalCacheTime = 100 * time.Millisecond
	MinimalCacheTime      = 5 * time.Second
	AverageCacheTime      = 30 * time.Second
	MaximalCacheTime      = 1 * time.Minute
	SuperMaximalCacheTime = 10 * time.Minute
)
