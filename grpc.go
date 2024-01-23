package pathgtfsrt

import (
	"context"
	"time"

	sourceapi "github.com/jamespfennell/path-train-gtfs-realtime/proto/sourceapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcApiUrl = "path.grpc.razza.dev:443"
)

// GrpcSourceClient is a source client that gets data using the Razza gRPC API.
type GrpcSourceClient struct {
	conn          *grpc.ClientConn
	stations      *sourceapi.StationsClient
	routes        *sourceapi.RoutesClient
	timeoutPeriod time.Duration
}

func NewGrpcSourceClient(timeoutPeriod time.Duration) (*GrpcSourceClient, error) {
	conn, err := grpc.Dial(grpcApiUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	stationsClient := sourceapi.NewStationsClient(conn)
	routesClient := sourceapi.NewRoutesClient(conn)
	return &GrpcSourceClient{conn: conn, stations: &stationsClient, routes: &routesClient, timeoutPeriod: timeoutPeriod}, nil
}

func (client *GrpcSourceClient) GetStationToStopId(ctx context.Context) (stationToStopId map[sourceapi.Station]string, err error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeoutPeriod)
	defer cancel()
	response, err := (*client.stations).ListStations(ctx, &sourceapi.ListStationsRequest{})
	if err != nil {
		return
	}
	stationToStopId = map[sourceapi.Station]string{}
	for _, stationData := range response.Stations {
		stationToStopId[stationData.Station] = stationData.Id
	}
	return
}

func (client *GrpcSourceClient) GetRouteToRouteId(ctx context.Context) (routeToRouteId map[sourceapi.Route]string, err error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeoutPeriod)
	defer cancel()
	response, err := (*client.routes).ListRoutes(ctx, &sourceapi.ListRoutesRequest{})
	if err != nil {
		return
	}
	routeToRouteId = map[sourceapi.Route]string{}
	for _, routeData := range response.Routes {
		routeToRouteId[routeData.Route] = routeData.Id
	}
	return
}

func (client *GrpcSourceClient) Close() error {
	return client.conn.Close()
}

func (client *GrpcSourceClient) GetTrainsAtStation(ctx context.Context, station sourceapi.Station) ([]Train, error) {
	ctx, cancel := context.WithTimeout(ctx, client.timeoutPeriod)
	defer cancel()
	request := sourceapi.GetUpcomingTrainsRequest{Station: station}
	response, err := (*client.stations).GetUpcomingTrains(ctx, &request)
	if err != nil {
		return nil, err
	}
	var trains []Train
	for _, train := range response.UpcomingTrains {
		applyRouteQaToTrain(train)
		trains = append(trains, train)
	}
	return trains, nil
}
