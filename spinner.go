// Nebula Spinner server to maintain socket connections to Captains.
package spinner

import (
  "github.com/gorilla/mux"
  "github.com/phayes/freeport"
  "github.com/armadanet/comms"
  "github.com/armadanet/captain/dockercntrl"
  "net/http"
  "log"
  "strconv"
  "fmt"
  "time"
)

// Server for the Nebula Spinner
type Server interface {
  // Given a port of 0, assigns a free port to the server.
  Run(beaconURL string, port int)
}

type server struct {
  router          *mux.Router
  handler         *Handler
  state           *dockercntrl.State
  container_name  string
  overlay_name    string
}

// Produces a new Server interface of struct server
func New(container_name string) (Server, error) {
  router := mux.NewRouter().StrictSlash(true)
  handler := NewHandler()
  state, err := dockercntrl.New()
  if err != nil {return nil, err}
  router.HandleFunc("/join", join(handler)).Name("Join")
  router.HandleFunc("/spin", spin(handler)).Name("Spin")
  handler.Start()
  return &server{
    router: router,
    handler: handler,
    state: state,
    container_name: container_name,
  }, nil
}

type newSpinnerRes struct {
  SwarmToken        string  `json:"SwarmToken"`
  BeaconIp          string  `json:"BeaconIp"`
  BeaconOverlay     string  `json:"BeaconOverlay"`
  BeaconName        string  `json:"BeaconName"`
  SpinnerOverlay    string  `json:"SpinnerOverlay"`
}

// Runs the spinner server.
func (s *server) Run(beaconURL string, port int) {
  // Query beacon
  fmt.Printf("Query beacon /newSpinner...")
  var res newSpinnerRes
  err := comms.SendPostRequest(beaconURL+"/newSpinner", map[string]string{
    "SpinnerId":s.container_name,
  }, &res)
  if err!=nil {
    log.Println(err)
    return
  }
  s.overlay_name = res.SpinnerOverlay
  fmt.Print("Get beacon info: ")
  fmt.Println(res)

  // join beacon swarm and attach self to beacon overlay
  fmt.Println(s.container_name)
  fmt.Println(res.BeaconOverlay)
  err = s.state.JoinSwarmAndOverlay(res.SwarmToken, res.BeaconIp, s.container_name, res.BeaconOverlay)
  if err != nil {
    log.Println(err)
    return
  }

  // attach self to spinner_overlay
  err = s.state.JoinOverlay(s.container_name, res.SpinnerOverlay)
  if err != nil {
    log.Println(err)
    return
  }

  // go routine periodically ping beacon to notify the alive (wait 1s)
  go s.Ping(res.BeaconName)

  // go routine notify parent captain join finish (wait 1s)


  // start the server
  if port == 0 {
    port, err = freeport.GetFreePort()
    if err != nil {log.Println(err); return}
  }
  log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), s.router))
}


func (s *server) Ping(beaconName string) {
  for {
    // err := comms.SendPostRequest("http://localhost:8787/register", map[string]interface{}{
    err := comms.SendPostRequest("http://"+beaconName+":8787/register", map[string]interface{}{
      "Id":s.container_name,
      "OverlayName":s.overlay_name,
      "LastUpdate":time.Now(),
    }, nil)
    if err!=nil {
      panic(err)
      return
    }
    fmt.Println(beaconName)
    // ping every 3 seconds
    time.Sleep(3 * time.Second)
  }
}
