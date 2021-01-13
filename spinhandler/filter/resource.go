package filter

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"errors"
	"log"
)

type ResourceFilter struct {
	//Name string
}

//func InitResourceFilter() (*ResourceFilter, error) {
//	return &ResourceFilter{
//		Name: "Resource",
//	}, nil
//}

func (f *ResourceFilter) FilterNode(tq *spincomm.TaskRequest, clients map[string]spinclient.Client) error {
	//hard resource filtering
	newclients := make(map[string]spinclient.Client)
	for k, v := range clients {
		newclients[k] = v
	}
	if len(tq.GetTaskspec().ResourceMap) == 0 {
		//do nothing
		log.Println("passed")
		return nil
	}
	for id, client := range clients {
		isSufficient := true
		for res, requirement := range tq.GetTaskspec().ResourceMap {
			if status, ok := client.Info().HostResource[res]; ok {
				if status.Unassigned < requirement.Requested {
					isSufficient = false
					break
				}
			} else {
				isSufficient = false
				break
			}
		}
		if !isSufficient {
			delete(clients, id)
		}
	}
	if len(clients) == 0 {
		for k, v := range newclients {
			clients[k] = v
		}
		return errors.New("no nodes")
	}
	return nil
}