DEV_CONTAINER := spinresp

protobuf: protobuf_build
	cd .. && docker run -it -v `pwd`:/home/dev $(DEV_CONTAINER) protoc -I $(DEV_CONTAINER)/ $(DEV_CONTAINER)/$(DEV_CONTAINER).proto --go_out=plugins=grpc:$(DEV_CONTAINER) --go_opt=paths=source_relative

protobuf_build:
	cd .. && docker build -t $(DEV_CONTAINER) -f build/Dockerfile.dev .

test: build 
	# cd .. && docker run -t -v `pwd`:/home/dev $(DEV_CONTAINER) go test -v . 
	cd .. && docker run -t -v `pwd`:/home/dev $(DEV_CONTAINER) go test -v ./spinclient

build:
	cd .. && docker build -t $(DEV_CONTAINER) -f build/Dockerfile.dev . 