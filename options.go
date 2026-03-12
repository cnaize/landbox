package landbox

import (
	"strconv"
	"strings"
)

type Options struct {
	TCPListen   Ports `json:"tcp_listen" yaml:"tcp_listen"`   // nil: allow all, empty: deny all
	TCPConnect  Ports `json:"tcp_connect" yaml:"tcp_connect"` // nil: allow all, empty: deny all
	DenySockets bool  `json:"deny_sockets" yaml:"deny_sockets"`
	DenySignals bool  `json:"deny_signals" yaml:"deny_signals"`
	EnableDebug bool  `json:"-" yaml:"-"`
}

type Paths []string

func (p Paths) String() string {
	return strings.Join(p, ":")
}

type Ports []uint16

func (p Ports) String() string {
	ports := make([]string, 0, len(p))
	for _, port := range p {
		ports = append(ports, strconv.Itoa(int(port)))
	}

	return strings.Join(ports, ":")
}

func (o Options) Scope() string {
	const (
		socketsKey = "a"
		signalsKey = "s"
	)

	var scope []string
	if o.DenySockets {
		scope = append(scope, socketsKey)
	}
	if o.DenySignals {
		scope = append(scope, signalsKey)
	}

	return strings.Join(scope, ":")
}
