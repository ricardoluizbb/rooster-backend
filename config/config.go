package config

import "github.com/spf13/viper"

const (
	databaseDSN = "DATABASE_DSN"
	httpPort    = "HTTP_PORT"
	magicLink   = "MAGIC_LINK_URL"
)

func DSN() string {
	return viper.GetString(databaseDSN)
}

func HttpPort() string {
	return viper.GetString(httpPort)
}

func MagicLinkUrl() string {
	return viper.GetString(magicLink)
}
