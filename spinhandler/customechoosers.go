package spinhandler

import (
	"context"
	"errors"
	"github.com/ArmadaStore/comms/rpc/taskToCargoMgr"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler/filter"
	"github.com/armadanet/spinner/spinhandler/sort"
	"google.golang.org/grpc"
	"log"
)

type CustomChooser struct {
	// LastChoice	string
	filters map[string]filter.Filter
	sort    map[string]sort.Sort
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

	chooser.sort["LeastUsage"] = &sort.LeastRecSort{}
	return chooser
}

func (r *CustomChooser) F(c ClientMap, tq *spincomm.TaskRequest) (string, string, error) {
	newclients := make(map[string]spinclient.Client)
	for _, k := range c.Keys() {
		newclients[k], _ = c.Get(k)
	}

	var (
		soft bool
		err  error
		ErrNoNode = errors.New("no nodes")
	)

	filterPlugins := tq.GetTaskspec().GetFilters()
	for _, f := range filterPlugins {
		err = r.filters[f].FilterNode(tq.GetTaskspec(), newclients)

		if err != nil {
			if err.Error() == ErrNoNode.Error() {
				soft = true
				r.filters["SoftResource"].FilterNode(tq.GetTaskspec(), newclients)
			} else {
				return "", "", err
			}
		}

		if len(newclients) == 0 {
			err := errors.New("no clients")
			return "", "", err
		}
	}
	sortPlugin := tq.GetTaskspec().GetSort()
	sortResult := r.sort[sortPlugin].SortNode(tq.GetTaskspec(), newclients, soft)

	// Contact with cargo manager
	var service taskToCargoMgr.RpcTaskToCargoMgrClient
	var conn *grpc.ClientConn
	cargoFlag := false
	if tq.GetTaskspec().GetCargoSpec() != nil {
		cargoFlag = true
		//TODO: change to a dynamic address "cargoMgr:port"
		conn, err = grpc.Dial("127.0.0.1"+":"+"5913", grpc.WithInsecure())
		if err != nil {
			cargoFlag = false
			log.Printf("Cannot access to Cargo Manager")
		}
		service = taskToCargoMgr.NewRpcTaskToCargoMgrClient(conn)
	}

	//TODO: double check
	cargoAddr := ""
	for _, id := range sortResult {
		if client, ok := c.Get(id); ok {
			if cargoFlag {
				lat, lon, err := client.Location()
				if err != nil {
					continue
				}
				req := taskToCargoMgr.RequesterInfo{
					Lat: lat,
					Lon: lon,
					Size: tq.GetTaskspec().GetCargoSpec().GetSize(),
					NReplicas: tq.GetTaskspec().GetCargoSpec().GetNReplica(),
				}
				log.Println(req)
				res, err := service.RequestCargo(context.Background(), &req)
				if err != nil {
					log.Println(err)
				}
				cargoAddr = res.GetIPPort()
				log.Println("cargo Address: " + cargoAddr)
				conn.Close()
			}
			return id, cargoAddr, nil
		}
	}

	return "", cargoAddr, errors.New("no node")
}