package structs

import "strings"

// NomadInfo holds top level meta data about the cluster
type NomadInfo struct {
	Servers Servers
	Clients Clients
}

type Servers []string
type Clients []string

func (s Servers) GoString() string {
	return strings.Join(s, ",")
}

func (c Clients) GoString() string {
	return strings.Join(c, ",")
}
