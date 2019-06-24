FROM golang:1.11

LABEL maintainer="bitwangtuo@gmail.com"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/wangthomas/bloomfield

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Download dependencies
RUN go get -d -v ./...

# Install the package
RUN go install -v ./...

# This container exposes port 8080 to the outside world
EXPOSE 8679

CMD ["bloomfield", "-config_file=config.toml"]