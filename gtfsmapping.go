package pathgtfsrt

import (
	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
)

// Snaphot from: https://transitfeeds.com/p/port-authority-of-new-york-and-new-jersey/384/latest/stops
var sourceStationToGtfsStopId = map[sourceapi.Station]string{
	sourceapi.Station_FOURTEENTH_STREET:   "26722",
	sourceapi.Station_TWENTY_THIRD_STREET: "26723",
	sourceapi.Station_THIRTY_THIRD_STREET: "26724",
	sourceapi.Station_NINTH_STREET:        "26725",
	sourceapi.Station_CHRISTOPHER_STREET:  "26726",
	sourceapi.Station_EXCHANGE_PLACE:      "26727",
	sourceapi.Station_GROVE_STREET:        "26728",
	sourceapi.Station_HARRISON:            "26729",
	sourceapi.Station_HOBOKEN:             "26730",
	sourceapi.Station_JOURNAL_SQUARE:      "26731",
	sourceapi.Station_NEWPORT:             "26732",
	sourceapi.Station_NEWARK:              "26733",
	sourceapi.Station_WORLD_TRADE_CENTER:  "26734",
}

// Snaphot from: https://transitfeeds.com/p/port-authority-of-new-york-and-new-jersey/384/latest/routes
var sourceRouteToGtfsRouteId = map[sourceapi.Route]string{
	sourceapi.Route_HOB_33:     "859",
	sourceapi.Route_HOB_WTC:    "860",
	sourceapi.Route_JSQ_33:     "861",
	sourceapi.Route_NWK_WTC:    "862",
	sourceapi.Route_JSQ_33_HOB: "1024",
}
