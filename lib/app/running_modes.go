package app

const (
	modeFlag = "mode"

	developmentMode = "dev"
	productionMode  = "prod"
)

var validRunningModes = map[string]struct{}{
	developmentMode: {},
	productionMode:  {},
}
