package spinserver


import (
	"context"
	"errors"
	"github.com/armadanet/spinner/spinclient"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler"
	"google.golang.org/grpc"
	"io"
	"log"
	// "time"
	// "strconv"
)


type spinrequest struct {
	stream 		spincomm.Spinner_RequestServer
	complete	func()
}

func NewRequest(stream spincomm.Spinner_RequestServer, cancel func()) *spinrequest {
	return &spinrequest{
		stream: stream,
		complete: cancel, 
	}
}

type spinnerserver struct {
	spincomm.UnimplementedSpinnerServer
	handler     spinhandler.Handler
	ctx			context.Context
	chooser		spinhandler.Chooser
	router		map[string]*spinrequest
	taskMap     map[string]string
}
  
func New(ctx context.Context) *grpc.Server {
	chooser := spinhandler.InitCustomChooser()
	s := &spinnerserver{
	  handler: spinhandler.New(),
	  ctx: ctx,
	  chooser: &chooser,
	  router: make(map[string]*spinrequest),
	  taskMap: make(map[string]string),
	}
	grpcServer := grpc.NewServer()
	spincomm.RegisterSpinnerServer(grpcServer, s)
	return grpcServer
}

// Receive a new task
func (s *spinnerserver) Request(req *spincomm.TaskRequest, stream spincomm.Spinner_RequestServer) error {
	log.Printf("Got request: %v\n", req)
	taskID := req.GetTaskId().GetValue()
	if taskID == "" {
		return errors.New("No task id given")
	}
	ctx, cancel := context.WithCancel(s.ctx)
	request := NewRequest(stream, cancel)
	s.router[taskID] = request

	// Choose captains
	cid, cargo, err := s.handler.ChooseClient(s.chooser, req)
	if err != nil {return err}
	cl, ok := s.handler.GetClient(cid)
	if !ok {return errors.New("No such client")}
	log.Println(cid)
	log.Println(cargo)
	if cargo != nil {
		req.Taskspec.CargoSpec.IPs = cargo.GetIPs()
		req.Taskspec.CargoSpec.Ports = cargo.GetPorts()
	}
	if err = cl.SendTask(req); err != nil {return err}
	s.taskMap[cid] = taskID

	// Return seleted node info
	taskLog := spincomm.TaskLog{
		TaskId: req.GetTaskId(),
		Ip: cl.IP(),
		HostResource: cl.Info().HostResource,
	}

	if err = stream.Send(&taskLog); err != nil {
		log.Println(err)
		return err
	}

	<- ctx.Done()
	return nil
}

// Runtime task log
func (s *spinnerserver) Run(stream spincomm.Spinner_RunServer) error {
	for {
		taskLog, err := stream.Recv()
		if err == io.EOF {
			log.Println(err)
			return stream.SendAndClose(&spincomm.TaskCompletion{})
		}
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("TaskLog:", taskLog)
		id := taskLog.GetTaskId().GetValue()
		if id == "" {
			err = errors.New("No id given")
			log.Println(err)
			return err
		}
		//taskStream, ok := s.router[id]
		//if !ok {
		//	err = errors.New("No such task")
		//	log.Println(err)
		//	return err
		//}
		//if err = taskStream.stream.Send(taskLog); err != nil {
		//	log.Println(err)
		//	return err
		//}
	}
}

// Captain joins the spinner
func (s *spinnerserver) Attach(req *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) error {
	log.Println("Attaching")
	cl, err := spinclient.RequestClient(s.ctx, req, stream)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	err = s.handler.AddClient(cl)
	if err != nil {
		log.Fatalln(err)
		return err
	}
	
	err = cl.Run()
	//TODO: handling error
	log.Println(err)
	return err
}

// Captain status update
func (s *spinnerserver) Update(ctx context.Context, status *spincomm.NodeInfo) (*spincomm.PingResp, error) {
	err := s.handler.UpdateClient(status)
	res := spincomm.PingResp{
		Status: true,
	}

	go s.ReportTask(status)
	if err != nil {
		res.Status = false
		return &res, err
	}
	return &res, nil
}

// Register a scheduling policy
func (s *spinnerserver) RegisterScheduler(ctx context.Context, sp *spincomm.SchedulePolicy) (*spincomm.PingResp, error) {
	res := spincomm.PingResp{
		Status: true,
	}
	isSuccess := s.chooser.Register(sp.GetId(), sp.GetType(), sp.GetIP(), sp.GetPort())
	if !isSuccess {
		res.Status = false
	}
	return &res, nil
}

func (s *spinnerserver) ReportTask(status *spincomm.NodeInfo) {
	cid := status.GetCaptainId().GetValue()
	taskID, ok := s.taskMap[cid]
	if !ok {
		return
	}

	stream := s.router[taskID].stream
	cl, _ := s.handler.GetClient(cid)
	taskLog := spincomm.TaskLog{
		TaskId: &spincomm.UUID{Value: taskID},
		Ip: cl.IP(),
		HostResource: status.HostResource,
	}

	if err := stream.Send(&taskLog); err != nil {
		log.Println(err)
	}
}