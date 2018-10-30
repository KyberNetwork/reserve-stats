package deployment

/**
Deployment represents a collection of Kyber Network smart contracts deployments.
There might be multiple deployments in same network used for different purpose.
Deployment is a separated concept from running mode to allow developers to run any deployment in debug mode.
*/
const (
	Production = "production"
	Staging    = "staging"
)

// TODO: consider making this to an enum type
