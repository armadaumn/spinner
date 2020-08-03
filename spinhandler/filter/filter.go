package filter

import (
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
	"github.com/armadanet/spinner/spinclient"
)

type Filter interface {
	FilterNode(tq task.TaskRequirement, clients map[string]spinclient.Client) error
}
