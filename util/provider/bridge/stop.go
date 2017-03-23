package bridge

func Stop() error {
	if runningBridge == nil {
		return nil
	}

	if err := runningBridge.Process.Kill(); err != nil {
		return err
	}

	if err := runningBridge.Wait(); err != nil {
		return err
	}

	// if we killed it and released the resources 
	// remove running bridge
	runningBridge = nil
	return nil
}
