![example workflow](https://github.com/chyngyz-sydykov/go-web/actions/workflows/ci.yml/badge.svg)
![Go Coverage](https://github.com/chyngyz-sydykov/go-web/wiki/coverage.svg)


# Installation

 - clone the repo
 - install docker
 - copy `.env.dist` to `.env`
 - run `docker-compose up --build`
 - if everything is ok, please check `http://localhost:8000/api/v1/books` url in the browser

# Testing

On initial project setup, please manually create a database
`APP_ENV=test go test ./tests/`
run test without cache `go test -count=1 ./tests/`
running test within docker `docker exec -it go_rest_api bash -c "APP_ENV=test go test -count=1 ./tests"`
running the test coverage on local machine `docker exec -it go_rest_api bash "scripts/coverage.sh"`
                `go tool cover -html=coverage/filtered_coverage.out`

# GRPC

the protobuf files are stored in different repo https://github.com/chyngyz-sydykov/book-rating-protos and it is imported via following command.

generate grpc files `docker exec -it go_rest_api bash -c "./generate_protoc.sh"`

in order to communicate with the rate microservice on local machine, do following

1. create a local network `docker network inspect grpc-network`
2. after running `docker-compose up` run `docker network inspect grpc-network`
You should see a json with the list of containers ex:
```
"Containers": {
            "some_hash": {
                "Name": "go_rating_postgres_db","
            },
            "some_hash": {
                "Name": "go_rest_api",
            },
            "some_hash": {
                "Name": "go_postgres_db",
            },
            "some_hash": {
                "Name": "go_rating_server",
            }
},
```

# Handy commands

To install new package

`go get package_name`

to clean up go.sum run

`go mod tidy`

to run test

running project via docker
`docker-compose up --build`
`docker-compose down`

`docker-compose logs -f`