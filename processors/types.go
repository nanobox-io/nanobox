package processors

type DeployConfig struct {
	App      string
	Message  string
	Force    bool
	Endpoint string
}

type ConsoleConfig struct {
	App      string
	Host     string
	Endpoint string
}

type TunnelConfig struct {
	App       string
	Port      string
	Container string
	Endpoint  string
}
