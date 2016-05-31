package models

type Plan struct {
	IPs           []string `json:"ips"`
	Users         []User   `json:"users"`
	MountProtocol string   `json:"mount_protocol"`
	Behaviors     []string `json:"behaviors"`
}

type User struct {
	Username string                 `json:"username"`
	Password string                 `json:"password"`
	Meta     map[string]interface{} `json:"meta"`
}

type Service struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	ExternalIP string `json:"external_ip"`
	InternalIP string `json:"internal_ip`
	Plan       Plan   `json:"plan"`
	State      string `json:"state"`
}

func (p Plan) BehaviorPresent(b string) bool {
	for _, behavior := range p.Behaviors {
		if behavior == b {
			return true
		}
	}
	return false
}
