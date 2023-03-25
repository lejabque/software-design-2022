package exchange

import "errors"

// todo: grpc
var ErrStockNotFound = errors.New("stock not found")

var ErrNotEnoughStocks = errors.New("not enough stocks")
