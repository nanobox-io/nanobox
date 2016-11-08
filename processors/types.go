package processors

type DeployConfig struct {
	App      string
	Message  string
	Force    bool
}

type ConsoleConfig struct {
	App      string
	Host     string
}

type TunnelConfig struct {
	App       string
	Port      string
	Container string
}
