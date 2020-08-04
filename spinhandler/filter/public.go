package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
)

type PublicFilter struct {
}

func (f *PublicFilter) FilterNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client) error {
	for id, client := range clients {
		if len(client.Info().OpenPorts) == 0 {
			delete(clients, id)
		}
	}
	return nil
}
