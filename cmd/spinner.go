package main

import (
  "github.com/armadanet/spinner"
  "os"
  //"fmt"
)

// os.Getenv("URL")
// get the env var about where to find beacon
// URL=http://public_ip:port/newSpinner
// PORT=internal open port
func main() {
  // beaconURL := "http://public_ip:9898/newSpinner"
  beaconURL := os.Getenv("URL")
  spinnerId := os.Getenv("SPINNERID")
  sp, err := spinner.New(spinnerId)
  if err != nil {panic(err)}
  sp.Run(beaconURL, 5912)
}
