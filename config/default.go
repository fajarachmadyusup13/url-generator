package config

import "time"

const (
	// RetryAttempts :nodoc:
	RetryAttempts = 5
	// DefaultCacheLongerTTL in milliseconds
	DefaultCacheLongerTTL = 6 * time.Hour
	// DefaultMySQLMaxIdleConns min connection pool
	DefaultMySQLMaxIdleConns = 2
	// DefaultMySQLMaxOpenConns max connection pool
	DefaultMySQLMaxOpenConns = 5
	// DefaultMySQLConnMaxLifetime :nodoc:
	DefaultMySQLConnMaxLifetime = 1 * time.Hour
	// DefaultMySQLPingInterval :nodoc:
	DefaultMySQLPingInterval = 1 * time.Second
)
