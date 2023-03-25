package exchange

import context "context"

type StocksStorage interface {
	GetStockData(name string) (*StockData, error)
	BuyStock(name string, amount int64) (float64, error)
	SellStock(name string, amount int64) (float64, error)
	SetStockData(name string, price float64, amount int64) error
}

type exchangeServer struct {
	UnimplementedStockExchangeServer
	stocks StocksStorage
}

func NewExchangeServer(stocks StocksStorage) StockExchangeServer {
	return &exchangeServer{
		stocks: stocks,
	}
}

func (s *exchangeServer) GetStockPrice(ctx context.Context, req *StockPriceRequest) (*StockPriceResponse, error) {
	data, err := s.stocks.GetStockData(req.Name)
	if err != nil {
		return nil, err
	}
	return &StockPriceResponse{
		Price: data.Price,
	}, nil
}

func (s *exchangeServer) BuyStock(ctx context.Context, req *BuyRequest) (*BuyResponse, error) {
	price, err := s.stocks.BuyStock(req.Name, req.Amount)
	if err != nil {
		return nil, err
	}
	return &BuyResponse{
		Price: price,
	}, nil
}

func (s *exchangeServer) SellStock(ctx context.Context, req *SellRequest) (*SellResponse, error) {
	price, err := s.stocks.SellStock(req.Name, req.Amount)
	if err != nil {
		return nil, err
	}
	return &SellResponse{
		Price: price,
	}, nil
}

func (s *exchangeServer) SetStockData(ctx context.Context, req *SetStockDataRequest) (*SetStockDataResponse, error) {
	err := s.stocks.SetStockData(req.Name, req.Price, req.Amount)
	if err != nil {
		return nil, err
	}
	return &SetStockDataResponse{}, nil
}
