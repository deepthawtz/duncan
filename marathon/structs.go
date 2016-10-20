package marathon

// Group represents a Marathon group
type Group struct {
	Apps []App `json:"apps"`
}

// App represents a Marathon app
type App struct {
	ID        string    `json:"id"`
	Container Container `json:"container"`
}

// Container represents a Marathon container
type Container struct {
	Docker Docker `json:"docker"`
}

// Docker represents marathon docker metadata
type Docker struct {
	Image string `json:"image"`
}

// Deployment represents a Marathon deployment
type Deployment struct {
	ID string `json:"id"`
}

// DeploymentResponse represents a Marathon deployment response
// when updating/creating a new app/group
type DeploymentResponse struct {
	ID string `json:"deploymentId"`
}
