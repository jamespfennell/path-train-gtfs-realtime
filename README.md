# GTFS Realtime for the PATH Train

This repository hosts a simple Go application
that reads PATH train realtime data from [Matt Razza's public HTTP API](https://github.com/mrazza/path-data)
and outputs the data in the GTFS Realtime format.
Some important notes:

- You don't need to run the application yourself.
    The GTFS Realtime feed produced by this software can be accessed at 
    [`https://path.transitdata.nyc/gtfsrt`](https://path.transitdata.nyc/gtfsrt).

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
 
There are two configurations, both read as environment variables:

- `PATH_GTFS_RT_PERIODICITY_MILLISECS` - the periodicity, in milliseconds,
    with which to update the GTFS Realtime file. The value of this environment
    variable must be castable to an integer, otherwise the application will fail to start.
    The default is 5000 milliseconds (so 5 seconds).
    
- `PATH_GTFS_RT_OUTPUT_PATH` - the location on disk to output the file,
    either absolute or relative to the working directory when the application starts.
    The default is `./path.gtfsrt`.


### Running using Docker

As part of the Travis CI job for this repo, a Docker image is built and stored
in the `jamespfennell/path-train-gtfs-realtime` repository on Docker Hub.
There are both `latest` tags and `build<n>` tags, where `<n>` is the Travis build number.
You can also build the Docker image locally by running `docker build .` in the
root of the repo.

It is generally simplest to run the application using Docker.
In the container, the GTFS Realtime file is written to `/output/gtfsrt`.
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

The application has a few dependencies which should be installed first:
```
go get github.com/google/uuid
go get github.com/golang/protobuf/proto
go get github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs
```
After that, just run `go run path-train-gtfs-realtime.go`.


### Error handling and exit codes

A number of errors can prevent the application from running 100% correctly,
with the main source of errors being network failures when hitting the Razza API.
At startup, the application downloads some basic data from the API and
tries to write an empty file to disk. 
If this start-up fails, the application will exit with one of the following exit codes:

- `101` - the environment variable `PATH_GTFS_RT_PERIODICITY_MILLISECS` is not an integer.
- `102` - there was an error retrieving [routes data](https://path.api.razza.dev/v1/routes) from the API (for example, network failure).
- `103` - a routes data response was received from the API, but the response was malformed.
- `104` - there was an error retrieving [stations data](https://path.api.razza.dev/v1/stations) from the API.
- `105` - a stations data response was received from the API, but the response was malformed.
- `106` - unable to write to the specified file output path. This can often indicate a filesystem permissions error. 

After start-up, any errors encountered are handled gracefully, and the application will not exit until interrupted.
If, during a particular update, the realtime data for a specific stop cannot be retrieved, or is malformed,
then the previously retrieved data will be used.
