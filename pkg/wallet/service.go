package wallet

import (
	"Lesson17_ful_export/pkg/types"
	"bufio"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/google/uuid"
)

var ErrPhoneRegistered = errors.New("phone already registered")
var ErrAccountNotFound = errors.New("account not found")
var ErrAmountMostBePositive = errors.New("amount must be greater than zero")
var ErrNotEnouthBalance = errors.New("not enough balance in account")
var ErrPaymentNotFound = errors.New("payment not found")
var ErrFavoriteNotFound = errors.New("favorite not found")

type Service struct {
	nextAccountId int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}

func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, acc := range s.accounts {
		if acc.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountId++
	account := &types.Account{
		ID:      s.nextAccountId,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

func (s *Service) FindAccountById(accountId int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountId {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}

func (s *Service) Deposit(accountId int64, amount types.Money) (*types.Account, error) {
	if amount <= 0 {
		return nil, ErrAmountMostBePositive
	}
	account, err := s.FindAccountById(accountId)
	if err != nil {
		return nil, err
	}
	account.Balance += amount
	return account, nil
}

func (s *Service) Pay(acountId int64, category types.PaymentCategory, amount types.Money) (*types.Payment, error) {
	if amount < 0 {
		return nil, ErrAmountMostBePositive
	}
	account, err := s.FindAccountById(acountId)
	if err != nil {
		return nil, err
	}
	if account.Balance < amount {
		return nil, ErrNotEnouthBalance
	}
	account.Balance -= amount

	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		AccountID: acountId,
		Amount:    amount,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil
}

func (s *Service) FindPaymentById(paymentId string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentId {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	payment, err := s.FindPaymentById(paymentID)
	if err != nil {
		return nil, ErrPaymentNotFound
	}
	favoriteID := uuid.New().String()
	favorite := &types.Favorite{
		ID:        favoriteID,
		AccountId: payment.AccountID,
		Name:      name,
		Amount:    payment.Amount,
		Category:  payment.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil
}
func (s *Service) FindFavoriteByID(favoriteID string) (*types.Favorite, error) {
	for _, fav := range s.favorites {
		if fav.ID == favoriteID {
			return fav, nil
		}
	}
	return nil, ErrFavoriteNotFound
}
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	favorite, err := s.FindFavoriteByID(favoriteID)
	if err != nil {
		return nil, err
	}
	return s.Pay(favorite.AccountId, favorite.Category, favorite.Amount)
}

//Возврашает true если в Service есть данный
func (s *Service) IsServiceEmpty() (bool, bool, bool) {
	return len(s.accounts) > 0, len(s.payments) > 0, len(s.favorites) > 0
}

//UpdateService обновляет сервис аккаунт с указыном Id если нет то создаёт новый аккаунт с указыном данный
func (s *Service) UpdateServiceAccount(account *types.Account) {
	acc, err := s.FindAccountById(account.ID)
	if err != nil {
		s.nextAccountId++
		accountNew := &types.Account{
			ID:      s.nextAccountId,
			Phone:   account.Phone,
			Balance: account.Balance,
		}
		s.accounts = append(s.accounts, accountNew)
	}
	acc.Balance = account.Balance
	acc.Phone = account.Phone
}
func (s *Service) UpdateServicePayment(payment *types.Payment) {
	paym, err := s.FindPaymentById(payment.ID)
	if err != nil {
		s.payments = append(s.payments, payment)
	}
	paym.Amount = payment.Amount
	paym.Category = payment.Category
	paym.Status = payment.Status
}
func (s *Service) UpdateServiceFavorite(favorite *types.Favorite) {
	fav, err := s.FindPaymentById(favorite.ID)
	if err != nil {
		s.favorites = append(s.favorites, favorite)
	}
	fav.Amount = favorite.Amount
	fav.Category = favorite.Category
}
func (s *Service) Export(dir string) error {
	a, p, f := s.IsServiceEmpty()
	if a {
		var result string
		for _, acc := range s.accounts {
			result += strconv.Itoa(int(acc.ID)) + " " + string(acc.Phone) + " " + strconv.Itoa(int(acc.Balance)) + ";\n"
		}
		dir += "/accounts.dump"
		err := ioutil.WriteFile(dir, []byte(result), 0666)
		if err != nil {
			return err
		}
	}
	if p {
		var result string
		for _, payment := range s.payments {
			result += string(payment.ID) + " " + strconv.Itoa(int(payment.AccountID)) + " " + strconv.Itoa(int(payment.Amount))
			result += " " + string(payment.Category) + " " + string(payment.Status) + ";\n"
		}
		dir += "/payments.dump"
		err := ioutil.WriteFile(dir, []byte(result), 0666)
		if err != nil {
			return err
		}
	}
	if f {
		var result string
		for _, favorite := range s.favorites {
			result += string(favorite.ID) + " " + strconv.Itoa(int(favorite.AccountId)) + " " + favorite.Name + " "
			result += strconv.Itoa(int(favorite.Amount)) + string(favorite.Category) + ";\n"
		}
		dir += "/favorites.dump"
		err := ioutil.WriteFile(dir, []byte(result), 0666)
		if err != nil {
			return err
		}
	}
	return nil
}
func (s *Service) Import(dir string) error {
	a, p, f := s.IsServiceEmpty()
	if a {
		dir += "/accounts.dump"
		err := s.ImportFromFile(dir)
		if err != nil {
			return err
		}
	}
	if p {
		dir += "/payments.dump"
		err := s.ImportFromFile(dir)
		if err != nil {
			return err
		}
	}
	if f {
		dir += "/favorites.dump"
		err := s.ImportFromFile(dir)
		if err != nil {
			return err
		}
	}
	return nil

}
func (s *Service) ImportFromFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			if err == nil {
				err = cerr
			}
		}
	}()
	reader := bufio.NewReader(src)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if err == io.EOF {
			log.Print(line)
			break
		}
		log.Print(line)
	}
	return nil
}
