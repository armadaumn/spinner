package filter

import (
	"github.com/armadanet/spinner/spinclient"
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
)

type SoftResFilter struct {

}

func (f *SoftResFilter) FilterNode(tq task.TaskRequirement, clients map[string]spinclient.Client) error {
	// Do soft filtering
	for id, client := range clients {
		isSufficient := true
		for res, requirement := range tq.ResourceMap {
			if !requirement.Required {
				continue
			}
			if status, ok := client.Info().HostResource[res]; ok {
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
