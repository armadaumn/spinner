module github.com/armadanet/spinner

go 1.13

//replace github.com/armadanet/comms => /Users/lh/Desktop/go/src/github.com/armadanet/comms

//replace github.com/armadanet/captain/dockercntrl => /Users/lh/Desktop/go/src/github.com/armadanet/captain/dockercntrl
replace github.com/armadanet/spinner => /Users/zhiyingliang/Documents/armada/spinner

//replace github.com/ArmadaStore/comms => /Users/zhiyingliang/Documents/armada/store/comms

require (
	github.com/ArmadaStore/comms v0.0.0-20210108195012-cc92ea560945
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.1.2
	github.com/mmcloughlin/geohash v0.10.0
	github.com/stretchr/stew v0.0.0-20130812190256-80ef0842b48b
	golang.org/x/net v0.0.0-20200114155413-6afb5195e5aa // indirect
	golang.org/x/sync v0.0.0-20190423024810-112230192c58
	golang.org/x/sys v0.0.0-20190507160741-ecd444e8653b // indirect
	google.golang.org/grpc v1.34.0
	google.golang.org/protobuf v1.25.0
)
