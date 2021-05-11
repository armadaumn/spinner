// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
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

func CreateAndServe(regsistry string, beaconURL string) error {
  ctx := context.Background()
  return GracefulListen(ctx,  5912, regsistry, beaconURL)
}

func GracefulListen(ctx context.Context, port int, registryURL string, beaconURL string) error {
  ctx, cancel := context.WithCancel(ctx)
  defer cancel()
  server := spinserver.New(ctx, registryURL)
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