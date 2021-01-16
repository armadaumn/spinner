package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"log"
)

type AffinityFilter struct {

}

func (f *AffinityFilter) FilterNode(tq *spincomm.TaskRequest, clients map[string]spinclient.Client) error {
	appID := tq.GetAppId().GetValue()
	for id, client := range clients {
		isQualified := true
		apps := client.GetApps()
		log.Printf("taskID: %s, client: %s, apps: ", tq.GetTaskId().GetValue(), client.Id())
		log.Println(apps)
		for _, app := range apps {
			if appID == app {
				isQualified = false
				break
			}
		}

		if !isQualified {
			delete(clients, id)
		}
	}
	return nil
}