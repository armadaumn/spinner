## What is this?
Spinner is a computation resource manager in the project, [Armada](https://armadanet.github.io). It manages a group of worker nodes, [captain](https://github.com/armadanet/captain), in the system. 
It is responsible for node registration, node management, and task instance placement.

## Quick Start
**Prerequisites**: Docker

To download the spinner image run: \
`docker pull armadaumn/spinner:latest` \
To start the spinner just run: \
`docker run -it --rm -p 5912:5912 armadaumn/spinner:latest`

## Build from the source
**Prerequisites**: Go environment, open port 5912

Build and run spinner:
```
git clone https://github.com/armadanet/spinner.git
cd spinner/build
make run
```
