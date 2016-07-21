// +build !windows

package vbox

// DetectVBoxManageCmd tries to find VBoxManage 
func DetectVBoxManageCmd() string {
	return detectVBoxManageCmdInPath()
}
