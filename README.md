# Authorization Service

GRPC service for user authorization. Easily integrated into microservice architecture.

## Run

Install protobuf: `brew install protobuf`

Generate protobuf: `make generate-proto`

Create config: `make env`

Run app with environment: `docker-compose up --build`

Run environment only: `docker-compose up postgres redis`

## Migrations

Starts with main application.

Package: https://github.com/golang-migrate/migrate

Create new migration: `make new-migration`

## Testing

Package: https://github.com/golang/mock

Run tests: `make test`

Integration tests [here](./test).

## TODO

- Unit-tests