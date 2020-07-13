package spinhandler

import (
	"github.com/armadanet/spinner/spinresp"
	"github.com/armadanet/spinner/spinclient"
	"sync"
	// "errors"
)

type handler struct {
	mutex 		*sync.Mutex
	clientmap	*clientmap
}

type Handler interface{
	AddClient(request *spinresp.JoinRequest, stream spinresp.Spinner_AttachServer) error
	RemoveClient(id string) error
	ChooseClient(ch Chooser) (string, error)
	// ConnectClient(id string) error
}

func New() Handler {
	return &handler{
		mutex: &sync.Mutex{},
		clientmap: newclientmap(),
	}
}

func (h *handler) AddClient(request *spinresp.JoinRequest, stream spinresp.Spinner_AttachServer) error {
	client, err := spinclient.RequestClient(request, stream)
	if err != nil {return err}
	err = h.clientmap.add(client)
	return err
}

func (h *handler) RemoveClient(id string) error {
	err := h.clientmap.remove(id)
	return err
}

func (h *handler) ChooseClient(ch Chooser) (string, error) {
	return ch.F(h.clientmap)
}