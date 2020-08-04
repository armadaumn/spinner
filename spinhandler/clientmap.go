package spinhandler

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"sync"
	"errors"
)

type clientmap struct {
	mutex		*sync.Mutex
	clients		map[string]spinclient.Client
}

type ClientMap interface {
	Get(string) (spinclient.Client, bool)
	Keys()		[]string
	Len()		int
}

func newclientmap() *clientmap {
	return &clientmap{
		mutex: &sync.Mutex{},
		clients: make(map[string]spinclient.Client),
	}
}

func (cm *clientmap) Get(key string) (spinclient.Client, bool) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	val, ok := cm.clients[key]
	return val, ok
}

func (cm *clientmap) Keys() []string {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	keys := []string{}
	for key, _ := range cm.clients {
		keys = append(keys, key)
	}
	return keys
}

func (cm *clientmap) Len() int {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	return len(cm.clients)
}

func (cm *clientmap) add(client spinclient.Client) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if _, ok := cm.clients[client.Id()]; ok {
		return errors.New("Requested client already in the system")
	}
	cm.clients[client.Id()] = client
	return nil
}

func (cm *clientmap) remove(id string) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	if _, ok := cm.clients[id]; !ok {
		return errors.New("No such client")
	}
	delete(cm.clients, id)
	return nil
}

func (cm *clientmap) update(status *spincomm.NodeInfo) error {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()
	id := status.GetCaptainId().Value
	client, ok := cm.clients[id]
	if !ok {
		return errors.New("No such client")
	}
	return client.UpdateStatus(status)
}
