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

func TestWellformedClient(t *testing.T) {
	req := &spinresp.JoinRequest{
		CaptainId: &spinresp.UUID{
			Value: "fake_id",
		},
	}
	resp, err := spinclient.RequestClient(req, nil);
	if  err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if resp.Id() != "fake_id" {
		t.Errorf("ID should be 'fake_id', not '%v'", resp.Id())
	}
}