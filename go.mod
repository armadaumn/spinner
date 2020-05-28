module github.com/armadanet/spinner

go 1.13

//replace github.com/armadanet/comms => /Users/lh/Desktop/go/src/github.com/armadanet/comms

//replace github.com/armadanet/captain/dockercntrl => /Users/lh/Desktop/go/src/github.com/armadanet/captain/dockercntrl

require (
	github.com/armadanet/captain/dockercntrl v0.0.0-20200528074631-abe3500b4269
	github.com/armadanet/comms v0.0.0-20200528090635-3089ce4375d7
	github.com/armadanet/spinner/spinresp v0.0.0-20200130235212-5ec32922cd99
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/phayes/freeport v0.0.0-20180830031419-95f893ade6f2
)
