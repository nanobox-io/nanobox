package models

// ComponentPlan ...
type ComponentPlan struct {
	IPs           []string            `json:"ips"`
	Users         []ComponentPlanUser `json:"users"`
	MountProtocol string              `json:"mount_protocol"`
	Behaviors     []string            `json:"behaviors"`
	DefaultUser   string              `json:"user"`
}

// ComponentPlanUser ...
type ComponentPlanUser struct {
	Username string                 `json:"username"`
	Password string                 `json:"password"`
	Meta     map[string]interface{} `json:"meta"`
}

// BehaviorPresent ...
func (p ComponentPlan) BehaviorPresent(b string) bool {
	for _, behavior := range p.Behaviors {
		if behavior == b {
			return true
		}
	}

	return false
}
