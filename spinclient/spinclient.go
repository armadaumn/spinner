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
	info     nodeInfo
	ip       string
	port     string
	lat      float64 //latitude
	lon      float64 //longitude
	geoid    string
	tasks    []string // existing tasks
	apps     []string // existring apps
}

type request struct {
	Requirement *spincomm.TaskRequest
	Stream      spincomm.Spinner_RequestServer
}

type Client interface {
	Id() string
	SendTask(task *spincomm.TaskRequest) error
	Run() error
	Info() nodeInfo
	UpdateStatus(status *spincomm.NodeInfo) error
	Location() (float64, float64, error)
	IP() string
	Geoid() string
	GetTasks() []string
	GetApps() []string
	AppendApps(appid string)
}

func RequestClient(ctx context.Context, request *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) (Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
		taskchan: make(chan *spincomm.TaskRequest),
		cancel: cancel,
		ctx: ctx,
		info: nodeInfo{
			HostResource: make(map[string]*spincomm.ResourceStatus),
			UsedPorts: make(map[string]string),
		},
		ip: request.GetIP(),
		port: request.GetPort(),
		lat: request.GetLat(),
		lon: request.GetLon(),
		tasks: make([]string, 0),
		apps: make([]string, 0),
	}
	//c.ip = "0.0.0.0"
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}

	if c.lat == 0 && c.lon == 0 {
		// TODO: fetch lat and lon
		c.lat = 45.0196
		c.lon = -93.2402
	}
	c.geoid = c.genGeoHashID(c.lat, c.lon)

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

func (c *client) Info() nodeInfo {
	return c.info
}

func (c *client) UpdateStatus(status *spincomm.NodeInfo) error {
	//TODO: Remove testing output information
	c.info.UsedPorts = status.GetUsedPorts()
	c.info.ActiveContainer = status.GetContainerStatus().GetActiveContainer()
	c.info.Images = status.GetContainerStatus().GetImages()
	c.info.HostResource = status.GetHostResource()
	//c.apps = status.GetAppIDs()
	c.tasks = status.GetTaskIDs()
	log.Println("after:", c.info)
	//log.Println("app: ", c.apps)
	return nil
}

func (c *client) Location() (float64, float64, error) {
	if c.lat == 0 && c.lon == 0 {
		return 0, 0, errors.New("No availabe location")
	}
	return c.lat, c.lon, nil
}

func (c *client) IP() string {
	return c.ip
}

func (c *client) Geoid() string {
	return c.geoid
}

func (c *client) genGeoHashID(lat float64, lon float64) string {
	geohashIDstr := geohash.Encode(lat, lon)
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