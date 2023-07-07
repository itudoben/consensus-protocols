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

Format all in the current directory
gofmt -w .

# Issues on Fri, July 7, 2023
using document raft-consensus-protocol-thesis.pdf stored on PC

From a third container to talk to 2 others node by broadcasting a command
curl 172.17.0.2:8000/status
or
echo q | nc -bu -w 1 172.17.255.255 8872 // to quit the app
echo i | nc -bu -w 1 172.17.255.255 8872 // to print the IP on the server logs

Next:
- Communicate securely by providing a public key
- Check ES how it's done.
- 