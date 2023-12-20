package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/fajarachmadyusup13/url-generator/config"
	"github.com/jpillora/backoff"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

var (
	MySQL *gorm.DB

	StopTickerCh chan bool

	sqlRegexp = regexp.MustCompile(`(\$\d+)|\?`)
)

func InitializeMySQLConn() {
	conn, err := openMySQLConn(config.DatabaseDSN())
	if err != nil {
		logrus.WithField("databaseDSN", config.DatabaseDSN()).Fatal("failed to connect MySQL", err)
	}

	MySQL = conn
	StopTickerCh = make(chan bool)

	// go checkConnection(time.NewTicker(config.DefaultMySQLPingInterval))

	MySQL.Logger = NewGormCustomLogger()

	switch config.LogLevel() {
	case "error":
		MySQL.Logger = MySQL.Logger.LogMode(gormLogger.Error)
	case "warn":
		MySQL.Logger = MySQL.Logger.LogMode(gormLogger.Warn)
	default:
		MySQL.Logger = MySQL.Logger.LogMode(gormLogger.Info)

	}

	log.Info("Connection to Cockroach Server success...")
}

func openMySQLConn(dsn string) (*gorm.DB, error) {
	dialector := mysql.Open(dsn)
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	conn, err := db.DB()
	if err != nil {
		logrus.Fatal(err)
	}

	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(1 * time.Hour)
	conn.SetMaxOpenConns(5)

	return db, nil
}

func checkConnection(ticker *time.Ticker) {
	fmt.Println("MASHOOK: ")

	for {
		select {
		case <-StopTickerCh:
			fmt.Println("STOPPED: ")

			ticker.Stop()
			return
		case <-ticker.C:
			if _, err := MySQL.DB(); err != nil {
				// reconnectMySQLConn()
			}
		}
	}
}

func reconnectMySQLConn() {
	b := backoff.Backoff{
		Factor: 2,
		Jitter: true,
		Min:    100 * time.Millisecond,
		Max:    1 * time.Second,
	}

	for b.Attempt() < config.RetryAttempts {
		conn, err := openMySQLConn(config.DatabaseDSN())
		if err != nil {
			logrus.WithField("databaseDSN", config.DatabaseDSN()).Fatal("failed to connect MySQL", err)
		}

		if conn != nil {
			MySQL = conn
			break
		}
		time.Sleep(b.Duration())
	}

	if b.Attempt() >= config.RetryAttempts {
		log.Fatal("maximum retry to connect database")
	}
	b.Reset()
}

// GormCustomLogger override gorm logger
type GormCustomLogger struct {
	gormLogger.Config
}

// NewGormCustomLogger :nodoc:
func NewGormCustomLogger() *GormCustomLogger {
	return &GormCustomLogger{
		Config: gormLogger.Config{
			LogLevel: gormLogger.Info,
		},
	}
}

// LogMode :nodoc:
func (g *GormCustomLogger) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	g.LogLevel = level
	return g
}

// Info :nodoc:
func (g *GormCustomLogger) Info(ctx context.Context, message string, values ...interface{}) {
	if g.LogLevel >= gormLogger.Info {
		logrus.WithFields(logrus.Fields{"data": values}).Info(message)
	}
}

// Warn :nodoc:
func (g *GormCustomLogger) Warn(ctx context.Context, message string, values ...interface{}) {
	if g.LogLevel >= gormLogger.Warn {
		logrus.WithFields(logrus.Fields{"data": values}).Warn(message)
	}

}

// Error :nodoc:
func (g *GormCustomLogger) Error(ctx context.Context, message string, values ...interface{}) {
	if g.LogLevel >= gormLogger.Error {
		logrus.WithFields(logrus.Fields{"data": values}).Error(message)
	}
}

// Trace :nodoc:
func (g *GormCustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc()
	if g.LogLevel <= 0 {
		return
	}

	elapsed := time.Since(begin)
	logger := logrus.WithFields(logrus.Fields{
		"took": elapsed,
	})

	sqlLog := sqlRegexp.ReplaceAllString(sql, "%v")
	if rows >= 0 {
		logger.WithField("rows", rows)
	} else {
		logger.WithField("rows", "-")
	}

	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound) && g.LogLevel >= gormLogger.Error:
		logger.WithField("sql", sqlLog).Error(err)
	case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && g.LogLevel >= gormLogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", g.SlowThreshold)
		logger.WithField("sql", sqlLog).Warn(slowLog)
	case g.LogLevel >= gormLogger.Info:
		logger.Info(sqlLog)

	}
}
