package models

type (
	// TunnelConfig contains the endpoint information.
	TunnelConfig struct {
		AppName    string // name of app to tunnel to
		ListenPort int    // local port to listen on
		DestPort   int    // port to tunnel to
		Component  string // component to tunnel to
	}

	TunnelInfo struct {
		Name  string `json:"name,omitempty"`  // component name being tunneled to
		Token string `json:"token,omitempty"` // token to complete the tunnel
		URL   string `json:"url,omitempty"`   // url/ip of nanoagent
		Port  int    `json:"port,omitempty"`  // port to tunnel to. *must* omitempty. odin will use whatever port is passed (if 0, don't send)
	}
)
