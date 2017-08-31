package processors

type DeployConfig struct {
	App     string
	Message string
	Force   bool
}

type ConsoleConfig struct {
	App  string
	Host string
}
