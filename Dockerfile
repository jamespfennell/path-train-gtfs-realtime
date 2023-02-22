FROM golang:1.19 AS build

WORKDIR /build

COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Install all the code generation tools.
RUN go install \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc
RUN curl -sSL "https://github.com/bufbuild/buf/releases/download/v1.13.1/buf-$(uname -s)-$(uname -m)" \
    -o "/usr/bin/buf"
RUN chmod +x "/usr/bin/buf"

COPY . .

RUN cd proto/gtfsrt && buf generate
RUN cd proto/sourceapi && buf generate
RUN go build cmd/pathgtfsrt.go
RUN go test ./...

FROM debian:buster
COPY --from=build /build/pathgtfsrt /usr/local/bin/
ENTRYPOINT ["pathgtfsrt"]
