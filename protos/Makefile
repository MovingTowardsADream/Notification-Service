protoc_version:
	@echo "Checking protoc version"
	protoc --version

generate_proto:
	@echo "Generating proto files"
	protoc -I proto proto/notify/notify.proto --go_out=./gen/go --go_opt=paths=source_relative --go-grpc_out=./gen/go/ --go-grpc_opt=paths=source_relative