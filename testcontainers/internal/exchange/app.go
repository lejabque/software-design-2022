package exchange

import (
	"fmt"
	"net"

	"github.com/lejabque/software-design-2022/testcontainers/internal/api/exchangeapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/lib"
	"github.com/lejabque/software-design-2022/testcontainers/internal/repos"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(args *lib.CliArgs) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", args.Port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	logger := zap.NewExample()
	defer logger.Sync()

	exchangeapi.RegisterStockExchangeServer(s, NewExchangeServer(repos.NewInMemoryStocksStorage()))
	reflection.Register(s)
	logger.Info("starting server", zap.Uint16("port", args.Port))
	if err := s.Serve(listener); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
