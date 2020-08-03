package sort

import (
	"github.com/armadanet/spinner/spinclient"
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
)

type Sort interface {
	SortNode(tq task.TaskRequirement, clients map[string]spinclient.Client, soft bool) []string
}
