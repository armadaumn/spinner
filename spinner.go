// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
  "google.golang.org/grpc"
  "github.com/armadanet/spinner/spinserver"
  "net"
  "log"
  "strconv"
  "context"
  "os"
  "os/signal"
  "syscall"
  "golang.org/x/sync/errgroup"
)

func CreateAndServe() error {
  server := spinserver.New()
  ctx := context.Background()
  return GracefulListen(ctx, server, 5912)
}

func GracefulListen(ctx context.Context, server *grpc.Server, port int) error {
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()
  g, ctx := errgroup.WithContext(ctx)

  interrupt := make(chan os.Signal, 1)
  signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
  defer signal.Stop(interrupt)

  g.Go(func() error {
    portVal := ":" + strconv.Itoa(port)
    lis, err := net.Listen("tcp", portVal)
    if err != nil {return err}
    log.Printf("Listening on TCP port %d\n", port)

    return server.Serve(lis)
  })
  
  select {
  case <-interrupt:
    break
  case <-ctx.Done():
    break
  }

  log.Println("Shutting Down Server")
  cancel()
  server.GracefulStop()

  err := g.Wait()
  return err
}

