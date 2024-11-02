# Start from the official Go image
FROM golang:1.23-alpine

# Install dependencies
RUN apk add --no-cache git bash postgresql-client
RUN go env -w GOFLAGS=-buildvcs=false
RUN go install github.com/air-verse/air@latest

# Set the Current Working Directory inside the container
WORKDIR /app
COPY . .
# RUN adduser -D -g '' golang
# RUN mkdir -p /app && chown -R golang:golang /app /go/pkg/mod
# WORKDIR /app
# USER golang
# COPY --chown=golang:golang . .

EXPOSE 8000

# Run the executable
CMD ["air", "-c", ".air.toml"]