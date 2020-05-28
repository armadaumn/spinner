IMAGE_NAME := spinner
CONTAINER_NAME := spinneruniqueid

run: build
	docker run -it --name $(CONTAINER_NAME) -e URL=http://localhost:9898 -e SPINNERID=spinnerid -v /var/run/docker.sock:/var/run/docker.sock $(IMAGE_NAME)

build:
	cd .. && docker build -t $(IMAGE_NAME) -f build/Dockerfile .

clean:
		docker rm $(CONTAINER_NAME)
		docker rmi $(IMAGE_NAME)
