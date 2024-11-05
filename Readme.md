To install new package

`go get package_name`

to clean up go.sum run

`go mod tidy`

to run test
#Testing

On initial project setup, please manually create a database
`APP_ENV=test go test ./tests/`
run test without cache `go test -count=1 ./tests/`
running test within docker `docker exec -it go_rest_api bash -c "APP_ENV=test go test -count=1 ./tests"`

running project via docker
`docker-compose up --build`
`docker-compose down`

`docker-compose logs -f`