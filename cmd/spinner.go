package main

import (
  "github.com/armadanet/spinner"
  "log"
  "os"
)


func main() {
  registryURL := ""
  beaconURL := ""
  if len(os.Args) > 1 {
    registryURL = os.Args[1]
    if len(os.Args) > 2 {
      beaconURL = os.Args[2]
    }
  }
  err := spinner.CreateAndServe(registryURL, beaconURL)
  if err != nil {log.Fatalln(err)}
}


// beaconURL := os.Getenv("URL")
// spinnerId := os.Getenv("SPINNERID")
// sp, err := spinner.New(spinnerId)
// if err != nil {panic(err)}
// sp.Run(beaconURL, 5912)