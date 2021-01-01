module github.com/armadanet/spinner

go 1.13

//replace github.com/armadanet/comms => /Users/lh/Desktop/go/src/github.com/armadanet/comms

//replace github.com/armadanet/captain/dockercntrl => /Users/lh/Desktop/go/src/github.com/armadanet/captain/dockercntrl
replace github.com/armadanet/spinner => /Users/zhiyingliang/Documents/armada/spinner

require (
	github.com/ArmadaStore/comms v0.0.0-20201231053020-c7ef8cc8487b
	github.com/golang/protobuf v1.4.3
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190507160741-ecd444e8653b // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)
