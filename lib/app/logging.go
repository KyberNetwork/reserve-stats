package app

import (
	"fmt"
	"github.com/urfave/cli"
	"go.uber.org/zap"
)

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
