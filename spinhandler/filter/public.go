package filter

import (
	"github.com/armadanet/spinner/spinclient"
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
)

type PublicFilter struct {
}

func (f *PublicFilter) FilterNode(tq task.TaskRequirement, clients map[string]spinclient.Client) error {
	for id, client := range clients {
		if len(client.Info().OpenPorts) == 0 {
			delete(clients, id)
		}
	}
	return nil
}
