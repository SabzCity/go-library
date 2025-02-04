/* For license and copyright information please see the LEGAL file in the code repository */

package net_p

import (
	capsule_p "memar/computer/capsule/protocol"
	error_p "memar/error/protocol"
	"memar/time/duration"
)

/*
**********************************************************************************
Transport (OSI Layer 4: Transport)

https://en.wikipedia.org/wiki/Transport_layer
**********************************************************************************
*/

// OSI_Transport provides services such as connection-oriented communication, reliability, flow control.
// It must also implement chunks managing like sRPC, QUIC, TCP, UDP, ...
type OSI_Transport interface {
	capsule_p.LifeCycle
	Framer
	NetworkAddress // string form of address (for example, "tcp://443", "srpc://1254872653")

	OSI_Transport_LowLevelAPIs
}

// OSI_Transport_LowLevelAPIs is low level APIs, don't use them in the services layer, if you don't know how it can be effect the application.
// It will use in chunks managing packages e.g. sRPC, TCP, UDP, ... or protocols wrappers e.g. QUIC, HTTPv2, ...
type OSI_Transport_LowLevelAPIs interface {
	FrameWriter
	OSI_Transport_Options
}

type OSI_Transport_Options interface {
	// release any underling data reference until call time without need to release socket itself
	Discard(n int) (discarded int, err error_p.Error)
	SetLinger(d duration.Second) error
	SetKeepAlivePeriod(d duration.Second) error
	SetNoDelay(noDelay bool) error
}
