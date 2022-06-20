package types

//Money представляет собой денежную сумму в минимпльных единицах (центы, копейки, дирами и.т)
type Money int64

//PaymentCategory представляет собой информацию о категориях который был совершён платёж (авто, аптеки, рестороны и.т)
type PaymentCategory string

//PaymentStatus представляет собой  статус платёжа.
type PaymentStatus string

// Предопределёный статус платёжей
const (
	PaymentStatusOk         PaymentStatus = "Ok"
	PaymentStatusFail       PaymentStatus = "Fail"
	PaymentStatusInProgress PaymentStatus = "InProgress"
)

//Payment представляет информайию о платёже.
type Payment struct {
	ID        string
	AccountID int64
	Amount    Money
	Category  PaymentCategory
	Status    PaymentStatus
}

type Phone string

//Account представляет информацию о счёте пользователья
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}

//Favorite представляет платёж в избранное.
type Favorite struct {
	ID        string
	AccountId int64
	Name      string
	Amount    Money
	Category  PaymentCategory
}
