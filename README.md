# Project staffinc

One Paragraph of project description goes here

## Getting Started

Run migration script and do watch script to hot reload when development
```bash
make watch
```

## How to generate mock automatically?

For mocking purposes, you can use [Mockery](https://github.com/vektra/mockery) to generate mock class(es) for your
interfaces.

```sh
make mock
```

## Component
For Creating new component or pages read guide on [Templ](https://templ.guide)

## API Docs

For API docs, we are using [Swagger](https://swagger.io/) with [Swag](https://github.com/swaggo/swag) Generator

- Modify api documentation based on [Swag](https://github.com/swaggo/swag#contents). 
  
- Install Swag

  ```sh
  go install github.com/swaggo/swag
  ```

- Generate api documentation using swag

  ```sh
  make swagger-doc
  ```


## Migration
[create migration](https://github.com/golang-migrate/migrate) installed

create migration file
```bash
migrate create -ext sql -dir files/db/migrations -seq migration_name
```
export POSTGRESQL_URL
```bash
export POSTGRESQL_URL='postgres://username:password@localhost:5432/db_name?sslmode=disable'
```
run migration command
```bash
migrate -database ${POSTGRESQL_URL} -path files/db/migrations up
```

## MakeFile

run all make commands with clean tests
```bash
make all build
```

build the application
```bash
make build
```

run the application
```bash
make run
```

Create DB container
```bash
make docker-run
```

Shutdown DB container
```bash
make docker-down
```

live reload the application
```bash
make watch
```

run the test suite
```bash
make test
```

clean up binary from the last build
```bash
make clean
```
