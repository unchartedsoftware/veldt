# Set the base image to ubuntu:trusty
FROM ubuntu:trusty

# Copy the server build files to the container's workspace
ADD ./server.bin /go/src/github.com/unchartedsoftware/prism/server.bin
ADD ./public /go/src/github.com/unchartedsoftware/prism/public

# Run the server by default when the container starts
ENTRYPOINT /go/src/github.com/unchartedsoftware/prism/server.bin \
    -alias="openstack=http://10.64.16.120:9200" \
    -port="8080" \
    -public="/go/src/github.com/unchartedsoftware/prism/public" \
    -redisHost=$PRISM_REDIS_1_PORT_6379_TCP_ADDR \
    -redisPort=$PRISM_REDIS_1_PORT_6379_TCP_PORT

# Document that the service listens on port 8080
EXPOSE 8080
