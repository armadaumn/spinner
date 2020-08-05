package spinhandler

import (
	"github.com/armadanet/spinner/spincomm"
	// "github.com/armadanet/spinner/spinclient"
	"sort"
	"errors"
)

type RoundRobinChooser struct {
	LastChoice	string
}

func (r *RoundRobinChooser) F (c ClientMap, tq *spincomm.TaskRequest) (string, error) {
	if c.Len() == 0 {
		r.LastChoice = ""
		return "", errors.New("No clients available")
	}
	clients := c.Keys()
	sort.Strings(clients)
	for _, v := range clients {
		if v > r.LastChoice {
			r.LastChoice = v
			return v, nil
		}
	}
	r.LastChoice = clients[0]
	return clients[0], nil
}