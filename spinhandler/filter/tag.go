package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"log"
)

type TagFilter struct {

}

func (f *TagFilter) FilterNode(tq *spincomm.TaskRequest, clients map[string]spinclient.Client) error {
	tags := tq.GetTaskspec().GetTags()
	storeClients := make(map[string]spinclient.Client)
	for k, v := range clients {
		storeClients[k] = v
	}

	for id, captain := range clients {
		isOverlapped := false
		nodeTags := captain.NodeInfo().Tags
		for _, tag := range tags {
			log.Printf(tag)
			for _, nodeTag := range nodeTags {
				log.Println(nodeTag)
				if tag == nodeTag {
					isOverlapped = true
					break
				}
			}
			if isOverlapped {
				break
			}
		}
		if !isOverlapped {
			delete(clients, id)
		}
	}

	if len(clients) == 0 {
		for k, v := range storeClients {
			clients[k] = v
		}
	}
	return nil
}