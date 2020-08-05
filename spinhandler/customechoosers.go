package spinhandler

import (
	"errors"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler/filter"
	"github.com/armadanet/spinner/spinhandler/sort"
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

func (r *CustomChooser) F(c ClientMap, tq *spincomm.TaskRequest) (string, error) {
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
				return "", err
			}
		}

		if len(newclients) == 0 {
			err := errors.New("no clients")
			return "", err
		}
	}
	sortPlugin := tq.GetTaskspec().GetSort()
	sortResult := r.sort[sortPlugin].SortNode(tq.GetTaskspec(), newclients, soft)
	//TODO: double check
	for _, id := range sortResult {
		if _, ok := c.Get(id); ok {
			return id, nil
		}
	}
	return "", errors.New("no node")
}