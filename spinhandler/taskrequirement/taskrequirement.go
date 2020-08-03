package taskrequirement

type TaskRequirement struct {
	Filters []string
	Sort    string

	ResourceMap map[string]ResourceRequirement
	Ports       map[string]struct{}
	IsPublic    bool
	NumReplicas int64
}

type ResourceRequirement struct {
	Requested int64
	Weight    float64
	Required  bool
}
