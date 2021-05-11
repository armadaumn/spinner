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
	"strings"

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
	taskMap     map[string][]*spincomm.TaskRequest
	registry    spinhandler.Registry
}
  
func New(ctx context.Context, registryURL string) *grpc.Server {
	chooser := spinhandler.InitCustomChooser()
	s := &spinnerserver{
	  handler: spinhandler.New(),
	  ctx: ctx,
	  chooser: &chooser,
	  router: make(map[string]*spinrequest),
	  taskMap: make(map[string][]*spincomm.TaskRequest),
	}
	if registryURL != "" {
		s.registry = spinhandler.NewRegistry(registryURL)
		go s.registry.UpdateImageList()
	}

	grpcServer := grpc.NewServer()
	spincomm.RegisterSpinnerServer(grpcServer, s)
	return grpcServer
}

// Receive a new task
func (s *spinnerserver) Request(req *spincomm.TaskRequest, stream spincomm.Spinner_RequestServer) error {
	//log.Printf("Got request: %v\n", req)
	taskID := req.GetTaskId().GetValue()
	if taskID == "" {
		return errors.New("No task id given")
	}
	ctx, cancel := context.WithCancel(s.ctx)
	request := NewRequest(stream, cancel)
	s.router[taskID] = request

	// Choose captains
	cl, cargo, err := s.handler.ChooseClient(s.chooser, req)
	if err != nil {return err}
	//cl, ok := s.handler.GetClient(cid)
	//if !ok {return errors.New("No such client")}

	// If the local registry has the same image, pull it from the local one (Note that maybe slower than docker hub)
	image := req.GetImage()
	imageName := strings.Split(image, "/")
	localRepos := s.registry.GetRepos()
	log.Println(localRepos)
	if _, ok := localRepos[imageName[2]]; ok {
		localName := s.registry.GetUrl() + "/" + imageName[2]
		req.Image = localName
	}

	if cargo != nil {
		req.Taskspec.CargoSpec.IPs = cargo.GetIPs()
		req.Taskspec.CargoSpec.Ports = cargo.GetPorts()
	}
	if err = cl.SendTask(req); err != nil {return err}

	if _, ok := s.taskMap[cl.Id()]; !ok {
		s.taskMap[cl.Id()] = []*spincomm.TaskRequest{req}
	} else {
		s.taskMap[cl.Id()] = append(s.taskMap[cl.Id()], req)
	}

	// Return selected node info
	//taskLog := spincomm.TaskLog{
	//	TaskId: &spincomm.UUID{Value: taskID},
	//	Ip: cl.IP(),
	//	HostResource: cl.NodeStatus().HostResource,
	//}
	//if err = stream.Send(&taskLog); err != nil {
	//	log.Println(err)
	//	return err
	//}

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
	if err == nil || err.Error() != "context ended" {
		// handling error
		log.Println(cl.Id() + " left")

		s.handler.RemoveClient(cl.Id())
		if taskList, ok := s.taskMap[cl.Id()]; ok {
			for _, task := range taskList {
				taskid := task.GetTaskId().GetValue()
				s.router[taskid].complete()
			}
			delete(s.taskMap, cl.Id())

			// Restart deployment
			//amStream := s.router[task.GetTaskId().GetValue()]
			//delete(s.taskMap, cl.Id())
			//retry := 0
			//for true {
			//	err := s.Request(task, amStream.stream)
			//	retry++
			//	if err == nil || retry == 1{
			//		break
			//	}
			//	time.Sleep(10 * time.Second)
			//}
		}
	}

	return err
}

// Captain status update
func (s *spinnerserver) Update(ctx context.Context, status *spincomm.NodeInfo) (*spincomm.PingResp, error) {
	err := s.handler.UpdateClient(status)
	res := spincomm.PingResp{
		Status: true,
	}

	cid := status.GetCaptainId().GetValue()
	if taskList, ok := s.taskMap[cid]; ok {
		for _, task := range taskList {
			if _, ok := status.UsedPorts[task.GetTaskId().Value]; ok {
				go s.ReportTask(task.GetTaskId().GetValue(), cid, status)
			}
		}
	}
	if err != nil {
		res.Status = false
		return &res, err
	}
	return &res, nil
}

func (s *spinnerserver) ReportTask(taskID string, cid string, status *spincomm.NodeInfo) {
	stream := s.router[taskID].stream
	cl, _ := s.handler.GetClient(cid)
	taskLog := spincomm.TaskLog{
		TaskId: &spincomm.UUID{Value: taskID},
		Ip: cl.NodeInfo().Ip,
		Port: status.UsedPorts[taskID],
		HostResource: status.HostResource,
		Location: &spincomm.Location{Lat: cl.NodeInfo().Lat, Lon: cl.NodeInfo().Lon},
		Tag: cl.NodeInfo().Tags,
	}
	if cl.NodeInfo().ServerType == spincomm.Type_LocalServer {
		taskLog.NodeType = 2
	} else {
		taskLog.NodeType = 3
	}

	if err := stream.Send(&taskLog); err != nil {
		log.Println(err)
	}
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

//TODO: set a new registry url