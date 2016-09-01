package steps
	
func Build(name string, complete CompleteCheckFunc, cmd CmdFunc) {
	stepList[name] = step{
		complete: complete,
		cmd:      cmd,
	}
}