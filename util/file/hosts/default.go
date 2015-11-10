//
package hosts

type (
	host struct{}
	Host interface {
		HasDomain() bool
		AddDomain()
		RemoveDomain()
	}
)

var (
	Default Host = host{}
)

func (host) HasDomain() bool {
	return HasDomain()
}

func (host) AddDomain() {
	AddDomain()
}

func (host) RemoveDomain() {
	RemoveDomain()
}
