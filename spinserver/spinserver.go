package spinserver


import (
	"google.golang.org/grpc"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler"
	"github.com/armadanet/spinner/spinclient"
	"context"
	"log"
	"io"
	"errors"
	// "time"
	// "strconv"
)


type spinnerserver struct {
	spincomm.UnimplementedSpinnerServer
	handler     spinhandler.Handler
	ctx			context.Context
	chooser		spinhandler.Chooser
	router		map[string]spincomm.Spinner_RequestServer
  }
  
func New(ctx context.Context) *grpc.Server {
	s := &spinnerserver{
	  handler: spinhandler.New(),
	  ctx: ctx, 
	  chooser: &spinhandler.RoundRobinChooser{
		  LastChoice: "",
	  },
	}
	grpcServer := grpc.NewServer()
	spincomm.RegisterSpinnerServer(grpcServer, s)
	return grpcServer
}

func (s *spinnerserver) Request(req *spincomm.TaskRequest, stream spincomm.Spinner_RequestServer) error {
	id := req.GetTaskId().GetValue()
	if id == "" {
		return errors.New("No task id given")
	}
	s.router[id] = stream
	cid, err := s.handler.ChooseClient(s.chooser)
	if err != nil {return err}
	cl, ok := s.handler.GetClient(cid)
	if !ok {return errors.New("No such client")}
	return cl.SendTask(req)
}

func (s *spinnerserver) Run(stream spincomm.Spinner_RunServer) error {
	for {
		taskLog, err := stream.Recv()
		if err == io.EOF {return stream.SendAndClose(&spincomm.TaskCompletion{})}
		if err != nil {return err}
		log.Println(taskLog)
		id := taskLog.GetTaskId().GetValue()
		if id == "" {return errors.New("No id given")}
		taskStream, ok := s.router[id]
		if !ok {return errors.New("No such task")}
		if err = taskStream.Send(taskLog); err != nil {return err}
	}
}

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
	log.Println(err)
	return err
}

func (s *spinnerserver) Update(ctx context.Context, status *spincomm.NodeInfo) (*spincomm.PingResp, error) {
	err := s.handler.UpdateCLient(status)
	resp := spincomm.PingResp{
		Status: true,
	}
	if err != nil {
		resp.Status = false
		return &resp, err
	}
	return &resp, nil
}

