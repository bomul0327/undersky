api.start:
	@go run app/api/*.go

colosseum.start:
	@go run app/colosseum/*.go

db.migrate:
	@go run app/migrate/*.go up

test:
	@go test -v -cover ./...

docker.build:
	@make docker.colosseum.build

docker.colosseum.build:
	@docker build -f app/colosseum/Dockerfile -t hellodhlyn/undersky-colosseum .

docker.colosseum.run:
	@docker build --rm hellodhlyn/undersky-colosseum ./colosseum

proto.build:
	@protoc -I gamer/ gamer/gamer.proto --go_out=plugins=grpc:gamer
	@python -m grpc_tools.protoc -Igamer --python_out=gamer/python3 --grpc_python_out=gamer/python3 gamer/gamer.proto
