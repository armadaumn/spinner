package sort

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/mmcloughlin/geohash"
	"log"
	"sort"
	"strings"
)

type GeoSort struct {

}

func (s *GeoSort) SortNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client, soft bool) []string {
	result := make([]struct {
		id    string
		score int
	}, len(clients))

	//dataSources := tq.GetDataSources()
	index := 0
	ds := tq.GetDataSources()
	sourceGeoID := geohash.Encode(ds.GetLat(), ds.GetLon())
	log.Println(sourceGeoID)
	for id, captain := range clients {
		result[index].id = id
		captainGeoID := strings.SplitN(captain.Geoid(), "-", 2)[0]
		log.Printf("captain: %s, geoID: %s", captain.Id(), captainGeoID)
		totalScore := 0
		dist := proximityComparison([]rune(sourceGeoID), []rune(captainGeoID))
		totalScore += dist
		result[index].score = totalScore
		index++
	}
	sort.Slice(result, func(i, j int) bool { return result[i].score > result[j].score })
	log.Println(result)
	var ids []string
	for _, r := range result {
		ids = append(ids, r.id)
	}
	return ids
}

func proximityComparison(ghSrc, ghDst []rune) int {
	ghSrcLen := len(ghSrc)

	prefixMatchCount := 0

	for i := 0; i < ghSrcLen; i++ {
		if ghSrc[i] == ghDst[i] {
			prefixMatchCount++
		} else {
			break
		}
	}
	return prefixMatchCount
}