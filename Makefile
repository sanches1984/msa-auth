.PHONY: generate-proto
generate-proto:
	@cd proto && make generate

.PHONY: test
test:
	go test -cover --tags=ci ./...

.PHONY: env
env:
	@cp .env.example .env