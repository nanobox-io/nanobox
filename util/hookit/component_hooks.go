package hookit

import (
	"fmt"
)


// RunPlanHook runs the plan hook inside of the specified container
func RunPlanHook(container, payload string) (string, error) {
	// run the plan hook
	res, err := Exec(container, "plan", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute plan hook: %s", err.Error())
	}

	return res, nil
}

// RunStartHook runs the start hook inside of the specified container
func RunStartHook(container, payload string) (string, error) {
	// run the start hook
	res, err := Exec(container, "start", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute start hook: %s", err.Error())
	}

	return res, nil
}

// RunUpdateHook runs the update hook inside of the specified container
func RunUpdateHook(container, payload string) (string, error) {
	// run the update hook
	res, err := Exec(container, "update", payload, "info")
	if err != nil {
		return "", fmt.Errorf("failed to execute update hook: %s", err.Error())
	}

	return res, nil
}
