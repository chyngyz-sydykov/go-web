#Installation

 - clone the repo
 - install docker
 - copy `.env.dist` to `.env`
 - run `docker-compose up --build`
 - if everything is ok, please check `http://localhost:8000/api/v1/books` url in the browser

#Testing

On initial project setup, please manually create a database
`APP_ENV=test go test ./tests/`
run test without cache `go test -count=1 ./tests/`
running test within docker `docker exec -it go_rest_api bash -c "APP_ENV=test go test -count=1 ./tests"`

#GRPC

the protobuf files are stored in different repo https://github.com/chyngyz-sydykov/book-rating-protos and it is imported via following command.

generate grpc files `docker exec -it go_rest_api bash -c "./generate_protoc.sh"`


#Handy commands

To install new package

`go get package_name`

to clean up go.sum run

`go mod tidy`

to run test

running project via docker
`docker-compose up --build`
`docker-compose down`

`docker-compose logs -f`