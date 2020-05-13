# Essentially a collection of small Bash scripts.

get-protos:
	rm -rf git
	mkdir -p git
	git clone --single-branch --branch master https://github.com/google/transit.git git/transit
	git clone --single-branch --branch master https://github.com/mrazza/path-data git/path-data
	git clone --single-branch --branch master https://github.com/googleapis/googleapis.git git/googleapis

	mkdir -p gtfsrt
	mv git/transit/gtfs-realtime/proto/*.proto gtfsrt

	mkdir -p sourceapi
	mv git/path-data/proto/*.proto sourceapi
	rm -rf sourceapi/google
	mv git/googleapis/google sourceapi

	rm -rf git

build-protos:
	protoc --go_out=./gtfsrt  --proto_path=./gtfsrt ./gtfsrt/*.proto
	protoc --go_out=plugins=grpc:./sourceapi  --proto_path=./sourceapi ./sourceapi/*.proto

docker-hub:
	docker login --username ${DOCKER_USERNAME} --password ${DOCKER_PASSWORD}
	docker tag jamespfennell/path-train-gtfs-realtime:latest jamespfennell/path-train-gtfs-realtime:build${TRAVIS_BUILD_NUMBER}
	docker push jamespfennell/path-train-gtfs-realtime:latest
	docker push jamespfennell/path-train-gtfs-realtime:build${TRAVIS_BUILD_NUMBER}
