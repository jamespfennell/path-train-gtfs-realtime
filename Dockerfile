FROM golang

# Copy the local package files to the container's workspace.
ADD path-train-gtfs-realtime.go /go/src/github.com/jamespfennell/path-train-gtfs-realtime/

# Dependencies
RUN go get github.com/google/uuid
RUN go get github.com/golang/protobuf/proto
RUN go get github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs

# Compile
RUN go install github.com/jamespfennell/path-train-gtfs-realtime

RUN mkdir /output
ENV PATH_GTFS_REALTIME_OUTPUT_PATH=/output/gtfsrt

# Run the command by default when the container starts.
ENTRYPOINT /go/bin/path-train-gtfs-realtime
