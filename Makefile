generate-proto:
	@cd proto && make generate

test:
	go test -cover ./...