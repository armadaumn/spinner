package sort

import (
	task "github.com/armadanet/spinner/spinhandler/taskrequirement"
	"github.com/armadanet/spinner/spinclient"
	"sort"
)

type LeastRecSort struct {

}

func (s *LeastRecSort) SortNode(tq task.TaskRequirement, clients map[string]spinclient.Client, soft bool) []string {
	result := make([]struct {
		id    string
		score float64
	}, len(clients))
	i := 0
	for id, client := range clients {
		var score, weightSum float64
		for res, requirement := range tq.ResourceMap {
			var avail float64
			resStatus := client.Info().Status.HostResource[res]
			if soft {
				avail = resStatus.Available
			} else {
				avail = float64(resStatus.Unassigned) / float64(resStatus.Total) * 100.0
			}
			resScore := avail - float64(requirement.Requested)/ float64(resStatus.Total) * 100.0
			weightSum += requirement.Weight
			score = score + resScore* requirement.Weight
		}
		score = score / weightSum

		result[i].id = id
		result[i].score = score
		i++
	}
	sort.Slice(result, func(i, j int) bool { return result[i].score > result[j].score })
	var ids []string
	for _, r := range result {
		ids = append(ids, r.id)
	}
	return ids
}