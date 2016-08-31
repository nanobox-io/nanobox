package hookit

import (
	"fmt"
)

// RunUserHook runs the user hook inside of the specified container
func RunUserHook(container, payload string) (string, error) {
	// run the user hook
	res, err := Exec(container, "user", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute user hook: %s", err.Error())
	}

	return res, nil
}

// RunConfigureHook runs the configure hook inside of the specified container
func RunConfigureHook(container, payload string) (string, error) {
	// run the configure hook
	res, err := Exec(container, "configure", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute configure hook: %s", err.Error())
	}

	return res, nil
}

// RunFetchHook runs the fetch hook inside of the specified container
func RunFetchHook(container, payload string) (string, error) {
	// run the fetch hook
	res, err := Exec(container, "fetch", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute fetch hook: %s", err.Error())
	}

	return res, nil
}

// RunSetupHook runs the setup hook inside of the specified container
func RunSetupHook(container, payload string) (string, error) {
	// run the setup hook
	res, err := Exec(container, "setup", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute setup hook: %s", err.Error())
	}

	return res, nil
}

// RunBoxfileHook runs the boxfile hook inside of the specified container
func RunBoxfileHook(container, payload string) (string, error) {
	// run the boxfile hook
	res, err := Exec(container, "boxfile", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute boxfile hook: %s", err.Error())
	}

	return res, nil
}

// RunPrepareHook runs the prepare hook inside of the specified container
func RunPrepareHook(container, payload string) (string, error) {
	// run the prepare hook
	res, err := Exec(container, "prepare", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute prepare hook: %s", err.Error())
	}

	return res, nil
}

// RunCompileHook runs the compile hook inside of the specified container
func RunCompileHook(container, payload string) (string, error) {
	// run the compile hook
	res, err := Exec(container, "compile", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute compile hook: %s", err.Error())
	}

	return res, nil
}

// RunPackAppHook runs the pack-app hook inside of the specified container
func RunPackAppHook(container, payload string) (string, error) {
	// run the pack-app hook
	res, err := Exec(container, "pack-app", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute pack-app hook: %s", err.Error())
	}

	return res, nil
}

// RunPackBuildHook runs the pack-build hook inside of the specified container
func RunPackBuildHook(container, payload string) (string, error) {
	// run the pack-build hook
	res, err := Exec(container, "pack-build", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute pack-build hook: %s", err.Error())
	}

	return res, nil
}

// RunCleanHook runs the clean hook inside of the specified container
func RunCleanHook(container, payload string) (string, error) {
	// run the clean hook
	res, err := Exec(container, "clean", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute clean hook: %s", err.Error())
	}

	return res, nil
}

// RunPackDeployHook runs the pack-deploy hook inside of the specified container
func RunPackDeployHook(container, payload string) (string, error) {
	// run the pack-deploy hook
	res, err := Exec(container, "pack-deploy", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute pack-deploy hook: %s", err.Error())
	}

	return res, nil
}
