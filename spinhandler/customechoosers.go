package spinhandler

import (
	"context"
	"errors"
	"github.com/ArmadaStore/comms/rpc/taskToCargoMgr"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler/filter"
	"github.com/armadanet/spinner/spinhandler/sort"
	"github.com/stretchr/stew/slice"
	"google.golang.org/grpc"
	"log"
)

type CustomChooser struct {
	// LastChoice	string
	filters   map[string]filter.Filter
	sort      map[string]sort.Sort
	filterKey []string
}

func InitCustomChooser() CustomChooser {
	chooser := CustomChooser{
		filters: make(map[string]filter.Filter),
		sort: make(map[string]sort.Sort),
	}

	chooser.filters["FreePorts"] = &filter.PortFilter{}
	chooser.filters["Public"] = &filter.PublicFilter{}
	chooser.filters["Resource"] = &filter.ResourceFilter{}
	chooser.filters["SoftResource"] = &filter.SoftResFilter{}
	chooser.filters["Tag"] = &filter.TagFilter{}
	chooser.filters["FirstDeployment"] = &filter.AffinityFilter{}

	chooser.filterKey = []string{"FreePorts", "Public", "Resource", "FirstDeployment", "Tag"}

	chooser.sort["LeastUsage"] = &sort.LeastRecSort{}
	chooser.sort["Geolocation"] = &sort.GeoSort{}
	return chooser
}

func (r *CustomChooser) Register(id string, kind string, ip string, port string) bool {
	//conn, err := grpc.Dial(ip+":"+port, grpc.WithInsecure())
	//if err != nil {
	//	log.Printf("Cannot access to remote scheduler.")
	//	return false
	//}
	////TODO: modify
	//service := taskToCargoMgr.NewRpcTaskToCargoMgrClient(conn)
	//if kind == "filter" {
	//	c := filter.CustomFilter{service: service}
	//	r.filters[id] = &c
	//}
	return true
}

func (r *CustomChooser) F(c ClientMap, tq *spincomm.TaskRequest) (string, *taskToCargoMgr.Cargos, error) {
	newclients := make(map[string]spinclient.Client)
	for _, k := range c.Keys() {
		newclients[k], _ = c.Get(k)
	}

	var (
		soft bool
		err  error
		ErrNoNode = errors.New("no resource")
	)

	requiredFilters := tq.GetTaskspec().GetFilters()

	for _, f := range r.filterKey {
		ok := slice.Contains(requiredFilters, f)
		if !ok {
			continue
		}
		filter, _ := r.filters[f]
		err = filter.FilterNode(tq, newclients)

		if err != nil {
			if err.Error() == ErrNoNode.Error() {
				soft = true
				log.Println("considering real time resource")
				r.filters["SoftResource"].FilterNode(tq, newclients)
			} else {
				return "", nil, err
			}
		}

		if len(newclients) == 0 {
			err := errors.New("no resource")
			return "", nil, err
		}
	}
	sortPlugin := tq.GetTaskspec().GetSort()
	sortResult := make([]string, len(newclients))
	if sortPolicy, ok := r.sort[sortPlugin]; ok {
		sortResult = sortPolicy.SortNode(tq.GetTaskspec(), newclients, soft)
	} else {
		for key, _ := range newclients {
			sortResult = append(sortResult, key)
		}
	}

	// Contact with cargo manager
	var service taskToCargoMgr.RpcTaskToCargoMgrClient
	var conn *grpc.ClientConn
	cargoFlag := false
	if tq.GetTaskspec().GetCargoSpec() != nil {
		cargoFlag = true
		//TODO: change to a dynamic address "cargoMgr:port"
		conn, err = grpc.Dial("128.101.118.101"+":"+"8080", grpc.WithInsecure())
		if err != nil {
			cargoFlag = false
			log.Printf("Cannot access to Cargo Manager")
		}
		service = taskToCargoMgr.NewRpcTaskToCargoMgrClient(conn)
	}

	//double check
	var cargos *taskToCargoMgr.Cargos
	for _, id := range sortResult {
		if client, ok := c.Get(id); ok {
			if cargoFlag {
				nodeInfo := client.NodeInfo()
				lat, lon := nodeInfo.Lat, nodeInfo.Lon
				if err != nil {
					continue
				}
				req := taskToCargoMgr.RequesterInfo{
					Lat: lat,
					Lon: lon,
					Size: tq.GetTaskspec().GetCargoSpec().GetSize(),
					NReplicas: tq.GetTaskspec().GetCargoSpec().GetNReplica(),
					AppID: tq.GetAppId().GetValue(),
				}
				cargos, err = service.RequestCargo(context.Background(), &req)
				if err != nil {
					log.Println(err)
				}
				conn.Close()
			}
			return id, cargos, nil
		}
	}

	return "", nil, errors.New("no resource")
}
