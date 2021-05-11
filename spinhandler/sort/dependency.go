package sort

import (
	"encoding/json"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type LayerSort struct {
	ImageLayer map[string][]string
}

func (s *LayerSort) SortNode(tq *spincomm.TaskSpec, clients map[string]spinclient.Client, soft bool) []string {
	image := tq.GetImage()
	tag := tq.GetImageVersion()
	dependency, ok := s.ImageLayer[image+":"+tag]
	if !ok {
		token, err := getUserToken(image)
		layers, err := getLayerDigest(image, tag, token)
		if err != nil {
			log.Println(err)
		}
		if layers != nil {
			log.Println(layers)
			for _, layer := range layers {
				l := layer.(map[string]interface{})
				// Remove sha256 prefix
				blobsum := l["blobSum"].(string)[7:19]
				log.Println(blobsum)
				dependency = append(dependency, blobsum)
			}
		}
		s.ImageLayer[image+":"+tag] = dependency
	}

	// sort nodes
	result := make([]struct {
		id         string
		score      int
	}, len(clients))

	index := 0
	for id, captain := range clients {
		result[index].id = id
		layers := captain.NodeStatus().Layers
		overlap := 0

		for _, d := range dependency {
			if _, ok = layers[d]; ok {
				overlap = overlap + 1
			}
		}
		result[index].score = overlap
		index++
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].score > result[j].score
	})
	log.Println(result)
	var ids []string
	for _, r := range result {
		ids = append(ids, r.id)
	}
	return ids
}

func getUserToken(image string) (string, error) {
	params := "service=registry.docker.io&scope=repository:" + image + ":pull"
	url := "https://auth.docker.io/token?"+params
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	//req.SetBasicAuth(dummyUser, dummyPassword)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	//buf := new(strings.Builder)
	//_, err = io.Copy(buf, resp.Body)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	//log.Println(m)
	token := m["token"].(string)
	//log.Printf(token)
	return token, nil
}

func getLayerDigest(image string, tag string, token string) ([]interface{}, error) {
	url := "https://registry-1.docker.io/v2/" + image
	url = url + "/manifests/" + tag
	log.Println(url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	if m["fsLayers"] == nil {
		return nil, nil
	}
	layers := m["fsLayers"].([]interface{})
	return layers, nil
}


