package spinclient

import (
	"github.com/armadanet/spinner/spincomm"
	"context"
	"log"
	"errors"
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
}

type Client interface {
	Id() string
	SendTask(task *spincomm.TaskRequest) error
	Run() error
	Info() nodeInfo
	UpdateStatus(status *spincomm.NodeInfo) error
	Location() (float64, float64, error)
	IP() string
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
		},
		ip: request.GetIP(),
		port: request.GetPort(),
		lat: request.GetLat(),
		lon: request.GetLon(),
	}
	c.ip = "0.0.0.0"
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}

	if c.lat == 0 && c.lon == 0 && c.ip != "" {
		// TODO: fetch lat and lon
		c.lat = 45.0196
		c.lon = -93.2402
	} else {
		c.lat = 45.0196
		c.lon = -93.2402
	}

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
		case <- ctx.Done():
			log.Printf("Context ended for %s: %v\n", c.Id(), ctx.Err())
			c.err = ctx.Err()
			return c.err 
		}
	}
}

func (c *client) Info() nodeInfo {
	return c.info
}

func (c *client) UpdateStatus(status *spincomm.NodeInfo) error {
	//TODO: Remove testing output information
	//log.Println("before:", c.info)
	c.info.UsedPorts = status.GetUsedPorts()
	c.info.ActiveContainer = status.GetContainerStatus().GetActiveContainer()
	c.info.Images = status.GetContainerStatus().GetImages()
	c.info.HostResource = status.GetHostResource()
	log.Println("after:", c.info)
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