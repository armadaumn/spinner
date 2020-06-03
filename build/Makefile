IMAGE_NAME := spinner
CONTAINER_NAME := spinnerid
URL := http://3.23.114.100:9898/newSpinner
SELFSPIN := true

build:
	cd .. && docker build -t $(IMAGE_NAME) -f build/Dockerfile .


run: build
	docker run -it --name $(CONTAINER_NAME) -e URL=$(URL) -e SPINNERID=$(CONTAINER_NAME) -e SELFSPIN=$(SELFSPIN) -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME)

clean:
		docker rm $(CONTAINER_NAME)
		docker rmi $(IMAGE_NAME)
		docker swarm leave -f
