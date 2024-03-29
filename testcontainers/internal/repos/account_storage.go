package repos

import (
	sync "sync"

	"github.com/lejabque/software-design-2022/testcontainers/internal/api/accountapi"
)

type inMemoryAccounts struct {
	lock     *sync.Mutex
	accounts map[string]*accountapi.Account
}

func NewInMemoryAccountsStorage() *inMemoryAccounts {
	return &inMemoryAccounts{
		lock:     &sync.Mutex{},
		accounts: make(map[string]*accountapi.Account),
	}
}

func (s *inMemoryAccounts) UpdateAccount(account *accountapi.Account) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.accounts[account.Id] = account
}

func (s *inMemoryAccounts) GetAccount(id string) *accountapi.Account {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.accounts[id]
}

func (s *inMemoryAccounts) Deposit(accountId string, amount float64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	account := s.accounts[accountId]
	account.Balance += amount
	s.accounts[accountId] = account
}

func (s *inMemoryAccounts) BuyStocks(accountId string, stockName string, price float64, amount int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	account := s.accounts[accountId]
	// TODO: we of course should check if we have enough money but it's not necessary for this example
	account.Balance -= float64(amount) * price
	account.Stocks[stockName] += amount
	s.accounts[accountId] = account
	return nil
}

func (s *inMemoryAccounts) SellStocks(accountId string, stockName string, price float64, amount int64) error {
	s.lock.Lock()
	defer s.lock.Unlock()
	account := s.accounts[accountId]
	// TODO: we of course should check if we have enough stocks but it's not necessary for this example
	account.Balance += float64(amount) * price
	account.Stocks[stockName] -= amount
	s.accounts[accountId] = account
	return nil
}
