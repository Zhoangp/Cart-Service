proto:
	protoc  --go_out=. --plugin=/Users/nguyenhoangphuc/go/bin/protoc-gen-go --go_opt=paths=source_relative --go-grpc_out=. --plugin=/Users/nguyenhoangphuc/go/bin/protoc-gen-go-grpc --go-grpc_opt=paths=source_relative $(Input)


server:
	go run main.go