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
	case prodMode.String():
		return zap.NewProduction()
	case devMode.String():
		return zap.NewDevelopment()
	default:
		return nil, fmt.Errorf("invalid running mode: %q", mode)
	}
}
