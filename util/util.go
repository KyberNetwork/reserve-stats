package util

import (
	"path"
	"path/filepath"
	"runtime"
)

// CurrentDir returns current directory of the caller.
func CurrentDir() string {
	_, current, _, _ := runtime.Caller(1)
	return filepath.Join(path.Dir(current))
}
