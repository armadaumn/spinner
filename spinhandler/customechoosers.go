package spinhandler

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spinhandler/filter"
	"github.com/armadanet/spinner/spinhandler/sort"
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
	"errors"
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

func (r *CustomChooser) F(c *clientmap, tq task.TaskRequirement) (string, error) {
	clients := c.clients
	newclients := make(map[string]spinclient.Client)
	for k, v := range clients {
		newclients[k] = v
	}

	var (
		soft bool
		err  error
		ErrNoNode = errors.New("no nodes")
	)

	for _, f := range tq.Filters {
		err = r.filters[f].FilterNode(tq, newclients)

		if err.Error() == ErrNoNode.Error() {
			soft = true
			r.filters["SoftResource"].FilterNode(tq, clients)
		} else if err != nil {
			return "", err
		}

		if len(newclients) == 0 {
			err := errors.New("no clients")
			return "", err
		}
	}
	sortResult := r.sort[tq.Sort].SortNode(tq, newclients, soft)
	return sortResult[0], nil
}