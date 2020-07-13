package spinclient

import (
	"github.com/armadanet/spinner/spinresp"
)

type client struct {
	id			string
	stream 		spinresp.Spinner_AttachServer
}

type Client interface {
	Id() string 	
}

func RequestClient(request *spinresp.JoinRequest, stream spinresp.Spinner_AttachServer) (Client, error) {
	client := &client{
		id: request.GetCaptainId().GetValue(),
		stream: stream, 
	}
	if client.id == "" {
		return nil, &MalformedClientRequestError{
			err: "No Client ID given",
		}
	}
	return client, nil
}

func (c *client) Id() string {
	return c.id
}