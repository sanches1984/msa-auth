.PHONY: generate-proto
generate-proto:
	@cd proto && make generate

.PHONY: test
test:
	go test -cover --tags=ci ./...

.PHONY: env
env:
	@cp .env.example .env

.PHONY: mocks
mocks:
	mockgen -package=mocks -source internal/pkg/storage/interface.go -destination internal/pkg/storage/mocks/mock.go
	mockgen -package=mocks -source internal/app/service/interface.go -destination internal/app/service/mocks/mock.go