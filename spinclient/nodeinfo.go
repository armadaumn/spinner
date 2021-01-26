package spinclient

import "github.com/armadanet/spinner/spincomm"

type nodeStatus struct {
	OpenPorts       []string
	HostResource    map[string]*spincomm.ResourceStatus
	UsedPorts       map[string]string
	ActiveContainer []string
	Images          []string
}

type nodeInfo struct {
	Ip         string
	Port       string
	Lat        float64 //latitude
	Lon        float64 //longitude
	Geoid      string
	ServerType spincomm.Type
	Tags       []string
}

//type resourceStatus struct {
//	Total      int64
//	Assigned   int64
//	Unassigned int64
//	Available  float64
//}
