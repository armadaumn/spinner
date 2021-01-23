package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
)

type SoftResFilter struct {

}

func (f *SoftResFilter) FilterNode(tq *spincomm.TaskRequest, clients map[string]spinclient.Client) error {
	// Do soft filtering
	for id, client := range clients {
		isSufficient := true
		for res, requirement := range tq.GetTaskspec().ResourceMap {
			if !requirement.Required {
				continue
			}
			if status, ok := client.NodeStatus().HostResource[res]; ok {
				if res == "CPU" || res == "Memory" {
					percent := float64(requirement.Requested) / float64(status.Total) * 100.0
					if status.Available < percent {
						isSufficient = false
					}
				} else if status.Unassigned < requirement.Requested {
					isSufficient = false
				}
			} else {
				isSufficient = false
			}
			if !isSufficient {
				break
			}
		}
		if !isSufficient {
			delete(clients, id)
		}
	}
	return nil
}
