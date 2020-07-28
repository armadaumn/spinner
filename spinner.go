// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
  // "google.golang.org/grpc"
  "github.com/armadanet/spinner/spinresp"
  "github.com/armadanet/spinner/spinhandler"
)

type spinnerserver struct {
  spinresp.UnimplementedSpinnerServer
  handler     *spinhandler.Handler
}

func New() spinresp.SpinnerServer {
  return &spinnerserver{
    handler: spinhandler.New(),
  }
}

func (s *spinnerserver) Attach(req *spinresp.JoinRequest, stream spinresp.Spinner_AttachServer) error {
  if err := s.handler.AddClient(req, stream); err != nil {
    return err
  }
  return nil
}

func (s *spinnerserver) ReportStatus(ctx context.Context, req *pb.NodeInfo) (*pb.PingResp, error) {
  resp := pb.PingResp{}
  if err := s.handler.UpdateClient(req); err != nil {
    resp.value = false
    return resp, err
  } 
  resp.value = true
  return resp, nil
}