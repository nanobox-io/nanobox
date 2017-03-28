package bridge

func Stop() error {
	if runningBridge == nil {
		return nil
	}

	if err := runningBridge.Process.Kill(); err != nil {
		return err
	}

	if err := runningBridge.Wait(); err != nil {
		// it gets a signal but it shows up as an error
		// we dont want that
		return nil
	}

	// if we killed it and released the resources
	// remove running bridge
	runningBridge = nil
	return nil
}
