package dns

import "github.com/nanobox-io/nanobox/processor"

// processDevDNSList
type processDevDNSList struct {
	control processor.ProcessControl
}

//
func init() {
	processor.Register("dev_dns_list", devDNSListFunc)
}

//
func devDNSListFunc(conf processor.ProcessControl) (processor.Processor, error) {
	return processDevDNSList{conf}, nil
}

//
func (devDNSList processDevDNSList) Results() processor.ProcessControl {
	return devDNSList.control
}

//
func (devDNSList processDevDNSList) Process() error {
	return nil
}
