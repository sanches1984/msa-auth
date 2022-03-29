generate:
	@cd api && rm -rf *
	@protoc *.proto  -I . -I .. --go_out=plugins=grpc:api
	@echo "Client generated"