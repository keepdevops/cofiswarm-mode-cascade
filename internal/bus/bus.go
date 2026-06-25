// Package bus wires the mode relay onto the NATS observer bus via the shared
// cofiswarm-observer-sdk service component: it announces presence and serves
// .mode.<name>.{info,health} alongside the relay's HTTP API.
package bus

import (
	"github.com/keepdevops/cofiswarm-observer-sdk/pkg/servicecomponent"
)

// Routes wires the relay's capability subjects. Reply field names carry schema_version for the
// major-version gate, mirroring the other service components.
func Routes(name string) map[string]servicecomponent.Handler {
	return map[string]servicecomponent.Handler{
		servicecomponent.Prefix + ".mode." + name + ".info":   infoHandler(name),
		servicecomponent.Prefix + ".mode." + name + ".health": healthHandler(),
	}
}

func infoHandler(name string) servicecomponent.Handler {
	return func([]byte) (any, error) {
		return infoReply{SchemaVersion: servicecomponent.SchemaVersion, OK: true, Mode: name}, nil
	}
}

func healthHandler() servicecomponent.Handler {
	return func([]byte) (any, error) {
		return healthReply{SchemaVersion: servicecomponent.SchemaVersion, OK: true, Status: "ok"}, nil
	}
}

type infoReply struct {
	SchemaVersion string `json:"schema_version"`
	OK            bool   `json:"ok"`
	Error         string `json:"error,omitempty"`
	Mode          string `json:"mode"`
}

type healthReply struct {
	SchemaVersion string `json:"schema_version"`
	OK            bool   `json:"ok"`
	Error         string `json:"error,omitempty"`
	Status        string `json:"status"`
}
