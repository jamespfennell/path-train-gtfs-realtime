# GTFS Realtime for the PATH Train

This repository hosts a simple Go application
that reads PATH train realtime data from Matt Razza's public API
and outputs the data in the GTFS Realtime format.
Some important notes:

- You don't need to run the application yourself.
    The GTFS Realtime feed produced by this software can be accessed at `https://transitdata.nyc/path.gtfsrt`.

- The outputted data is compatible with the official GTFS Static data published by the Port Authority
    in the sense that the stop IDs and route IDs match up.
    The feed should work correctly for software that integrates realtime and static data.

- Unfortunately the Port Authority doesn't distribute full realtime data set, and so the GTFS 
  Realtime feed has some missing pieces:
  - There trip data: all the Port Authority shares are stops, and arrival times at those stops.
    There is no easy way to connect arrival times for the same trip at multiple stops.
    So, in the GTFS Realtime feed, the "trips" are dummy trips with a random ID and a single 
    stop time update. This should be sufficient for consumers that want to show arrival times at stops,
    but of course blocks other uses like tracking trains through the system.
  - The GTFS Static feed describes all the tracks/platforms at each of the PATH stations
    but in the realtime data we don't known which platform a train will stop at.
    In the realtime feed, all of the trains stop at the "station" stop (i.e., the stop in the static
    feed with location type `1`).


## Exit codes

- 1, could not download routes data
- 2, result was errored
- 3, could not download stations data
- 4, result was errored