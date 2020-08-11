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
  }
  
func New(ctx context.Context) *grpc.Server {
	s := &spinnerserver{
	  handler: spinhandler.New(),
	  ctx: ctx, 
	  chooser: &spinhandler.RoundRobinChooser{
		  LastChoice: "",
	  },
	  router: make(map[string]*spinrequest),
	}
	grpcServer := grpc.NewServer()
	spincomm.RegisterSpinnerServer(grpcServer, s)
	return grpcServer
}

func (s *spinnerserver) Request(req *spincomm.TaskRequest, stream spincomm.Spinner_RequestServer) error {
	log.Printf("Got request: %v\n", req)
	id := req.GetTaskId().GetValue()
	if id == "" {
		return errors.New("No task id given")
	}
	ctx, cancel := context.WithCancel(s.ctx)
	request := NewRequest(stream, cancel)
	s.router[id] = request
	cid, err := s.handler.ChooseClient(s.chooser)
	if err != nil {return err}
	cl, ok := s.handler.GetClient(cid)
	if !ok {return errors.New("No such client")}
	if err = cl.SendTask(req); err != nil {return err}
	<- ctx.Done()
	return nil
}

func (s *spinnerserver) Run(stream spincomm.Spinner_RunServer) error {
	setEnd := false
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
		taskStream, ok := s.router[id]
		if !ok {
			err = errors.New("No such task")
			log.Println(err)
			return err
		}
		if !setEnd {
			setEnd = true 
			complete := taskStream.complete
			defer complete()
		}
		if err = taskStream.stream.Send(taskLog); err != nil {
			log.Println(err)
			return err
		}
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

