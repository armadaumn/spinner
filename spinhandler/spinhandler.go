package spinhandler

import (
	"errors"
	"github.com/ArmadaStore/comms/rpc/taskToCargoMgr"
	// "github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"sync"
	// "context"
	// "errors"
)

type handler struct {
	mutex 		*sync.Mutex
	clientmap	*clientmap
}

type Handler interface{
	AddClient(client spinclient.Client) error
	RemoveClient(id string) error
	ChooseClient(ch Chooser, req *spincomm.TaskRequest) (spinclient.Client, *taskToCargoMgr.Cargos, error)
	ListClientIds() []string
	GetClient(id string) (spinclient.Client, bool)
	// ConnectClient(id string) error
	UpdateClient(status *spincomm.NodeInfo) error
}

func New() Handler {
	return &handler{
		mutex: &sync.Mutex{},
		clientmap: newclientmap(),
	}
}

func (h *handler) AddClient(client spinclient.Client) error {
	err := h.clientmap.add(client)
	return err
}

func (h *handler) RemoveClient(id string) error {
	err := h.clientmap.remove(id)
	return err
}

func (h *handler) ChooseClient(ch Chooser, req *spincomm.TaskRequest) (spinclient.Client, *taskToCargoMgr.Cargos, error) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	cid, cargos, err := ch.F(h.clientmap, req)
	if err != nil {
		return nil, cargos, err
	}
	cl, ok := h.GetClient(cid)
	if !ok {
		return nil, cargos, errors.New("no resource")
	}
	//apps := cl.GetApps()
	//appid := req.GetAppId().GetValue()
	//for _, app := range apps {
	//	if appid == app {
	//		return nil, cargos, errors.New("task is present")
	//	}
	//}
	cl.AppendApps(req.GetAppId().Value)
	cl.UpdateAllocation(req.GetTaskspec().GetResourceMap())
	return cl, cargos, nil
}

func (h *handler) ListClientIds() []string {
	return h.clientmap.Keys()
}

func (h *handler) GetClient(id string) (spinclient.Client, bool) {
	return h.clientmap.Get(id)
}

func (h *handler) UpdateClient(status *spincomm.NodeInfo) error {
	return h.clientmap.update(status)
}

