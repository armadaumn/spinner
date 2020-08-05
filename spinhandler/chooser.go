package spinhandler

import "github.com/armadanet/spinner/spincomm"

type Chooser interface {
	F(ClientMap, *spincomm.TaskRequest) (string, error)
}