# Set the base image to redis
FROM redis

# Copy the config file to the container
COPY ./redis.conf /usr/local/etc/redis/redis.conf

# Start the server with the conf file
CMD [ "redis-server", "/usr/local/etc/redis/redis.conf" ]
