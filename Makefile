test:
	@go test -v -cover ./...

proto.build:
	@protoc -I gamer/ gamer/gamer.proto --go_out=plugins=grpc:gamer
	@python -m grpc_tools.protoc -Igamer --python_out=gamer/python3 --grpc_python_out=gamer/python3 gamer/gamer.proto
