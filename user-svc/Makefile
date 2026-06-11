server:
	go run main.go

gen-proto:
	protoc --proto_path=internal/api/grpc/ --go_out=. --go-grpc_out=. internal/api/grpc/order/*.proto

test:
	go test ./... -v