# GTFS Realtime for the PATH Train

This repository hosts a simple Go application
that reads PATH train realtime data from Matt Razza's public API
and outputs the data in the GTFS Realtime format.
Some important notes:

- You don't need to run the application yourself.
    The GTFS Realtime feed produced by this software can be accessed at `https://path.transitdata.nyc/gtfsrt`.

- The outputted data is compatible with the official GTFS Static data published by the Port Authority
    in the sense that the stop IDs and route IDs match up.
    The feed should work correctly for software that integrates realtime and static data.

- Unfortunately the Port Authority doesn't distribute the full realtime data set, and so the GTFS
  Realtime feed has some big missing pieces:
  - There is no trip data: all the Port Authority communicates are stops, and arrival times at those stops.
    There is no easy way to connect arrival times for the same train at multiple stops.
    So, in the GTFS Realtime feed, the "trips" are dummy trips with a random ID and a single 
    stop time update. This should be sufficient for consumers that want to show arrival times at stops,
    but of course blocks other uses like tracking trains through the system.
  - The GTFS Static feed describes all the tracks/platforms at each of the PATH stations
    but in the realtime data we don't known which platform a train will stop at.
    In the realtime feed, all of the trains stop at the "station" stop (i.e., the stop in the static
    feed with location type `1`).

## Running the application

Env var `GTFS_REALTIME_OUTPUT_PATH`

## Exit codes

- 101, could not download routes data
- 102, result was errored
- 103, could not download stations data
- 104, result was errored