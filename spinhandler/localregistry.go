package spinhandler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

type registry struct {
	mutex *sync.Mutex
	url   string
	repos map[string]interface{}
}

type Registry interface {
	GetUrl() string
	SetUrl(url string)
	GetRepos() map[string]interface{}
	SetRepos(repos []string)
	UpdateImageList()
}

func NewRegistry(url string) Registry {
	return &registry{
		mutex: &sync.Mutex{},
		url: url,
		repos: make(map[string]interface{}),
	}
}

func (r *registry) GetUrl() string {
	return r.url
}

func (r *registry) SetUrl(url string) {
	r.url = url
}

func (r *registry) GetRepos() map[string]interface{} {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	return r.repos
}

func (r *registry) SetRepos(repoList []string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	for _, image := range repoList {
		if _, ok := r.repos[image]; !ok {
			r.repos[image] = nil
		}
	}
}

func (r *registry) UpdateImageList() {
	for {
		url := "http://" + r.url + "/v2/_catalog"
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		resp.Body.Close()

		var f interface{}
		err = json.Unmarshal(body, &f)
		if err != nil {
			log.Println(err)
		}
		m := f.(map[string]interface{})
		repos := m["repositories"].([]interface{})

		repoList := make([]string, len(repos))
		for _, r := range repos {
			repoList = append(repoList, r.(string))
		}
		r.SetRepos(repoList)
		time.Sleep(5 * time.Second)
	}
}
