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
	ctx			context.Context
}

type Client interface {
	Id() string
	SendTask(task *spincomm.TaskRequest) error
	Run() error 
}

func RequestClient(ctx context.Context, request *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) (Client, error) {
	ctx, cancel := context.WithCancel(ctx)
	c := &client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
		taskchan: make(chan *spincomm.TaskRequest),
		cancel: cancel,
		ctx: ctx, 
	}
	if c.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
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