package steps

func Build(name string, private bool, complete CompleteCheckFunc, cmd CmdFunc) {
	stepList[name] = step{
		private: private,
		complete: complete,
		cmd:      cmd,
	}
}
