package exchange

import (
	context "context"

	"github.com/lejabque/software-design-2022/testcontainers/internal/api/exchangeapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/repos"
)

type StocksStorage interface {
	GetStockData(name string) (*repos.StockData, error)
	BuyStock(name string, amount int64) (float64, error)
	SellStock(name string, amount int64) (float64, error)
	SetStockData(name string, price float64, amount int64) error
}

type exchangeServer struct {
	exchangeapi.UnimplementedStockExchangeServer
	stocks StocksStorage
}

func NewExchangeServer(stocks StocksStorage) exchangeapi.StockExchangeServer {
	return &exchangeServer{
		stocks: stocks,
	}
}

func (s *exchangeServer) GetStockPrice(ctx context.Context, req *exchangeapi.StockPriceRequest) (*exchangeapi.StockPriceResponse, error) {
	data, err := s.stocks.GetStockData(req.Name)
	if err != nil {
		return nil, err
	}
	return &exchangeapi.StockPriceResponse{
		Price: data.Price,
	}, nil
}

func (s *exchangeServer) BuyStock(ctx context.Context, req *exchangeapi.BuyRequest) (*exchangeapi.BuyResponse, error) {
	price, err := s.stocks.BuyStock(req.Name, req.Amount)
	if err != nil {
		return nil, err
	}
	return &exchangeapi.BuyResponse{
		Price: price,
	}, nil
}

func (s *exchangeServer) SellStock(ctx context.Context, req *exchangeapi.SellRequest) (*exchangeapi.SellResponse, error) {
	price, err := s.stocks.SellStock(req.Name, req.Amount)
	if err != nil {
		return nil, err
	}
	return &exchangeapi.SellResponse{
		Price: price,
	}, nil
}

func (s *exchangeServer) SetStockData(ctx context.Context, req *exchangeapi.SetStockDataRequest) (*exchangeapi.SetStockDataResponse, error) {
	err := s.stocks.SetStockData(req.Name, req.Price, req.Amount)
	if err != nil {
		return nil, err
	}
	return &exchangeapi.SetStockDataResponse{}, nil
}
