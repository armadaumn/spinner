package sort

import (
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/mmcloughlin/geohash"
	"log"
	"math"
	"sort"
	"strings"
)

type GeoSort struct {

}

func (s *GeoSort) SortNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client, soft bool) []string {
	result := make([]struct {
		id         string
		serverType spincomm.Type
		score      int
		dist       float64
		availCpu   int64
	}, len(clients))

	//dataSources := tq.GetDataSources()
	index := 0
	ds := tq.GetDataSources()
	requiredCpu := tq.GetResourceMap()["CPU"].Requested
	sourceGeoID := geohash.EncodeWithPrecision(ds.GetLat(), ds.GetLon(), 4)
	log.Println(sourceGeoID)
	for id, captain := range clients {
		result[index].id = id
		captainInfo := captain.NodeInfo()
		result[index].serverType = captainInfo.ServerType
		if captain.NodeStatus().HostResource["CPU"].Unassigned <= requiredCpu {
			result[index].availCpu = captain.NodeStatus().HostResource["CPU"].Unassigned
		} else {
			result[index].availCpu = requiredCpu
		}

		// Get neighbors of the data source
		neighbor := geohash.Neighbors(sourceGeoID[:4])
		neighbor = append(neighbor, sourceGeoID)

		captainGeoID := strings.SplitN(captainInfo.Geoid, "-", 2)[0]
		//log.Printf("captain: %s, geoID: %s", captain.Id(), captainGeoID)
		result[index].score = proximityComparison(neighbor, []rune(captainGeoID))
		nodeLat, nodeLon := captainInfo.Lat, captainInfo.Lon
		dist := getDistance(ds.GetLat(), ds.GetLon(), nodeLat, nodeLon)
		result[index].dist = dist
		index++
	}
	// Sorting the result
	sort.Slice(result, func(i, j int) bool {
		if result[i].score != result[j].score {
			return result[i].score > result[j].score
		}
		i_weightedScore := 0.5 * (float64(result[i].serverType) * 2 + 2) + 0.5 * float64(result[i].availCpu)
		j_weightedScore := 0.5 * (float64(result[j].serverType) * 2 + 2) + 0.5 * float64(result[j].availCpu)
		return i_weightedScore > j_weightedScore
	})

	log.Println(result)
	var ids []string
	for _, r := range result {
		ids = append(ids, r.id)
	}
	return ids
}

func proximityComparison(neighbor []string, ghDst []rune) int {
	maxCount := 0

	for _, src := range neighbor {
		ghSrc := []rune(src)
		ghSrcLen := len(ghSrc)
		prefixMatchCount := 0

		for i := 0; i < ghSrcLen && i < len(ghDst); i++ {
			if ghSrc[i] == ghDst[i] {
				prefixMatchCount++
			} else {
				break
			}
		}

		if prefixMatchCount > maxCount {
			maxCount = prefixMatchCount
		}
	}
	return maxCount
}

func getDistance(srcLat, srcLon, nodeLat, nodeLon float64) float64 {
	radiusSrcLat := math.Pi * srcLat / 180
	radiusNodeLat := math.Pi * nodeLat / 180

	radiusTheta := math.Pi * (srcLon - nodeLon) / 180

	dist := math.Sin(radiusSrcLat)*math.Sin(radiusNodeLat) +
		math.Cos(radiusSrcLat)*math.Cos(radiusNodeLat)*math.Cos(radiusTheta)
	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	return dist
}