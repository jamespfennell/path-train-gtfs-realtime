# The first step is to retrieve the proto files for the GTFS Realtime spec and the
# source gRPC API. The protoc Docker image used to build these is based on Alpine
# Linux and doesn't have the necessary facilities (git, make) to get the protos.
# This is why we do it in a separate stage.
FROM buildpack-deps:buster AS get-protos

WORKDIR /build
ADD Makefile .
RUN make get-protos

FROM namely/grpc-cli as build-protos

COPY --from=get-protos /build /build
WORKDIR /build
RUN protoc --go_out=./gtfsrt  --proto_path=./gtfsrt ./gtfsrt/*.proto
RUN protoc --go_out=plugins=grpc:./sourceapi  --proto_path=/opt/include --proto_path=./sourceapi ./sourceapi/*.proto

FROM golang:1.14 AS build-go

# Dependencies. These come first to take advantage of Docker caching.
ENV GO111MODULE=on
WORKDIR /build
COPY go.mod .
COPY go.sum .
RUN go mod download

# We're intentionally only copying over the files that are also kept in source control
# so that the Docker build replicates the bare metal build.
COPY --from=build-protos /build/gtfsrt/* /build/gtfsrt/
COPY --from=build-protos /build/sourceapi/* /build/sourceapi/

COPY . .
RUN go build cmd/path_gtfsrt.go

# As is standard, we copy over the built binary to its own Docker image so the
# resulting image does not have redundant Go build infrastructure and artifacts.
FROM debian:buster

COPY --from=build-go /build/path_gtfsrt /usr/local/bin/

RUN mkdir /output
ENV PATH_GTFS_RT_OUTPUT_PATH=/output/gtfsrt

ENTRYPOINT ["path_gtfsrt"]
