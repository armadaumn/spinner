package spinclient

type nodeInfo struct {
	OpenPorts       []string
	HostResource    map[string]*resourceStatus
	UsedPorts       []string
	ContainerStatus *containerStatus
}

type resourceStatus struct {
	Total      int64
	Assigned   int64
	Unassigned int64
	Available  float64
}

type containerStatus struct {
	ActiveContainer []string
	Images          []string
}
