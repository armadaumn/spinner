package spinhandler

import (
	"github.com/ArmadaStore/comms/rpc/taskToCargoMgr"
	"github.com/armadanet/spinner/spincomm"
)

type Chooser interface {
	F(ClientMap, *spincomm.TaskRequest) (string, *taskToCargoMgr.Cargos, error)
}