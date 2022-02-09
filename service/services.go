/* For license and copyright information please see LEGAL file in repository */

package service

import (
	"../protocol"
)

// Services store all application service
type Services struct {
	poolByID        map[uint64]protocol.Service
	poolByMediaType map[string]protocol.Service
}

// Init use to initialize
func (ss *Services) Init() {
	ss.poolByID = make(map[uint64]protocol.Service, 512)
	ss.poolByMediaType = make(map[string]protocol.Service, 512)
}

// RegisterService use to set or change specific service detail!
func (ss *Services) RegisterService(s protocol.Service) {
	ss.registerServiceByMediaType(s)
	ss.registerServiceByURI(s)
}

func (ss *Services) registerServiceByMediaType(s protocol.Service) {
	var serviceID = s.ID()
	if ss.GetServiceByID(serviceID) != nil {
		// This condition will just be true in the dev phase.
		panic("ID associated for '" + s.MediaType().MediaType() + "' Used before for other service and not legal to reuse same ID for other services\n" +
			"Exiting service MediaType is: " + ss.GetServiceByID(serviceID).MediaType().MediaType())
	} else {
		ss.poolByID[serviceID] = s
	}
}

func (ss *Services) registerServiceByURI(s protocol.Service) {
	var serviceURI = s.URI()
	if serviceURI != "" {
		if ss.GetServiceByMediaType(serviceURI) != nil {
			// This condition will just be true in the dev phase.
			panic("URI associated for '" + s.MediaType().MediaType() + " service with `" + serviceURI + "` as URI, Used before for other service and not legal to reuse URI for other services\n" +
				"Exiting service MediaType is: " + ss.GetServiceByMediaType(serviceURI).MediaType().MediaType())
		} else {
			ss.poolByMediaType[serviceURI] = s
		}
	}
}

// GetServiceByID use to get specific service handler by service ID
func (ss *Services) GetServiceByID(serviceID uint64) protocol.Service {
	return ss.poolByID[serviceID]
}

// GetServiceByMediaType use to get specific service handler by service URI
func (ss *Services) GetServiceByMediaType(uri string) protocol.Service {
	return ss.poolByMediaType[uri]
}

// DeleteService use to delete specific service in services list.
func (ss *Services) DeleteService(s protocol.Service) {
	delete(ss.poolByID, s.ID())
	delete(ss.poolByMediaType, s.URI())
}
