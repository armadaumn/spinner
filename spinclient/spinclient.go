package spinclient

import (
	"github.com/aramdanet/spinresp"
	"errors"
)

type client struct {
	id			string
	stream 		spinresp.Spinner_AttachServer
	info        Node
}

type Client interface {
	Id() string 	
}

type Node struct {
	openPorts   []string
	status      spinresp.NodeInfo
}

func RequestClient(request *spinresp.JoinRequest, stream spinresp.Spinner_AttachServer) (*Client, error) {
	client := &Client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
	}
	if !client.id {
		return nil, errors.New("No Client ID given")
	}
	return client, nil
}

func (c *client) Id() string {
	return c.id
}