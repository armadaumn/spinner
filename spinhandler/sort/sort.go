package sort

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
)

type Sort interface {
	SortNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client, soft bool) []string
}
