package filter

import (
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinclient"
)

type Filter interface {
	FilterNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client) error
}
