package spinclient

import (
	"github.com/armadanet/spinner/spincomm"
	"context"
	"log"
)

type client struct {
	id			string
	stream 		spincomm.Spinner_AttachServer
	cancel		func()
	taskchan	chan *spincomm.TaskRequest
	err			error
}

type Client interface {
	Id() string
	SendTask(task *spincomm.TaskRequest) error
}

func RequestClient(ctx context.Context, request *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) (Client, error) {
	c := &client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
	}
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}
	go c.run(ctx)
	return c, nil
}

func (c *client) Id() string {
	return c.id
}

func (c *client) SendTask(task *spincomm.TaskRequest) error {
	if c.err != nil {
		return c.err
	}
	return c.stream.Send(task)
}

func (c *client) run(ctx context.Context) {
	if c.cancel != nil {c.cancel()}
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()
	c.cancel := cancel
	c.taskchan := make(chan *spincomm.TaskRequest)

	for {
		select {
		case task, ok := <- c.taskchan:
			if !ok {
				log.Printf("Task Channel closed for %s\n.", c.Id())
				return
			}
			if err := c.stream.Send(task); err != nil {
				log.Printf("Error for %s: %v\n", c.Id(), err)
				return
			}
		case <- ctx.Done():
			log.Printf("Context ended for %s: %v\n", c.Id(), ctx.Err())
			return
		}
	}
}