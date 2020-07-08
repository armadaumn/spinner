package spinhandler

import (
	"github.com/armadanet/client"
	"sort"
	"errors"
)

type RoundRobinChooser struct {
	LastChoice	string
}

func (r *RoundRobinChooser) F (c ClientMap) (string, error) {
	if c.Len() == 0 {
		r.LastChoice = ""
		return "", errors.New("No clients available")
	}
	clients := sort.Strings(c.Keys())
	for _, v := range clients {
		if v > r.LastChoice {
			r.LastChoice = v
			return v, nil
		}
	}
	r.LastChoice = clients[0]
	return clients[0], nil
}