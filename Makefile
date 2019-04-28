api.start:
	@go run app/api/*.go

colosseum.start:
	@go run app/colosseum/*.go

db.migrate:
	@go run app/migrate/*.go up

test:
	@go test -v -cover ./...

docker.build:
	@make docker.api.build
	@make docker.colosseum.build
	@make docker.migrate.build

docker.push:
	@make docker.api.push
	@make docker.colosseum.push
	@make docker.migrate.push

docker.api.build:
	@docker build -f app/api/Dockerfile -t hellodhlyn/undersky-api .

docker.api.push:
	@docker push hellodhlyn/undersky-api

docker.colosseum.build:
	@docker build -f app/colosseum/Dockerfile -t hellodhlyn/undersky-colosseum .

docker.colosseum.push:
	@docker push hellodhlyn/undersky-colosseum

docker.migrate.build:
	@docker build -f app/migrate/Dockerfile -t hellodhlyn/undersky-migrate .

docker.migrate.push:
	@docker push hellodhlyn/undersky-migrate

proto.build:
	@protoc -I gamer/ gamer/gamer.proto --go_out=plugins=grpc:gamer
	@python -m grpc_tools.protoc -Igamer --python_out=gamer/python3 --grpc_python_out=gamer/python3 gamer/gamer.proto
