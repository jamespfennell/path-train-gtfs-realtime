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

- Unfortunately the Port Authority


## Exit codes

- 1, could not download routes data
- 2, result was errored
- 3, could not download stations data
- 4, result was errored