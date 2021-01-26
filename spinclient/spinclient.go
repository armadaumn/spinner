package spinclient

import (
	"context"
	"errors"
	"github.com/armadanet/spinner/spincomm"
	"github.com/google/uuid"
	"github.com/mmcloughlin/geohash"
	"log"
)

type client struct {
	id       string
	stream   spincomm.Spinner_AttachServer
	cancel   func()
	taskchan chan *spincomm.TaskRequest
	err      error
	ctx      context.Context
	status   nodeStatus
	info     nodeInfo
	tasks    []string // existing tasks
	apps     []string // existring apps
}

type Client interface {
	Id() string
	SendTask(task *spincomm.TaskRequest) error
	Run() error
	NodeInfo() nodeInfo
	NodeStatus() nodeStatus
	UpdateStatus(status *spincomm.NodeInfo) error
	GetTasks() []string
	GetApps() []string
	AppendApps(appid string)
	UpdateAllocation(req map[string]*spincomm.ResourceRequirement)
}

func RequestClient(ctx context.Context, request *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) (Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
		taskchan: make(chan *spincomm.TaskRequest),
		cancel: cancel,
		ctx: ctx,
		status: nodeStatus{
			HostResource: make(map[string]*spincomm.ResourceStatus),
			UsedPorts:    make(map[string]string),
		},
		info: nodeInfo{
			Ip:         request.GetIP(),
			Port:       request.GetPort(),
			Lat:        request.GetLat(),
			Lon:        request.GetLon(),
			ServerType: request.GetType(),
			Tags:       request.GetTags(),
		},
		tasks: make([]string, 0),
		apps: make([]string, 0),
	}
	//c.ip = "0.0.0.0"
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}

	if c.info.Lat == 0 && c.info.Lon == 0 {
		// TODO: fetch lat and lon
		c.info.Lat = 45.0196
		c.info.Lon = -93.2402
	}
	c.info.Geoid = c.genGeoHashID(c.info.Lat, c.info.Lon)

	return c, nil
}

func (c *client) Id() string {
	return c.id
}

func (c *client) SendTask(task *spincomm.TaskRequest) error {
	if c.err != nil {
		return c.err
	}
	if task == nil {
		return errors.New("Nil task sent")
	}
	c.taskchan <- task
	return nil
}

func (c *client) Run() error {
	cancel := c.cancel
	defer cancel()
	ctx := c.ctx
	for {
		select {
		case task, ok := <- c.taskchan:
			log.Println("Task sent")
			if !ok {
				log.Printf("Task Channel closed for %s\n.", c.Id())
				c.err = errors.New("Task Channel closed")
				return c.err
			}
			if err := c.stream.Send(task); err != nil {
				log.Printf("Error for %s: %v\n", c.Id(), err)
				c.err = err
				return c.err
			}
		case <-ctx.Done():
			log.Printf("Context ended for %s: %v\n", c.Id(), ctx.Err())
			c.err = errors.New("context ended")
			return c.err

		case <-c.stream.Context().Done():
			log.Printf("Stream ended by %s: %v\n", c.Id(), ctx.Err())
			c.err = errors.New("captain left")
			return c.err
		}
	}
}

func (c *client) NodeStatus() nodeStatus {
	return c.status
}

func (c *client) UpdateStatus(status *spincomm.NodeInfo) error {
	//TODO: Remove testing output information
	c.status.UsedPorts = status.GetUsedPorts()
	c.status.ActiveContainer = status.GetContainerStatus().GetActiveContainer()
	c.status.Images = status.GetContainerStatus().GetImages()
	c.status.HostResource = status.GetHostResource()
	//c.apps = status.GetAppIDs()
	c.tasks = status.GetTaskIDs()
	log.Println("after:", c.status)
	//log.Println("app: ", c.apps)
	return nil
}

func (c *client) NodeInfo() nodeInfo {
	return c.info
}

func (c *client) genGeoHashID(lat float64, lon float64) string {
	geohashIDstr := geohash.EncodeWithPrecision(lat, lon, 4)
	uuID, err := uuid.NewUUID()
	if err != nil {
		log.Println(err)
	}
	geoID := geohashIDstr + "-" + uuID.String()
	return geoID
}

func (c *client) GetTasks() []string {
	return c.tasks
}

func (c *client) GetApps() []string {
	return c.apps
}

func (c *client) AppendApps(appid string) {
	c.apps = append(c.apps, appid)
}

func (c *client) UpdateAllocation(req map[string]*spincomm.ResourceRequirement) {
	for key, value := range req {
		c.status.HostResource[key].Assigned += value.Requested
		c.status.HostResource[key].Unassigned = c.status.HostResource[key].Total - c.status.HostResource[key].Assigned
	}
	return
}