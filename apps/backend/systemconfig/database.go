package systemconfig

import (
	"strings"

	"resuming/env"
)

var DatabaseType string
var DatabaseUsername string
var DatabaseHost string
var DatabaseName string
var DatabasePassword string
var DatabasePort string
var DatabaseDSN string

func init() {
	DatabaseType = env.GetEnv("DATABASE_TYPE")
	DatabaseUsername = env.GetEnv("DATABASE_USERNAME")
	DatabaseHost = env.GetEnv("DATABASE_HOST")
	DatabaseName = env.GetEnv("DATABASE_NAME")
	DatabasePassword = env.GetEnv("DATABASE_PASSWORD")
	DatabasePort = env.GetEnv("DATABASE_PORT")

	dsn_details := []string{
		"host=" + DatabaseHost,
		"user=" + DatabaseUsername,
		"password=" + DatabasePassword,
		"dbname=" + DatabaseName,
		"port=" + DatabasePort,
		"ssl_mode=disabled",
		"TimeZone=Asia/Shanghai",
	}
	dsn := strings.Join(dsn_details, " ")

	DatabaseDSN = dsn
}
