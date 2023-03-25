package exchange

import sync "sync"

type StockData struct {
	Price float64
	Count int64
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
		return nil, ErrStockNotFound
	}
	return data, nil
}

func (s *inMemoryStocks) BuyStock(name string, amount int64) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.stocks[name]
	if !ok {
		return 0, ErrStockNotFound
	}
	if data.Count < amount {
		return 0, ErrNotEnoughStocks
	}
	data.Count -= amount
	return data.Price, nil
}

func (s *inMemoryStocks) SellStock(name string, amount int64) (float64, error) {
	s.lock.Lock()
	defer s.lock.Unlock()
	data, ok := s.stocks[name]
	if !ok {
		return 0, ErrStockNotFound
	}
	data.Count += amount
	return data.Price, nil
}

func (s *inMemoryStocks) SetStockData(name string, price float64, amount int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stocks[name].Price = price
	s.stocks[name].Count = amount
	return nil
}
