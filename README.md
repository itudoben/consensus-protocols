# consensus-protocols

A Go project to learn Go and consensus protocols.

https://pkg.go.dev allows to find packages from other modules.

# Dev Containers

## App Container

All the work is done through a development container. Docker desktop or other docker engine must be running.
The application container is started as follows:

```bash
./devcontainer run
```

The source code can be executed as follows:

```bash
./devcontainer build
```

## Development Container

To develop the [cluster using Raft consensus](./node/README.md) one must use a dev container.
This command builds the dev container:

```bash
./devcontainer builddev
```

This command spawns a container in Docker for development purpose:

```bash
./devcontainer dev
```

# Format code

gofmt -w file_name
https://www.tothenew.com/blog/gofmt-formatting-the-go-code/#:~:text=commands%20and%20options,to%20the%20source%20before%20reformatting.

Format all in the current directory
gofmt -w .
