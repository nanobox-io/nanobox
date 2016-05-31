package portal

import(
	"crypto/tls"
	"fmt"
	"net/http"
	"io"
	"io/ioutil"
	"bytes"
	"encoding/json"
)

type PortalClient struct {
	host string
	token string
}

func init() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
}

func New(host, token string) PortalClient {
	return PortalClient{host, token}
}

func (self PortalClient) GetServices() ([]Service, error) {
	// Get /services
	services := []Service{}
	err := self.do("GET", "/services", nil, &services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
func (self PortalClient) CreateService(service Service) error {
	// Post /services
	return self.do("POST", "/services", service, nil)
}
func (self PortalClient) UpdateServices(services []Service) error {
	// Put /services
	return self.do("PUT", "/services", services, nil)
}
func (self PortalClient) UpdateService(id string, service Service) error {
	// Put /services/:service_id
	return self.do("PUT", "/services/"+id, service, nil)
}
func (self PortalClient) GetService(id string) (*Service, error) {
	// Get /services/:service_id
	services := &Service{}
	err := self.do("GET", "/services/"+id, nil, services)
	if err != nil {
		return nil, err
	}
	return services, nil
}
func (self PortalClient) DeleteService(id string) error {
	// Delete /services/:service_id
	return self.do("DELETE", "/services/"+id, nil, nil)
}
func (self PortalClient) GetServer(id string) ([]Server, error) {
	// Get /services/:service_id/servers
	server := []Server{}
	err := self.do("GET", "/services/"+id+"/servers", nil, &server)
	if err != nil {
		return nil, err
	}
	return server, nil
}
func (self PortalClient) CreateServer(id string, server Server) error {
	// Post /services/:service_id/servers
	return self.do("POST", "/services/"+id+"/servers", server, nil)
}
func (self PortalClient) UpdateServers(id string, servers []Server) error {
	// Put /services/:service_id/servers
	return self.do("PUT", "/services/"+id+"/servers", servers, nil)
}
func (self PortalClient) UpdateServer(serviceId, id string, server Server) error {
	// Get /services/:service_id/servers/:server_id
	return self.do("PUT", "/services/"+serviceId+"/servers/"+id, server, nil)
}
func (self PortalClient) DeleteServer(serviceId, id string) error {
	// Delete /services/:service_id/servers/:server_id
	return self.do("DELETE", "/services/"+serviceId+"/servers/"+id, nil, nil)
}
func (self PortalClient) GetRoutes() ([]Route, error) {
	// Get /routes
	routes := []Route{}
	err := self.do("GET", "/routes", nil, &routes)
	if err != nil {
		return nil, err
	}
	return routes, nil
}
func (self PortalClient) CreateRoute(route Route) error {
	// Post /routes
	return self.do("POST", "/routes", route, nil)
}
func (self PortalClient) UpdateRoutes(routes []Route) error {
	// Put /routes
	return self.do("PUT", "/routes", routes, nil)
}
func (self PortalClient) DeleteRoute(route Route) error {
	// Delete /routes
	return self.do("DELETE", "/routes", route, nil)
}
func (self PortalClient) GetCert() ([]CertBundle, error) {
	// Delete /certs
	certs := []CertBundle{}
	err := self.do("GET", "/certs", nil, &certs)
	if err != nil {
		return nil, err
	}
	return certs, nil
}
func (self PortalClient) CreateCert(cert CertBundle) error {
	// Get /certs
	return self.do("POST", "/certs", cert, nil)
}
func (self PortalClient) UpdateCert(certs []CertBundle) error {
	// Post /certs
	return self.do("PUT", "/certs", certs, nil)
}
func (self PortalClient) DeleteCert(cert CertBundle) error {
	// Put /certs
	return self.do("DELETE", "/certs", cert, nil)
}


func (self PortalClient) do(method, path string, requestBody, responseBody interface{}) error {
	var rbodyReader io.Reader
	if requestBody != nil {
		jsonBytes, err := json.Marshal(requestBody)
		if err != nil {
			return err
		}
		rbodyReader = bytes.NewBuffer(jsonBytes)
	}

	req, err := http.NewRequest(method, fmt.Sprintf("https://%s/%s", self.host, path), rbodyReader)
	if err != nil {
		return err
	}
	req.Header.Add("X-NANOBOX-TOKEN", self.token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if responseBody != nil {
		b, err := ioutil.ReadAll(res.Body)
		err = json.Unmarshal(b, responseBody)
		if err != nil {
			return err
		}
	}
	return nil
}

