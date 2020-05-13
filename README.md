# GTFS Realtime for the PATH Train

This repository hosts a simple Go application
that reads PATH train realtime data from [Matt Razza's public API](https://github.com/mrazza/path-data)
and outputs the data in the GTFS Realtime format.
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

The application works by periodically retrieving data from the Razza API
and writing the GTFS Realtime file to disk.
If you want to serve the file over the internet, you can use Nginx or a similar web server.
 
There are a couple of environment variable configurations:

- `PATH_GTFS_RT_SOURCE_API` - either `grpc` or `http`;
    denotes which API to use for retrieving the source data. 
    Defaults to the recommended setting `grpc`.
    The `http` option is supported because it's easier to debug.

- `PATH_GTFS_RT_PERIODICITY_MILLISECONDS` - the periodicity, in milliseconds,
    with which to update the GTFS Realtime file. The value of this environment
    variable must be castable to an integer.
    The default is 5000 milliseconds (so 5 seconds).
    Remember that the more frequent you update, the more stress you place
    on the source API, so be nice.
    
- `PATH_GTFS_RT_OUTPUT_PATH` - the location on disk to output the file,
    either absolute or relative to the working directory when the application starts.
    The default is `./path.gtfsrt`.


### Running using Docker

The [Travis CI job](https://travis-ci.org/github/jamespfennell/path-train-gtfs-realtime)
builds a Docker image and stores it
in the `jamespfennell/path-train-gtfs-realtime` 
[repository on Docker Hub](https://hub.docker.com/repository/docker/jamespfennell/path-train-gtfs-realtime).
There are both `latest` tags and `build<n>` tags, where `<n>` is the Travis build number.
You can also build the Docker image locally by running `docker build .` in the
root of the repo.

It is generally simplest to run the application using Docker.
The application in the container writes the GTFS Realtime file to `/output/gtfsrt`.
To access this from outside the container, for example to serve it via Nginx,
just use a Docker volumne.
This is a functioning Docker compose configuration that does this:
```
version: '3.5'

services:
  path-train-gtfs-realtime:
    image: jamespfennell/path-train-gtfs-realtime:latest
    volumes:
      - ./output:/output
    restart: always
```

### Running using `go run`

When doing dev work it is generally necessary to run the application on "bare metal".
The application has a few dependencies which should be installed first:
```
go get github.com/google/uuid
go get github.com/golang/protobuf/proto
go get google.golang.org/grpc
```
After that, just run `go run path_gtfsrt.go`.

The source gRPC API and the GTFS Realtime format are both built
on `proto` files. 
Getting these `proto` files and compiling them to `go` files
is a bit of a pain, so they're kept in source control.
To regenerate them, it's probably just simplest to use the Docker build process.

### Error handling and exit codes

A number of errors can prevent the application from running 100% correctly,
with the main source of errors being network failures when hitting the source API.
At start-up, the application downloads some basic data from the API and
tries to write an empty file to disk. 
If this start-up fails, the application will exit with one of the following exit codes:

- `101` - the environment variable `PATH_GTFS_RT_PERIODICITY_MILLISECS` is not an integer.
- `102` - the environment variable `PATH_GTFS_RT_SOURCE_API` is set, but is not equal to either `grpc` or `http`.
- `103` - unable to write to the specified file output path. This can often indicate a filesystem permissions error. 
- `104` - there was an error retrieving [routes data](https://path.api.razza.dev/v1/routes) from the API (for example, network failure).
- `105` - there was an error retrieving [stations data](https://path.api.razza.dev/v1/stations) from the API.

After start-up, any errors encountered are handled gracefully, and the application will not exit until interrupted.
If, during a particular update, the realtime data for a specific stop cannot be retrieved, or is malformed,
then the previously retrieved data will be used.

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
