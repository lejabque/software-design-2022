package account

import (
	context "context"
	"testing"
	"time"

	"github.com/lejabque/software-design-2022/testcontainers/internal/api/accountapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/api/exchangeapi"
	"github.com/lejabque/software-design-2022/testcontainers/internal/repos"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AccountServerSuite struct {
	suite.Suite
	accounts AccountsStorage
	server   accountapi.AccountServiceServer

	exchange     testcontainers.Container
	exchangeConn *grpc.ClientConn
}

func (s *AccountServerSuite) equalFloat64(a, b float64) {
	s.LessOrEqual(a-b, 0.0001)
}

func (s *AccountServerSuite) SetupTest() {
	s.accounts = repos.NewInMemoryAccountsStorage()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "exchange:latest",
		ExposedPorts: []string{"30080/tcp"},
		WaitingFor:   wait.ForLog("starting server"),
	}
	exchangeContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
		Logger:           testcontainers.Logger,
	})
	s.Require().NoError(err)
	s.exchange = exchangeContainer
	time.Sleep(1 * time.Second)

	endpoint, err := exchangeContainer.Endpoint(ctx, "")
	s.Require().NoError(err)
	exchangeConn, err := grpc.Dial(endpoint, grpc.WithInsecure())
	s.Require().NoError(err)
	s.exchangeConn = exchangeConn

	s.server = NewAccountServer(s.accounts, exchangeapi.NewStockExchangeClient(exchangeConn))
}

func (s *AccountServerSuite) TearDownTest() {
	ctx := context.Background()
	s.exchange.Terminate(ctx)
	s.exchangeConn.Close()
}

func (s *AccountServerSuite) setStockData(name string, price float64, amount int64) {
	client := exchangeapi.NewStockExchangeClient(s.exchangeConn)
	ctx := context.Background()
	_, err := client.SetStockData(ctx, &exchangeapi.SetStockDataRequest{
		Name:   name,
		Price:  price,
		Amount: amount,
	})
	s.Require().NoError(err)
}

func (s *AccountServerSuite) TestBasicAccount() {
	ctx := context.Background()
	_, err := s.server.CreateAccount(ctx, &accountapi.CreateAccountRequest{
		Id: "existing",
	})
	s.Require().NoError(err)

	resp, err := s.server.GetAccount(ctx, &accountapi.GetAccountRequest{
		Id: "existing",
	})
	if s.NoError(err) {
		acc := resp.GetAccount()
		s.Equal("existing", acc.Id)
		s.Equal(float64(0), acc.Balance)
		s.Len(acc.Stocks, 0)
	}

	_, err = s.server.GetAccount(ctx, &accountapi.GetAccountRequest{
		Id: "unknown",
	})
	s.Error(err)
	s.Equal(codes.NotFound, status.Code(err))

	_, err = s.server.Deposit(ctx, &accountapi.DepositRequest{
		Id:     "existing",
		Amount: 100,
	})
	s.Require().NoError(err)

	resp, err = s.server.GetAccount(ctx, &accountapi.GetAccountRequest{
		Id: "existing",
	})
	if s.NoError(err) {
		acc := resp.GetAccount()
		s.equalFloat64(float64(100), acc.Balance)
	}
}

func (s *AccountServerSuite) TestBuyAndSellStocks() {
	accountID := "test"
	ctx := context.Background()

	_, err := s.server.CreateAccount(ctx, &accountapi.CreateAccountRequest{
		Id: accountID,
	})
	s.Require().NoError(err)

	_, err = s.server.Deposit(ctx, &accountapi.DepositRequest{
		Id:     accountID,
		Amount: 1000,
	})
	s.Require().NoError(err)

	s.setStockData("YNDX", 82.84, 10)
	_, err = s.server.BuyStocks(ctx, &accountapi.BuyRequest{
		Id:        accountID,
		Amount:    5,
		StockName: "YNDX",
	})
	s.Require().NoError(err)

	resp, err := s.server.GetAccount(ctx, &accountapi.GetAccountRequest{
		Id: accountID,
	})
	balance := float64(1000 - 82.84*5)
	if s.NoError(err) {
		acc := resp.GetAccount()
		s.Len(acc.Stocks, 1)
		s.Equal(int64(5), acc.Stocks["YNDX"])
		s.equalFloat64(balance, acc.Balance)
	}

	totalResp, err := s.server.TotalBalance(ctx, &accountapi.TotalBalanceRequest{
		Id: accountID,
	})
	if s.NoError(err) {
		s.equalFloat64(float64(1000), totalResp.GetBalance())
	}

	// something happened
	s.setStockData("YNDX", 18.94, 5)
	totalResp, err = s.server.TotalBalance(ctx, &accountapi.TotalBalanceRequest{
		Id: accountID,
	})
	if s.NoError(err) {
		s.equalFloat64(balance+18.94*5, totalResp.GetBalance())
	}

	_, err = s.server.SellStocks(ctx, &accountapi.SellRequest{
		Id:        accountID,
		Amount:    3,
		StockName: "YNDX",
	})
	s.Require().NoError(err)

	resp, err = s.server.GetAccount(ctx, &accountapi.GetAccountRequest{
		Id: accountID,
	})
	balance = balance + 18.94*3
	if s.NoError(err) {
		acc := resp.GetAccount()
		s.Len(acc.Stocks, 1)
		s.Equal(int64(2), acc.Stocks["YNDX"])
		s.equalFloat64(balance, acc.Balance)
	}
}

func TestAccountServerSuite(t *testing.T) {
	suite.Run(t, new(AccountServerSuite))
}
