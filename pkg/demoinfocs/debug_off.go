//go:build !debugdemoinfocs
// +build !debugdemoinfocs

// This file is just a bunch of NOPs for the release build, see debug_on.go for debugging stuff

package demoinfocs

import (
	msg "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/msg"
	st "github.com/rechedev9/demo-parser-rsn/pkg/demoinfocs/sendtables"
)

func debugGameEvent(descriptor *msg.CMsgSource1LegacyGameEventListDescriptorT, ge *msg.CMsgSource1LegacyGameEvent) {
	// NOP
}

func debugIngameTick(tickNr int) {
	// NOP
}

func debugAllServerClasses(classes st.ServerClasses) {
	// NOP
}
