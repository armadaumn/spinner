package main

import (
  "github.com/armadanet/spinner"
  "log"
)


func main() {
  err := spinner.CreateAndServe()
  if err != nil {log.Fatalln(err)}
}


// beaconURL := os.Getenv("URL")
// spinnerId := os.Getenv("SPINNERID")
// sp, err := spinner.New(spinnerId)
// if err != nil {panic(err)}
// sp.Run(beaconURL, 5912)