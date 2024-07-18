gen:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative usermgmt/usermgmt.proto

server:
	go run usermgmt_server/usermgmt_server.go 
client:
	go run usermgmt_client/usermgmt_client.go 