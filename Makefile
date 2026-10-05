.PHONY: gen-proto serve build db-up db-down protoset

db-up:
	docker compose -f docker-compose.db.yml up -d

db-down:
	docker compose -f docker-compose.db.yml down 

build:
	go build -o build/apid ./cmd/apid

proto-gen:
	cd proto && buf generate

protoset:
	cd proto && buf build -o ../api.protoset

serve:
	go run ./cmd/apid

## preserve for grpc client
serve-client:
	go run ./cmd/client/main.go
