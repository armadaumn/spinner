package spinclient

import (
	"github.com/armadanet/spinner/spincomm"
	"context"
	"log"
	"errors"
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
		err: errors.New("Not yet initialized"),
	}
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}
	completed := make(chan bool)
	go c.run(ctx, completed)
	<- completed

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

func (c *client) run(ctx context.Context, completed chan bool) {
	if c.cancel != nil {c.cancel()}
	ctx, cancel := context.WithCancel(ctx)

	defer cancel()
	c.cancel = cancel
	c.taskchan = make(chan *spincomm.TaskRequest, 2)

	log.Println("Ready for tasks")
	c.err = nil
	close(completed)

	for {
		select {
		case task, ok := <- c.taskchan:
			log.Println("Task sent")
			if !ok {
				log.Printf("Task Channel closed for %s\n.", c.Id())
				c.err = errors.New("Task Channel closed")
				return
			}
			if err := c.stream.Send(task); err != nil {
				log.Printf("Error for %s: %v\n", c.Id(), err)
				c.err = err
				return
			}
		case <- ctx.Done():
			log.Printf("Context ended for %s: %v\n", c.Id(), ctx.Err())
			c.err = ctx.Err()
			return
		}
	}
}