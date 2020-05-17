# Bast Cloud Controller

This is the code for the Bast Cloud Controller.

The Controller is responsible for keeping track of user settings and authentication.

# Running

You can use the Dockerfiles to build and run in containers:

```sh
# start the service
# (or use docker-compose)
podman-compose up

# stop the service
podman-compose down
```
Alternatively, you can build from source and run it.

## Build and Run

After installing the prerequisites:

* Golang
* MariaDB or MySQL

It can be started by running the go module with:

```sh
go build
./controller
```

