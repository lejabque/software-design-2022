package exchange

import (
	"fmt"
	"net"

	"github.com/lejabque/software-design-2022/testcontainers/internal/app"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(args *app.CliArgs) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", args.Port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	logger := zap.NewExample()
	defer logger.Sync()

	RegisterStockExchangeServer(s, &exchangeServer{})
	reflection.Register(s)
	logger.Info("starting server", zap.Uint16("port", args.Port))
	if err := s.Serve(listener); err != nil {
		logger.Fatal("failed to serve", zap.Error(err))
	}
}
