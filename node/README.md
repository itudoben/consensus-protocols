# node

This module provides code to run server that implements the Raft protocol to form a cluster.

Here are the steps to test a cluster:
- build the dev container image
  - ./devcontainer builddev
- start a dev container
  - ./devcontainer dev
- at the prompt cd into ./node
- build the node source code
  - go build . or go build node.go
- start the node
  - ./node


Start another dev container
- cd /Users/jhujol/Projects/itudoben/consensus-protocols
- ./devcontainer dev
- check the IP of the first node
- execute the command
- curl http://172.17.0.2:8000/status