# Project staffinc

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.


### API Docs

For API docs, we are using [Swagger](https://swagger.io/) with [Swag](https://github.com/swaggo/swag) Generator

- Modify api documentation based on [Swag](https://github.com/swaggo/swag#contents). **Required:** Api documentation
  must be placed in
  this [folder](https://bitbucket.org/hypefast-tech/marketing-promotion-service/src/master/docs/swagger/)

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
