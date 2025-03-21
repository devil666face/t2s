// automatically generated by stateify.

//go:build linux
// +build linux

package stopfd

import (
	"context"

	"gvisor.dev/gvisor/pkg/state"
)

func (sf *StopFD) StateTypeName() string {
	return "pkg/tcpip/link/stopfd.StopFD"
}

func (sf *StopFD) StateFields() []string {
	return []string{
		"EFD",
	}
}

func (sf *StopFD) beforeSave() {}

// +checklocksignore
func (sf *StopFD) StateSave(stateSinkObject state.Sink) {
	sf.beforeSave()
	stateSinkObject.Save(0, &sf.EFD)
}

func (sf *StopFD) afterLoad(context.Context) {}

// +checklocksignore
func (sf *StopFD) StateLoad(ctx context.Context, stateSourceObject state.Source) {
	stateSourceObject.Load(0, &sf.EFD)
}

func init() {
	state.Register((*StopFD)(nil))
}
