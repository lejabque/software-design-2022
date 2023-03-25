package lib

import (
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var ErrStockNotFound = status.Error(codes.NotFound, "stock not found")

var ErrNotEnoughStocks = status.Error(codes.FailedPrecondition, "not enough stocks")
