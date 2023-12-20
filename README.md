# GTFS Realtime for the PATH Train

This repository hosts a simple Go application
that reads PATH train realtime data from [Matt Razza's public API](https://github.com/mrazza/path-data)
and outputs a feed of the data in the GTFS Realtime format.
Some important notes:

- **You don't need to run the application yourself.**
    The GTFS Realtime feed produced by this software can be accessed at
    [`https://path.transitdata.nyc/gtfsrt`](https://path.transitdata.nyc/gtfsrt).
    It's updated every 5 seconds.

- The outputted data is compatible with [the official GTFS Static data](https://old.panynj.gov/path/developers.html)
    published by the Port Authority
    in the sense that the stop IDs and route IDs match up.
    The feed should work correctly for software that integrates realtime and static data.

- Unfortunately the Port Authority doesn't distribute the full realtime data set, and so the GTFS
  Realtime feed has some big missing pieces:
  - There is no trip data: all the Port Authority communicates are stops, and arrival times at those stops.
    There is no easy way to connect arrival times for the same train at multiple stops.
    So, in the GTFS Realtime feed, the "trips" are dummy trips with a random ID and a single
    stop time update. This should be sufficient for consumers that want to show arrival times at stops,
    but of course prevents other uses like tracking trains through the system.
  - The GTFS Static feed describes all the tracks/platforms at each of the PATH stations
    but in the realtime data we don't known which platform a train will stop at.
    In the realtime feed, all of the trains stop at the "station" stop (i.e., the stop in the static
    feed with location type `1`).


## Running the application

The application is an HTTP server with the
    GTFS Realtime feed available at the `/gtfsrt` path.
In the background, the program periodically retrieves data from the Razza API
    and updates the feed.
By default, this update occurs every 5 seconds.

There are a couple flags that can be passed to the binary:

- `--port <int>`: the port to bind the HTTP server to (default `8080`)

- `--timeout_period <duration>`:
        the maximum duration to wait for a response from the source API (default 5s)

- `--update_period <duration>`:
        how often to update the feed (default 5s).
    Remember that the more frequently you update, the more stress you place
    on the source API, so be nice.

- `--use_http_source_api`
    use the HTTP source API instead of the default gRPC API.

- `--use_panynj_api`:
    use the PANYNJ API.

### Running using Docker

The CI process (using Github actions) builds a Docker image and stores it
at the `jamespfennell/path-train-gtfs-realtime:latest`
[tag on Docker Hub](https://hub.docker.com/repository/docker/jamespfennell/path-train-gtfs-realtime).
You can also build the Docker image locally by running `docker build .` in the
root of the repo.

It is generally simplest to run the application using Docker.
The only thing you need to do is port forward the HTTP server's port outside of the container.
This is a functioning Docker compose configuration that does this:
```
version: '3.5'

services:
  path-train-gtfs-realtime:
    image: jamespfennell/path-train-gtfs-realtime:latest
    port: 8080:9001
    restart: always
```

### Running using `go run`

When doing dev work it is generally necessary to run the application on "bare metal",
which you can do simply with  `go run cmd/pathgtfsrt.go`.

The source gRPC API and the GTFS Realtime format are both built
on `proto` files.
Getting these `proto` files and compiling them to `go` files
is a bit of a pain, so they're kept in source control.
To regenerate them, it's probably just simplest to use the Docker build process.

### Error handling and exit codes

A number of errors can prevent the application from running 100% correctly,
    with the main source of errors being network failures when hitting the source API.
At start-up, the application downloads static and realtime data from the API;
    if this fails, the application will exit.

After start-up, any further errors encountered are handled gracefully,
    and the server will not exit until interrupted.
If, during a particular update, the realtime data for a specific stop cannot be retrieved, or is malformed,
then the previously retrieved data will be used.

### Monitoring

The application exports metrics in Prometheus format on the `/metrics` endpoint.
See `cmd/pathgtfsrt.go` for the metric definitions.

## Licence notes

- All the code in the root directory of the repo is
released under the MIT License (see `LICENSE`).

- The `proto` files in the `sourceapi` directory are sourced from the
[mrazza/path-data Github repo](https://github.com/mrazza/path-data),
are released under the MIT License and are copyright Matthew Razza.

- The `proto` files in the `gtfsrt` directory are sourced from the
[google/tranist Github repo](https://github.com/google/transit),
are released under the Apache License 2.0 and are copyright Google Inc.

- My understanding is that the `proto` copyrights extend
to the compiled `go` files.
