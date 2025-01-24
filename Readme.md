![example workflow](https://github.com/chyngyz-sydykov/go-web/actions/workflows/ci.yml/badge.svg)
![Go Coverage](https://github.com/chyngyz-sydykov/go-web/wiki/coverage.svg)

# About the project

This is a one of the microservices for personal pet project as study practice. the whole system consists of 3 microservices:
 - **go-web** (current project) works as an api gateway. the endpoints include CRUD actions for book, create endpoint for saving rating
 - **go-rating** another microservice that saves new rating and return list of rating by book id. the communication between go-web and go-rating is via gRPC [link](https://github.com/chyngyz-sydykov/go-rating)
- **go-recommendation** On Progress third microservice that will hold business logic related with recommendation of books depending on rating and how often the book is edited or created. the communication will be done via RabbitMQ

# Installation

 - clone the repo
 - install docker
 - copy `.env.dist` to `.env`
 - run `docker-compose up --build`
 - if everything is ok, please check `http://localhost:8000/api/v1/books` url in the browser

# Testing

On initial project setup, please manually create a database for tests. check the database name in env.test file. to run use following commands:

run tests `APP_ENV=test go test ./tests/`

run tests without cache `go test -count=1 ./tests/`

run tests within docker (preferred way) `docker exec -it go_rest_api bash -c "APP_ENV=test go test -count=1 ./tests"`

run test coverage on local machine `docker exec -it go_rest_api bash "scripts/coverage.sh"`
`go tool cover -html=coverage/filtered_coverage.out`

# SWAGGER

To generate swagger documentatation `docker exec -it go_rest_api bash -c "swag init --generalInfo router.go --dir ./application/router,./application/handlers --parseDependency"`

# GRPC

the protobuf files are stored in different repo https://github.com/chyngyz-sydykov/book-rating-protos and it is imported via following command.

generate grpc files `docker exec -it go_rest_api bash -c "./generate_protoc.sh"`

in order to communicate with the rate microservice on local machine, do following

1. create a local network `docker network create grpc-network`
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