DEV_CONTAINER := spincomm
DEV_DOCKERFILE := Dockerfile.dev
CONTAINER := spinner
DOCKERFILE := Dockerfile
NETWORK := spinner-local-network

protobuf: protobuf_build
	cd .. && docker run -it -v `pwd`:/home/dev $(DEV_CONTAINER) protoc -I $(DEV_CONTAINER)/ $(DEV_CONTAINER)/$(DEV_CONTAINER).proto --go_out=plugins=grpc:$(DEV_CONTAINER) --go_opt=paths=source_relative

protobuf_build:
	cd .. && docker build -t $(DEV_CONTAINER) -f build/$(DEV_DOCKERFILE) .

test: test_build 
	# cd .. && docker run -t -v `pwd`:/home/dev $(DEV_CONTAINER) go test -v . 
	cd .. && docker run -t -v `pwd`:/home/dev $(DEV_CONTAINER) go test -v ./spinclient

test_build:
	cd .. && docker build -t $(DEV_CONTAINER) -f build/$(DEV_DOCKERFILE) . 

run: build network
	docker run -it --rm -p 5912:5912 --net $(NETWORK) --name $(CONTAINER) $(CONTAINER)

build: test_build
	cd .. && docker build -t $(CONTAINER) -f build/$(DOCKERFILE) --build-arg base=$(DEV_CONTAINER) .

network:
	docker network create $(NETWORK) || true

clean:
	docker rm -f $(CONTAINER) || true
	docker network rm $(NETWORK) || true
