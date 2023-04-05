package account

import (
	context "context"

	"github.com/lejabque/software-design-2022/testcontainers/internal/api/accountapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/api/exchangeapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/lib"
)

type AccountsStorage interface {
	UpdateAccount(account *accountapi.Account)
	GetAccount(id string) *accountapi.Account
	Deposit(accountId string, amount float64)
	BuyStocks(accountId string, stockName string, totalPrice float64, amount int64) error
	SellStocks(accountId string, stockName string, price float64, amount int64) error
}

type accountServer struct {
	accountapi.UnimplementedAccountServiceServer
	accounts       AccountsStorage
	exchangeClient exchangeapi.StockExchangeClient
}

func NewAccountServer(accounts AccountsStorage, exchangeClient exchangeapi.StockExchangeClient) accountapi.AccountServiceServer {
	return &accountServer{
		accounts:       accounts,
		exchangeClient: exchangeClient,
	}
}

func (s *accountServer) CreateAccount(ctx context.Context, req *accountapi.CreateAccountRequest) (*accountapi.CreateAccountResponse, error) {
	account := &accountapi.Account{
		Id:     req.Id,
		Stocks: make(map[string]int64),
	}
	s.accounts.UpdateAccount(account)
	return &accountapi.CreateAccountResponse{}, nil
}

func (s *accountServer) Deposit(ctx context.Context, req *accountapi.DepositRequest) (*accountapi.DepositResponse, error) {
	s.accounts.Deposit(req.Id, req.Amount)
	return &accountapi.DepositResponse{}, nil
}

func (s *accountServer) BuyStocks(ctx context.Context, req *accountapi.BuyRequest) (*accountapi.BuyResponse, error) {
	// yes, that's not transactional, but it's just an example
	resp, err := s.exchangeClient.BuyStock(ctx, &exchangeapi.BuyRequest{
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
	return &accountapi.BuyResponse{Price: resp.Price}, nil
}

func (s *accountServer) SellStocks(ctx context.Context, req *accountapi.SellRequest) (*accountapi.SellResponse, error) {
	// yes, that's not transactional, but it's just an example
	resp, err := s.exchangeClient.SellStock(ctx, &exchangeapi.SellRequest{
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
	return &accountapi.SellResponse{Price: resp.Price}, nil
}

func (s *accountServer) GetAccount(ctx context.Context, req *accountapi.GetAccountRequest) (*accountapi.GetAccountResponse, error) {
	account := s.accounts.GetAccount(req.Id)
	if account == nil {
		return nil, lib.ErrAccountNotFound
	}
	return &accountapi.GetAccountResponse{
		Account: account,
	}, nil
}

func (s *accountServer) TotalBalance(ctx context.Context, req *accountapi.TotalBalanceRequest) (*accountapi.TotalBalanceResponse, error) {
	account := s.accounts.GetAccount(req.Id)
	totalBalance := account.Balance
	for name, amount := range account.Stocks {
		resp, err := s.exchangeClient.GetStockPrice(ctx, &exchangeapi.StockPriceRequest{
			Name: name,
		})
		if err != nil {
			return nil, err
		}
		totalBalance += float64(amount) * resp.Price
	}
	return &accountapi.TotalBalanceResponse{
		Balance: totalBalance,
	}, nil
}
