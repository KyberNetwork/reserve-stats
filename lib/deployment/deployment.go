package deployment

/**
Deployment represents a collection of Kyber Network smart contracts deployments.
There might be multiple deployments in same network used for different purpose.
Deployment is a separated concept from running mode to allow developers to run any deployment in debug mode.
*/

//Deployment is a enum type for checking valid DeploymentMode
//go:generate stringer -type=Deployment -linecomment
type Deployment int

const (
	//Production is production mode for deployment
	Production Deployment = iota //production
	//Staging is staging mode for deployment
	Staging //staging
	//Ropsten is ropsten mode for deployment
	Ropsten //ropsten
)
