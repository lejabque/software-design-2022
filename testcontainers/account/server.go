package account

import (
	context "context"

	"github.com/lejabque/software-design-2022/testcontainers/exchange"
)

type AccountsStorage interface {
	UpdateAccount(account *Account)
	GetAccount(id string) *Account
	Deposit(accountId string, amount float64)
	BuyStocks(accountId string, stockName string, totalPrice float64, amount int64) error
	SellStocks(accountId string, stockName string, price float64, amount int64) error
}

type accountServer struct {
	UnimplementedAccountServiceServer
	accounts       AccountsStorage
	exchangeClient exchange.StockExchangeClient
}

func NewAccountServer(accounts AccountsStorage, exchangeClient exchange.StockExchangeClient) AccountServiceServer {
	return &accountServer{
		accounts:       accounts,
		exchangeClient: exchangeClient,
	}
}

func (s *accountServer) CreateAccount(ctx context.Context, req *CreateAccountRequest) (*CreateAccountResponse, error) {
	account := &Account{
		Id:     req.Id,
		Stocks: make(map[string]int64),
	}
	s.accounts.UpdateAccount(account)
	return &CreateAccountResponse{}, nil
}

func (s *accountServer) Deposit(ctx context.Context, req *DepositRequest) (*DepositResponse, error) {
	s.accounts.Deposit(req.Id, req.Amount)
	return &DepositResponse{}, nil
}

func (s *accountServer) BuyStocks(ctx context.Context, req *BuyRequest) (*BuyResponse, error) {
	// yes, that's not transactional, but it's just an example
	resp, err := s.exchangeClient.BuyStock(ctx, &exchange.BuyRequest{
		Name:   req.StockName,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, err
	}
	err = s.accounts.BuyStocks(req.Id, req.StockName, resp.Price, req.Amount)
	if err != nil {
		return nil, err
	}
	return &BuyResponse{Price: resp.Price}, nil
}

func (s *accountServer) SellStocks(ctx context.Context, req *SellRequest) (*SellResponse, error) {
	// yes, that's not transactional, but it's just an example
	resp, err := s.exchangeClient.SellStock(ctx, &exchange.SellRequest{
		Name:   req.StockName,
		Amount: req.Amount,
	})
	if err != nil {
		return nil, err
	}
	err = s.accounts.SellStocks(req.Id, req.StockName, resp.Price, req.Amount)
	if err != nil {
		return nil, err
	}
	return &SellResponse{Price: resp.Price}, nil
}

func (s *accountServer) GetAccount(ctx context.Context, req *GetAccountRequest) (*GetAccountResponse, error) {
	account := s.accounts.GetAccount(req.Id)
	return &GetAccountResponse{
		Account: account,
	}, nil
}

func (s *accountServer) TotalBalance(ctx context.Context, req *TotalBalanceRequest) (*TotalBalanceResponse, error) {
	account := s.accounts.GetAccount(req.Id)
	totalBalance := account.Balance
	for name, amount := range account.Stocks {
		resp, err := s.exchangeClient.GetStockPrice(ctx, &exchange.StockPriceRequest{
			Name: name,
		})
		if err != nil {
			return nil, err
		}
		totalBalance += float64(amount) * resp.Price
	}
	return &TotalBalanceResponse{
		Balance: totalBalance,
	}, nil
}
