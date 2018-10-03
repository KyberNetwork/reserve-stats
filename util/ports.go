package util

import (
	"fmt"
)

// HTTPPort define custom type for port
type HTTPPort int

const (
	// IPLocatorPort port number of IpLocator service
	IPLocatorPort HTTPPort = 8001
)

// GinPort return gin port input format
func (port HTTPPort) GinPort() string {
	return fmt.Sprintf(":%d", port)
}
