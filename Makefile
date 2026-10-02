.PHONY: gen-proto serve build

build:
	go build -o build/apid ./cmd/apid

gen-proto:
	protoc --proto_path=proto proto/*.proto --go_out=. --go-grpc_out=.

serve:
	go run ./cmd/apid/main.go

## preserve for grpc client
serve-client:
	go run ./cmd/client/main.go
