package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
)

type PublicFilter struct {
}

func (f *PublicFilter) FilterNode(tq *spincomm.TaskRequest, clients map[string]spinclient.Client) error {
	for id, client := range clients {
		if len(client.NodeStatus().OpenPorts) == 0 {
			delete(clients, id)
		}
	}
	return nil
}
