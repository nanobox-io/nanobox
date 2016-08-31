package processors

type DeployConfig struct {
	App     string
	Message string
}

type TunnelConfig struct {
	App       string
	Port      string
	Container string
}
