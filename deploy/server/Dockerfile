# Set the base image to golang
FROM golang

# Copy the server build files to the container's workspace.
ADD ./bin /go/src/github.com/unchartedsoftware/prism/bin
ADD ./public /go/src/github.com/unchartedsoftware/prism/public

# Run the server by default when the container starts.
ENTRYPOINT /go/src/github.com/unchartedsoftware/prism/bin/server.bin

# Document that the service listens on port 8080.
EXPOSE 8080
