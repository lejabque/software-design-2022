package repos

import (
	sync "sync"

	"github.com/lejabque/software-design-2022/testcontainers/internal/lib"
)

type StockData struct {
	Price  float64
	Amount int64
}

type inMemoryStocks struct {
	lock   *sync.Mutex
	stocks map[string]*StockData
}

func NewInMemoryStocksStorage() *inMemoryStocks {
	return &inMemoryStocks{
		lock:   &sync.Mutex{},
		stocks: make(map[string]*StockData),
	}
}

// GetStockData(name string) (*StockData, error)
// 	BuyStock(name string, amount int64) (float64, error)
// 	SellStock(name string, amount int64) (float64, error)
// 	SetStockData(name string, price float64, amount int64) error
// }

func (s *inMemoryStocks) GetStockData(name string) (*StockData, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.stocks[name]
	if !ok {
		return nil, lib.ErrStockNotFound
	}
	return data, nil
}

func (s *inMemoryStocks) BuyStock(name string, amount int64) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.stocks[name]
	if !ok {
		return 0, lib.ErrStockNotFound
	}
	if data.Amount < amount {
		return 0, lib.ErrNotEnoughStocks
	}
	data.Amount -= amount
	return data.Price, nil
}

func (s *inMemoryStocks) SellStock(name string, amount int64) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.stocks[name]
	if !ok {
		return 0, lib.ErrStockNotFound
	}
	data.Amount += amount
	return data.Price, nil
}

func (s *inMemoryStocks) SetStockData(name string, price float64, amount int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stocks[name] = &StockData{
		Price:  price,
		Amount: amount,
	}
	return nil
}
