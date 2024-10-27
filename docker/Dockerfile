# Start from the official Go image
FROM golang:1.23-alpine
# Set the Current Working Directory inside the container
WORKDIR /app


RUN apk add --no-cache postgresql-client
# COPY . .
# # Build the Go app
# RUN go build -o /app/bin/restApi /app/cmd/restApi/main.go
# Expose port 8080 to the outside world
EXPOSE 8000
# Run the executable
CMD ["go", "run", "cmd/restApi/main.go"]