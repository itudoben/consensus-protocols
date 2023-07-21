# node

This module provides code to run nodes that implement the Raft protocol to form a cluster.

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

# Nodes Communication

## REST API

curl 172.17.0.2:8000/status

## Broadcast with udp

echo i | nc -bu -w 1 172.17.255.255 8872 // to print the IP on the server logs
echo q | nc -bu -w 1 172.17.255.255 8872 // to quit the app

# Status on Fri, July 21, 2023

using document raft-consensus-protocol-thesis.pdf stored on PC

From a third container to talk to 2 others node by broadcasting a command
or
Next:

- Communicate securely by providing a public key
- Check ES how it's done.
- 