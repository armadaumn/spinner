package spinclient

import "github.com/armadanet/spinner/spincomm"

type nodeInfo struct {
	OpenPorts       []string
	HostResource    map[string]*spincomm.ResourceStatus
	UsedPorts       map[string]string
	ActiveContainer []string
	Images          []string
}

//type resourceStatus struct {
//	Total      int64
//	Assigned   int64
//	Unassigned int64
//	Available  float64
//}
