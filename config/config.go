package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// GetConf :nodoc:
func GetConf() {
	viper.AddConfigPath(".")
	viper.AddConfigPath("./..")
	viper.AddConfigPath("./../..")
	viper.AddConfigPath("./../../..")
	viper.SetConfigName("config")
	viper.SetEnvPrefix("svc")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		log.Warnf("%v", err)
	}
}

// HTTPPort HTTP Port :nodoc:
func HTTPPort() string {
	return viper.GetString("http_port")
}

// Env :nodoc:
func Env() string {
	return viper.GetString("env")
}

// CockroachHost :nodoc:
func MySQLHost() string {
	return viper.GetString("mysql.host")
}

// MySQLDatabase :nodoc:
func MySQLDatabase() string {
	return viper.GetString("mysql.database")
}

// MySQLUsername :nodoc:
func MySQLUsername() string {
	return viper.GetString("mysql.username")
}

// MySQLPassword :nodoc:
func MySQLPassword() string {
	return viper.GetString("mysql.password")
}

// MySQLSSLMode :nodoc:
func MySQLSSLMode() string {
	if viper.IsSet("mysql.sslmode") {
		return viper.GetString("mysql.sslmode")
	}
	return "disable"
}

// MySQLMaxIdleConns :nodoc:
func MySQLMaxIdleConns() int {
	if viper.GetInt("mysql.max_idle_conns") <= 0 {
		return DefaultMySQLMaxIdleConns
	}
	return viper.GetInt("mysql.max_idle_conns")
}

// MySQLMaxOpenConns :nodoc:
func MySQLMaxOpenConns() int {
	if viper.GetInt("mysql.max_open_conns") <= 0 {
		return DefaultMySQLMaxOpenConns
	}
	return viper.GetInt("mysql.max_open_conns")
}

// MySQLConnMaxLifetime :nodoc:
func MySQLConnMaxLifetime() time.Duration {
	if viper.GetInt("mysql.conn_max_lifetime") <= 0 {
		return DefaultMySQLConnMaxLifetime
	}
	return time.Duration(viper.GetInt("mysql.conn_max_lifetime")) * time.Millisecond
}

// MySQLPingInterval :nodoc:
func MySQLPingInterval() time.Duration {
	if viper.GetInt("mysql.ping_interval") <= 0 {
		return DefaultMySQLPingInterval
	}
	return time.Duration(viper.GetInt("mysql.ping_interval")) * time.Millisecond
}

// DatabaseDSN :nodoc:
func DatabaseDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		MySQLUsername(),
		MySQLPassword(),
		MySQLHost(),
		MySQLDatabase())
}

// LogLevel :nodoc:
func LogLevel() string {
	return viper.GetString("log_level")
}
