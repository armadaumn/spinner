package spinclient_test

import (
	"testing"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spinresp"
)

func TestMalformedClient(t *testing.T) {
	req := &spinresp.JoinRequest{
		CaptainId: nil,
	}
	_, err := spinclient.RequestClient(req, nil)
	if err == nil {
		t.Error("No error given")
	} else {
		switch v := err.(type) {
		case *spinclient.MalformedClientRequestError:
		default:
			t.Errorf("Expected MalformedClientRequestError, not %v", v)
		}
	}
}