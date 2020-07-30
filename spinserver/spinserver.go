package spinserver


import (
	"google.golang.org/grpc"
	"github.com/armadanet/spinner/spincomm"
	"github.com/armadanet/spinner/spinhandler"
	"context"
	"log"
	"time"
	"strconv"
  )


type spinnerserver struct {
	spincomm.UnimplementedSpinnerServer
	handler     spinhandler.Handler
  }
  
  func New() *grpc.Server {
	s := &spinnerserver{
	  handler: spinhandler.New(),
	}
	grpcServer := grpc.NewServer()
	spincomm.RegisterSpinnerServer(grpcServer, s)
	return grpcServer
  }

  func (s *spinnerserver) Attach(req *spincomm.JoinRequest, stream spincomm.Spinner_AttachServer) error {
	log.Println("Attaching")
	ctx := context.Background()
	if err := s.handler.AddClient(ctx, req, stream); err != nil {
		log.Fatalln(err)
		return err
	}
	log.Println(s.handler.ListClientIds())
	cl, ok := s.handler.GetClient("captain")
	if !ok {
		log.Println("No captain")
		return nil
	}
	i := 1
	for {
		t := &spincomm.TaskRequest{
			TaskId: &spincomm.UUID{
				Value: "simple_task_" + strconv.Itoa(i),
			},
		}
		err := cl.SendTask(t)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println("Task sent:", i)
		time.Sleep(2*time.Second)
		i += 1
	}
	return nil
  }