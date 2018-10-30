package deployment

var dpl Deployment

// GetDeployment returns current configured deployment.
func GetDeployment() Deployment {
	return dpl
}

// GetDeployment sets current configured deployment to the one in parameter.
func SetDeployment(newDpl Deployment) {
	dpl = newDpl
}

/**
Deployment represents a collection of Kyber Network smart contracts deployments.
There might be multiple deployments in same network used for different purpose.
Deployment is a separated concept from running mode to allow developers to run any deployment in debug mode.
 */
//go:generate stringer -type=deployment -linecomment
type Deployment int

const (
	ProdDeployment    Deployment = iota + 1 // production
	StagingDeployment                       // staging
)

var validDeployments = map[string]struct{}{
	ProdDeployment.String():    {},
	StagingDeployment.String(): {},
}
