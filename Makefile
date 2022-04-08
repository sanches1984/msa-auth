generate-proto:
	@cd proto && make generate

test:
	go test -cover ./...

env:
	@cp .env.example .env