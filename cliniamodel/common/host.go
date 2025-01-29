package common

import (
	"fmt"
)

type HostScheme string

const (
	HTTP HostScheme = "http"
	HTTPS HostScheme = "https"
)

type Host struct {
	Url string
	Port int
	Scheme HostScheme
}

// String returns the host in the format of "scheme://url:port"
func (h Host) String() string {
	return fmt.Sprintf("%s://%s:%d", h.Scheme, h.Url, h.Port)
}

// Host returns the host in the format of "url:port"
func (h Host) Host() string {
	return fmt.Sprintf("%s:%d", h.Url, h.Port)
}



