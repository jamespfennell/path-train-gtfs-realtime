# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD path-train-gtfs-realtime.go /go/src/github.com/jamespfennell/path-train-gtfs-realtime/

# Get dependencies
RUN go get github.com/golang/protobuf/proto
RUN go get github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs

# Compile
RUN go install github.com/jamespfennell/path-train-gtfs-realtime

# Run the outyet command by default when the container starts.
ENTRYPOINT /go/bin/path-train-gtfs-realtime
