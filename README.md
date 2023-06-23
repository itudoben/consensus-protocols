# consensus-protocols

A Go project to learn Go and consensus protocols.

https://pkg.go.dev allows to find packages from other modules.

# Dev Container

All the work is done through a development container. Docker desktop or other docker engine must be running.
The container is started as follows:

```bash
./devcontainer run
```

The source code can be executed as follows:
```bash
./devcontainer build
```

# Format code 
gofmt -w file_name
https://www.tothenew.com/blog/gofmt-formatting-the-go-code/#:~:text=commands%20and%20options,to%20the%20source%20before%20reformatting.

# Issues on Fri, June 23, 2023
using document raft-consensus-protocol-thesis.pdf stored on PC

From a third container to talk to 2 others node
curl 172.17.0.2:8000/status

Next:
- store the state from node.go
- get the packages in directory main and state working with the imports in node.go
