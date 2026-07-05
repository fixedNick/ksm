package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	HttpConfig
	SQLConfig
}
type HttpConfig struct{ Port int }
type SQLConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DbName   string
}

func (s SQLConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		s.User, s.Password, s.Host, s.Port, s.DbName)
}

func MustLoadConfig() *Config {
	cfg := &Config{}
	cfg.HttpConfig.Port = mustLookupEnvInt("KSM_CHAT_HTTP_PORT")
	cfg.SQLConfig.User = mustLookupEnv("KSM_CHAT_MYSQL_USER")
	cfg.SQLConfig.Password = mustLookupEnv("KSM_CHAT_MYSQL_PASSWORD")
	cfg.SQLConfig.Port = mustLookupEnvInt("KSM_CHAT_MYSQL_PORT")
	cfg.SQLConfig.DbName = mustLookupEnv("KSM_CHAT_MYSQL_DBNAME")
	cfg.SQLConfig.Host = mustLookupEnv("KSM_CHAT_MYSQL_HOST")
	return cfg
}
func mustLookupEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if ok && v != "" {
		return v
	}
	panic("missing env var: " + key)
}
func mustLookupEnvInt(key string) int {
	v := mustLookupEnv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		panic("missing env var: " + key)
	}
	return i
}
