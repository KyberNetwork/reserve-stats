package app

const (
	modeFlag = "mode"
)

//go:generate stringer -type=runningMode -linecomment
type runningMode int

const (
	devMode  runningMode = iota // dev
	prodMode                    // prod
)

var validRunningModes = map[string]struct{}{
	devMode.String(): {},
	prodMode.String():  {},
}
