// +build windows

package ui

// PPrompt calls prompt, because in windows the lib that hides the typed response
// cant be used
func PPrompt(p string) string {
	return Prompt(p)
}
