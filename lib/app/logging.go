package app

import (
	"fmt"

	"github.com/TheZeroSlave/zapsentry"
	"github.com/getsentry/sentry-go"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "error"
	fatalLevel = "fatal"

	sentryDSNFlag      = "sentry-dsn"
	sentryLevelFlag    = "sentry-lv"
	defaultSentryLevel = errorLevel
)

// NewSentryFlags returns flags to init sentry client
func NewSentryFlags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   sentryDSNFlag,
			EnvVar: "SENTRY_DSN",
			Usage:  "dsn for sentry client",
		}, cli.StringFlag{
			Name:   sentryLevelFlag,
			EnvVar: "SENTRY_LEVEL",
			Usage:  "log level report message to sentry (info, error, warn, fatal)",
			Value:  defaultSentryLevel,
		},
	}
}

type syncer interface {
	Sync() error
}

// NewFlusher creates a new syncer from given syncer that log a error message if failed to sync.
func NewFlusher(s syncer) func() {
	return func() {
		// ignore the error as the sync function will always fail in Linux
		// https://github.com/uber-go/zap/issues/370
		_ = s.Sync()
	}
}

// NewLogger creates a new logger instance.
// The type of logger instance will be different with different application running modes.
func NewLogger(c *cli.Context) (*zap.Logger, error) {
	mode := c.GlobalString(modeFlag)
	switch mode {
	case productionMode:
		return zap.NewProduction()
	case developmentMode:
		return zap.NewDevelopment()
	default:
		return nil, fmt.Errorf("invalid running mode: %q", mode)
	}
}

// NewSugaredLogger creates a new sugared logger and a flush function. The flush function should be
// called by consumer before quitting application.
// This function should be use most of the time unless
// the application requires extensive performance, in this case use NewLogger.
func NewSugaredLogger(c *cli.Context) (*zap.SugaredLogger, func(), error) {
	logger, err := NewLogger(c)
	if err != nil {
		return nil, nil, err
	}
	// init sentry if flag dsn exists
	if len(c.String(sentryDSNFlag)) != 0 {
		sentryClient, err := sentry.NewClient(
			sentry.ClientOptions{
				Dsn: c.String(sentryDSNFlag),
			},
		)
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to init sentry client")
		}

		cfg := zapsentry.Configuration{
			DisableStacktrace: false,
		}
		switch c.String(sentryLevelFlag) {
		case infoLevel:
			cfg.Level = zapcore.InfoLevel
		case warnLevel:
			cfg.Level = zapcore.WarnLevel
		case errorLevel:
			cfg.Level = zapcore.ErrorLevel
		case fatalLevel:
			cfg.Level = zapcore.FatalLevel
		default:
			return nil, nil, errors.Errorf("invalid log level %v", c.String(sentryLevelFlag))
		}

		core, err := zapsentry.NewCore(cfg, zapsentry.NewSentryClientFromClient(sentryClient))
		if err != nil {
			return nil, nil, errors.Wrap(err, "failed to init zap sentry")
		}
		// attach to logger core
		logger = zapsentry.AttachCoreToLogger(core, logger)
	}
	sugar := logger.Sugar()
	return sugar, NewFlusher(logger), nil
}
