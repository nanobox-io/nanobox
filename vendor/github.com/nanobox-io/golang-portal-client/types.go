package portal

type (
	Server struct {
		// todo: change "Id" to "name" (for clarity)
		Id             string `json:"id,omitempty"`
		Host           string `json:"host"`
		Port           int    `json:"port"`
		Forwarder      string `json:"forwarder"`
		Weight         int    `json:"weight"`
		UpperThreshold int    `json:"upper_threshold"`
		LowerThreshold int    `json:"lower_threshold"`
	}
	Service struct {
		Id          string   `json:"id,omitempty"`
		Host        string   `json:"host"`
		Interface   string   `json:"interface,omitempty"`
		Port        int      `json:"port"`
		Type        string   `json:"type"`
		Scheduler   string   `json:"scheduler"`
		Persistence int      `json:"persistence"`
		Netmask     string   `json:"netmask"`
		Servers     []Server `json:"servers,omitempty"`
	}

	Route struct {
		// defines match characteristics
		SubDomain string `json:"subdomain"` // subdomain to match on - "admin"
		Domain    string `json:"domain"`    // domain to match on - "myapp.com"
		Path      string `json:"path"`      // route to match on - "/admin"
		// defines actions
		Targets []string `json:"targets"` // ips of servers - ["http://127.0.0.1:8080/app1","http://127.0.0.2"] (optional)
		FwdPath string   `json:"fwdpath"` // path to forward to targets - "/goadmin" incoming req: test.com/admin -> 127.0.0.1/goadmin (optional)
		Page    string   `json:"page"`    // page to serve instead of routing to targets - "<HTML>We are fixing it</HTML>" (optional)
	}

	CertBundle struct {
		Cert string `json:"cert"`
		Key  string `json:"key"`
	}
)
