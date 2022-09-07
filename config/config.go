package config

import (
	"github.com/spf13/viper"
	"go.uber.org/fx"
	"log"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var Module = fx.Provide(New)

type Config interface {
	Get(key string) interface{}
	GetBool(key string) bool
	GetFloat64(key string) float64
	GetInt(key string) int
	GetUInt32(key string) uint32
	GetUInt8(key string) uint8
	GetIntSlice(key string) []int
	GetString(key string) string
	GetStringMap(key string) map[string]interface{}
	GetStringMapString(key string) map[string]string
	UnmarshalKey(key string, val interface{}) error
	GetStringSlice(key string) []string
	GetDuration(key string) time.Duration
}

type config struct {
	cfg *viper.Viper
}

func New() (Config, error) {
	cfg := viper.New()

	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")
	cfg.AddConfigPath("./config")

	cfg.AddConfigPath(RootDir() + "/config")

	if err := cfg.ReadInConfig(); err != nil {
		log.Printf("could not read config file: %v", err)
	}

	cfg.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.AutomaticEnv()

	cfg.WatchConfig()

	return &config{cfg: cfg}, nil
}

func (c *config) Get(key string) interface{} {
	return c.cfg.Get(key)
}

func (c *config) GetBool(key string) bool {
	return c.cfg.GetBool(key)
}

func (c *config) GetFloat64(key string) float64 {
	return c.cfg.GetFloat64(key)
}

func (c *config) GetInt(key string) int {
	return c.cfg.GetInt(key)
}

func (c *config) GetUInt32(key string) uint32 {
	return c.cfg.GetUint32(key)
}

func (c *config) GetUInt8(key string) uint8 {
	return uint8(c.cfg.GetUint(key))
}

func (c *config) GetIntSlice(key string) []int {
	return c.cfg.GetIntSlice(key)
}

func (c *config) GetString(key string) string {
	return c.cfg.GetString(key)
}

func (c *config) GetStringSlice(key string) []string {
	return c.cfg.GetStringSlice(key)
}

func (c *config) GetStringMap(key string) map[string]interface{} {
	return c.cfg.GetStringMap(key)
}
func (c *config) GetStringMapString(key string) map[string]string {
	return c.cfg.GetStringMapString(key)
}

func (c *config) UnmarshalKey(key string, val interface{}) error {
	return c.cfg.UnmarshalKey(key, &val)
}

func (c *config) GetDuration(key string) time.Duration {
	return c.cfg.GetDuration(key)
}

func RootDir() string {
	_, b, _, _ := runtime.Caller(0)
	d := path.Join(path.Dir(b))
	return filepath.Dir(d)
}
