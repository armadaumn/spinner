package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"log"
)

type PortFilter struct {
}

func (f *PortFilter) FilterNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client) error {
	//Ports filtering
	if len(tq.Ports) == 0 {
		//do nothing
		log.Println("passed")
		return nil
	}
	for id, client := range clients {
		isUsed := false
		for _, port := range client.Info().UsedPorts {
			if _, ok := tq.Ports[port]; ok {
				isUsed = true
				break
			}
		}
		if isUsed {
			delete(clients, id)
		}
	}
	return nil
}
