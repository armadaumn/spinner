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