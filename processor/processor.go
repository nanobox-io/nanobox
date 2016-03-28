package processor

type (
	ProcessConfig struct {
		DevMode bool
		Verbose bool
		Background bool
		Force bool
		Meta map[string]string 
	}

	ProcessBuilder func(ProcessConfig) (Processor, error)

	Processor interface {
		Process() error
		Results() ProcessConfig
	}
)

var (
	processors = map[string]ProcessBuilder{}
)

func Register(name string, sb ProcessBuilder) {
	processors[name] = sb
}

func Build(name string, pc ProcessConfig) Processor, error {
	proc, ok := processors[name]
	if !ok {
		return nil, fmt.Errorf("Invalid Processor %s", name)
	}
	return proc(pc)
}

func Run(name string, pc ProcessConfig) error {
	proc, err := Build(name, pc)
	if err != nil {
		return err
	}
	return proc.Process()
}