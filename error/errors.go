package my_error

import "errors"

var ErrNotFound = errors.New("resource not found")

var ErrgRpcServerDown = errors.New("gRPC server is down")
