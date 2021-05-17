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
     beaconURL = os.Args[1]
   }
  }


  log.SetFlags(log.LstdFlags | log.Lmicroseconds)
  //For evaluation, registryURL is used as BeaconURL
  // Spinner will register itself to Beacon
  //go connectBeacon(beaconURL)

  err := spinner.CreateAndServe(registryURL, beaconURL)
  if err != nil {log.Fatalln(err)}
}

//func connectBeacon(beaconURL string) {
//  body := strings.NewReader(`{"SpinnerId":"1","IP":"128.101.118.101","Port":"8081","GeoID":"9zvxy"}`)
//  req, err := http.NewRequest("POST", "http://" + beaconURL + "/newSpinner", body)
//  if err != nil {
//    log.Println(err)
//  }
//  req.Header.Set("Accept", "application/json")
//  req.Header.Set("X-Http-Method-Override", "PUT")
//  req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//  _, err = http.DefaultClient.Do(req)
//}

// beaconURL := os.Getenv("URL")
// spinnerId := os.Getenv("SPINNERID")
// sp, err := spinner.New(spinnerId)
// if err != nil {panic(err)}
// sp.Run(beaconURL, 5912)