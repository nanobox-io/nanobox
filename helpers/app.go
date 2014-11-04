package helpers

// DetermineAppType takes an app's 'type' and returns a '*' if it is 'tinker'
func DetermineAppType(t bool) string {

	// 'tinker' app
	if t {
		return "*"
	}

	// 'production' app (default)
	return ""
}

// DetermineAppFlation takes an app's 'flation' and returns a 'friendly' version
func DetermineAppFlation(f string) string {
	switch f {
	case "inflating", "inflated":
		return "Awake"
	case "deflating", "deflated":
		return "Asleep"
	default:
		return f
	}
}

// DetermineAppStatus takes an app's 'status' and returns a corresponding color
// to represent that status
func DetermineAppStatus(s, f string) string {
	switch s {

	//
	case "initialized", "created", "active":

		//
		if f == "deflated" {
			return "[yellow]"
		} else {
			switch s {

			//
			case "initialized", "created":
				return "[blue]"

			//
			case "active":
				return "[green]"
			}
		}

	//
	case "uninitialized", "inactive", "defunct", "hybernated":
		return "[red]"
	}

	return ""
}
