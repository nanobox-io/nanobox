package processors

type DeployConfig struct {
	App     string
	Message string
	Force   bool
}

type TunnelConfig struct {
	App       string
	Port      string
	Container string
}
