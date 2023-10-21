package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

type Order struct {
	Source      string // source student/cb
	Destination string // dest student
	Currency    string
	Amount      decimal.Decimal
	Reference   string // reason
}

type Gecko struct {
	name string
	usd  float32
}

type CryptoDecimal struct {
	Basis        decimal.Decimal `json:"basis"`
	CurrentPrice decimal.Decimal `json:"currentPrice"`
	Name         string          `json:"name"`
	Quantity     decimal.Decimal `json:"quantity"`
}

type Transaction struct {
	Ts             time.Time
	Source         string
	Destination    string
	CurrencySource string
	CurrencyDest   string
	AmountSource   decimal.Decimal
	AmountDest     decimal.Decimal
	XRate          decimal.Decimal
	Reference      string
	FromSource     bool
	Net            decimal.Decimal `json:"-"`
	Balance        decimal.Decimal `json:"-"`
	Ubuck          decimal.Decimal `json:"-"`
}

type CertificateOfDeposit struct {
	Ts           time.Time
	Principal    int32
	CurrentValue decimal.Decimal
	RefundValue  decimal.Decimal
	Interest     float32
	Maturity     time.Time
	Active       bool `json:"active,omitempty"`
}

func buyMarketItem(db *bolt.DB, clock Clock, student UserInfo, teacher UserInfo, itemId string) (purchaseId string, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		teacherBucket, err := getTeacherBucketTx(tx, teacher.SchoolId, teacher.Email)
		if err != nil {
			return err
		}

		market := teacherBucket.Bucket([]byte(KeyMarket))
		if market == nil {
			return fmt.Errorf("failed to find market for: %v", teacher.LastName)
		}

		dataBucket := market.Bucket([]byte(itemId))
		if dataBucket == nil {
			return fmt.Errorf("ERROR cannot get market data bucket")
		}
		var details MarketItem
		marketData := dataBucket.Get([]byte(KeyMarketData))
		err = json.Unmarshal(marketData, &details)
		if err != nil {
			return fmt.Errorf("ERROR cannot unmarshal market details")
		}

		if details.Count == 0 {
			return fmt.Errorf("ERROR there is nothing left to buy")
		}

		err = chargeStudentTx(tx, clock, student, decimal.NewFromInt32(details.Cost), teacher.Email, "Purchased "+details.Title, true)
		if err != nil {
			return err
		}

		details.Count = details.Count - 1
		marshal, err := json.Marshal(details)
		if err != nil {
			return err
		}

		err = dataBucket.Put([]byte(KeyMarketData), marshal)
		if err != nil {
			return err
		}

		buyersBucket, err := dataBucket.CreateBucketIfNotExists([]byte(KeyBuyers))
		if err != nil {
			return err
		}

		marshal, err = json.Marshal(Buyer{Active: true, Id: student.Email})
		if err != nil {
			return err
		}

		purchaseId = RandomString(7)

		err = buyersBucket.Put([]byte(purchaseId), marshal)
		if err != nil {
			return err
		}

		return nil
	})

	return

}

// adds ubucks from CB
// creates order, transaction into student account, update account balance, update ubuck balance
func addUbuck2Student(db *bolt.DB, clock Clock, userInfo UserInfo, amount decimal.Decimal, reference string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return addUbuck2StudentTx(tx, clock, userInfo, amount, reference)
	})
}

// register order, transactions in students account, transactions i CB
// your functions should sit in a separete file
func addUbuck2StudentTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, reference string) error {
	return pay2StudentTx(tx, clock, userInfo, amount, CurrencyUBuck, reference)

}

func addBuck2Student(db *bolt.DB, clock Clock, userInfo UserInfo, amount decimal.Decimal, currency, reference string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return addBuck2StudentTx(tx, clock, userInfo, amount, currency, reference)
	})
}

func addBuck2StudentTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, currency, reference string) error {
	return pay2StudentTx(tx, clock, userInfo, amount, currency, reference)

}

// addToHolderTx updates balance and adds transaction
// debit means to remove money
func addToHolderTx(holder *bolt.Bucket, account string, transaction Transaction, direction int, negBlock bool) (balance decimal.Decimal, basis decimal.Decimal, errR error) {
	accounts, err := holder.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		errR = err
		return
	}

	accountBucket, err := accounts.CreateBucketIfNotExists([]byte(account))
	if err != nil {
		errR = err
		return
	}

	balanceB := accountBucket.Get([]byte(KeyBalance))
	if balanceB != nil {
		err = balance.UnmarshalText(balanceB)
		if err != nil {
			errR = fmt.Errorf("cannot extract balance for the account %s: %v", account, err)
			return
		}
	} else {
		balance = decimal.Zero
	}

	basisB := accountBucket.Get([]byte(KeyBasis))
	if basisB != nil {
		err = basis.UnmarshalText(basisB)
		if err != nil {
			errR = fmt.Errorf("cannot extract basis for the account %s: %v", account, err)
			return
		}
	} else {
		basis = decimal.Zero
	}

	if direction == OperationCredit {
		if transaction.AmountDest.IsPositive() {
			basisValue := basis.Mul(balance)
			buyValue := transaction.AmountDest.Mul(transaction.XRate)
			numerator := basisValue.Add(buyValue)
			denominator := balance.Add(transaction.AmountDest)
			if !denominator.IsZero() {
				basis = numerator.Div(denominator)
			}
		}

		balance = balance.Add(transaction.AmountDest)
	} else {
		balance = balance.Sub(transaction.AmountSource)
	}

	if balance.IsZero() {
		basis = decimal.Zero
	}

	basisB, err = basis.MarshalText()
	if err != nil {
		errR = err
		return
	}
	err = accountBucket.Put([]byte(KeyBasis), basisB)
	if err != nil {
		errR = err
		return
	}

	if balance.Round(5).Sign() < 0 && negBlock {
		errR = fmt.Errorf("insufficient funds")
		return
	}

	balanceB, err = balance.MarshalText()
	if err != nil {
		errR = err
		return
	}
	err = accountBucket.Put([]byte(KeyBalance), balanceB)
	if err != nil {
		errR = err
		return
	}

	transactions, err := accountBucket.CreateBucketIfNotExists([]byte(KeyTransactions))
	if err != nil {
		errR = err
		return
	}

	tsB, err := transaction.Ts.MarshalText()
	if err != nil {
		errR = err
		return
	}

	oldTransaction := transactions.Get(tsB)
	for oldTransaction != nil {
		transaction.Ts = transaction.Ts.Add(time.Millisecond * 1)
		tsB, err = transaction.Ts.MarshalText()
		if err != nil {
			errR = err
			return
		}
		oldTransaction = transactions.Get(tsB)
	}

	transactionB, err := json.Marshal(transaction)
	if err != nil {
		errR = err
		return
	}

	errR = transactions.Put(tsB, transactionB)

	return
}

func addToCDHolderTx(holder *bolt.Bucket, transaction Transaction, CD CertificateOfDeposit, CD_id string) (err error) {
	accounts, err := holder.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		return
	}

	CDBucket, err := accounts.CreateBucketIfNotExists([]byte(KeyCertificateOfDeposit))
	if err != nil {
		return
	}

	transactions, err := CDBucket.CreateBucketIfNotExists([]byte(KeyTransactions))
	if err != nil {
		return
	}

	tsB, err := transaction.Ts.MarshalText()
	if err != nil {
		return
	}

	oldCD := CDBucket.Get(tsB)
	oldTransaction := transactions.Get(tsB)
	for oldTransaction != nil || oldCD != nil {
		transaction.Ts = transaction.Ts.Add(time.Millisecond * 1)
		CD.Ts = transaction.Ts
		tsB, err = transaction.Ts.MarshalText()
		if err != nil {
			return
		}
		oldTransaction = transactions.Get(tsB)
		oldCD = CDBucket.Get(tsB)
	}

	CDB, err := json.Marshal(CD)
	if err != nil {
		return
	}

	transactionB, err := json.Marshal(transaction)
	if err != nil {
		return
	}

	err = transactions.Put(tsB, transactionB)
	if err != nil {
		return
	}

	if CD_id != "" {
		newId, err := time.Parse(time.RFC3339, CD_id)
		if err != nil {
			return err
		}

		CD.Ts = newId

		tsB, err = newId.MarshalText()
		if err != nil {
			return err
		}

		CDB, err := json.Marshal(CD)
		if err != nil {
			return err
		}

		err = CDBucket.Put(tsB, CDB)

		return err
	}

	err = CDBucket.Put(tsB, CDB)

	return
}

func StudentNetWorth(db *bolt.DB, userName string) (res decimal.Decimal) {
	_ = db.View(func(tx *bolt.Tx) error {
		res = StudentNetWorthTx(tx, userName)
		return nil
	})
	return
}
func StudentNetWorthTx(tx *bolt.Tx, userName string) (res decimal.Decimal) {
	res = decimal.Zero
	student, err := getStudentBucketRx(tx, userName)
	if err != nil {
		return
	}

	userData, err := getUserInLocalStoreTx(tx, userName)
	if err != nil {
		return
	}
	accounts := student.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return
	}

	c := accounts.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {

		account := accounts.Bucket(k)
		if account == nil {
			continue
		}

		if string(k) == KeyCertificateOfDeposit {
			res = res.Add(CDSumTx(tx, account))
			continue
		}

		balance := account.Get([]byte(KeyBalance))
		if balance == nil {
			continue
		}

		var value decimal.Decimal
		err = value.UnmarshalText(balance)
		if err != nil {
			lgr.Printf("ERROR cannot unmarshal balance for ubuck of %s: %v", userName, err)
			return decimal.Zero
		}

		if value.IsZero() {
			continue
		}

		var ubuck decimal.Decimal
		if string(k) != CurrencyUBuck && string(k) != KeyDebt && !strings.Contains(string(k), "@") {
			usd, err := getCrypto(tx, string(k))
			if err != nil {
				lgr.Printf("ERROR cannot get crytpo to ubuck of %s: %v", userName, err)
				return decimal.Zero
			}

			basis, err := getStudentCryptoBasisRx(student, string(k))
			if err != nil {
				lgr.Printf("ERROR cannot get crytpo basis of %s: %v", userName, err)
				return decimal.Zero
			}
			ubuck = cryptoConvert(basis, usd, value)
			if err != nil {
				return
			}
		} else {
			ubuck, _, err = convertRx(tx, userData.SchoolId, string(k), "", value.InexactFloat64())
			if err != nil {
				return
			}
		}

		res = res.Add(ubuck)

	}

	return
}

func CertificateOfDepositIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		student, _ := getStudentBucketRx(tx, userDetails.Name)
		if student == nil {
			needToAdd = true
			return nil
		}

		needToAdd = IsDailyPayNeeded(student, clock)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		student, _ := getStudentBucketRx(tx, userDetails.Name)
		if student == nil {
			needToAdd = true
			return nil
		}

		needToAdd := IsDailyPayNeeded(student, clock)

		if !needToAdd {
			return nil
		}

		accountsBucket := student.Bucket([]byte(KeyAccounts))
		if accountsBucket == nil {
			return fmt.Errorf("cannot find accounts bucket")
		}

		CDS_bucket := accountsBucket.Bucket([]byte(KeyCertificateOfDeposit))
		if CDS_bucket == nil {
			return fmt.Errorf("cannot find CDS bucket")
		}

		c := CDS_bucket.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v == nil {
				continue
			}

			var deposit CertificateOfDeposit
			err := json.Unmarshal(v, &deposit)
			if err != nil {
				return err
			}

			if !deposit.Active {
				continue
			}

			if deposit.Maturity.Before(clock.Now()) {
				deposit.CurrentValue = decimal.NewFromInt32(deposit.Principal).Mul(decimal.NewFromFloat32(deposit.Interest).Pow(decimal.NewFromInt(InterestToTime(deposit.Interest))))
				deposit.RefundValue = deposit.CurrentValue
			} else {
				diff := clock.Now().Sub(deposit.Ts)
				days := math.Floor(diff.Hours() / 24)
				deposit.CurrentValue = decimal.NewFromInt32(deposit.Principal).Mul(decimal.NewFromFloat32(deposit.Interest).Pow(decimal.NewFromFloat(days)))
				interest := timeToInterest(int32(days))
				deposit.RefundValue = (decimal.NewFromInt32(deposit.Principal).Mul(decimal.NewFromFloat32(interest).Pow(decimal.NewFromFloat(days)))).Mul(decimal.NewFromFloat32(.9))
			}

			data, err := json.Marshal(deposit)
			if err != nil {
				return err
			}

			err = CDS_bucket.Put(k, data)
			if err != nil {
				return err
			}

		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR checking college on %s: %v", userDetails.Name, err)
		return false
	}
	return needToAdd
}

func DailyPayIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		student, _ := getStudentBucketRx(tx, userDetails.Name)
		if student == nil {
			needToAdd = true
			return nil
		}
		needToAdd = IsDailyPayNeeded(student, clock)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		student, err := getStudentBucketTx(tx, userDetails.Name)
		if err != nil {
			return err
		}
		needToAdd = IsDailyPayNeeded(student, clock)
		if !needToAdd {
			return nil
		}

		pay := decimal.NewFromFloat32(userDetails.Income)
		haveDebt, _, balance, err := IsDebtNeeded(student, clock)
		if err != nil {
			return err
		}

		payDate, err := clock.Now().Truncate(24 * time.Hour).MarshalText()
		if err != nil {
			return err
		}
		err = student.Put([]byte(KeyDayPayment), payDate)
		if err != nil {
			return fmt.Errorf("cannot save daily payment date: %v", err)
		}

		if haveDebt {
			garnish := pay.Mul(decimal.NewFromFloat32(.75))
			pay = pay.Mul(decimal.NewFromFloat32(.25))
			if garnish.GreaterThan(balance) {
				garnishRemainder := garnish.Sub(balance)
				pay = pay.Add(garnishRemainder)
				garnish = balance.Add(decimal.Zero)
			}
			err = chargeStudentTx(tx, clock, userDetails, garnish, KeyDebt, "Paycheck Garnishment", false)
			if err != nil {
				return err
			}
			return addUbuck2StudentTx(tx, clock, userDetails, pay, "daily payment")
		}
		return addUbuck2StudentTx(tx, clock, userDetails, pay, "daily payment")
	})

	if err != nil {
		lgr.Printf("ERROR daily payment not added to %s, Error: %v", userDetails.Name, err)
		return false
	}
	return true
}

func IsDailyPayNeeded(student *bolt.Bucket, clock Clock) bool {
	dayB := student.Get([]byte(KeyDayPayment))
	if dayB == nil {
		return true
	}

	var day time.Time
	err := day.UnmarshalText(dayB)
	if err != nil {
		return true
	}
	if clock.Now().Truncate(24 * time.Hour).After(day) {
		return true
	}
	return false
}

func CollegeIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		if users == nil {
			return nil
		}

		studentData := users.Get([]byte(userDetails.Name))
		if studentData == nil {
			return nil
		}

		needToAdd, _, _ = IsCollegeNeeded(studentData, clock)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		if users == nil {
			return nil
		}

		studentData := users.Get([]byte(userDetails.Name))
		if studentData == nil {
			return nil
		}

		needToAdd, student, err := IsCollegeNeeded(studentData, clock)
		if err != nil {
			return err
		}

		if !needToAdd {
			return nil
		}

		student.Job = getJobIdRx(tx, KeyCollegeJobs)
		jobDetails := getJobRx(tx, KeyCollegeJobs, student.Job)
		max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
		min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
		diff := max.Sub(min)
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		random := decimal.NewFromFloat32(r.Float32())

		student.Income = float32(random.Mul(diff).Add(min).Floor().InexactFloat64())
		student.CollegeEnd = time.Time{}
		marshal, err := json.Marshal(student)
		if err != nil {
			return err
		}

		err = users.Put([]byte(userDetails.Name), marshal)
		if err != nil {
			return fmt.Errorf("cannot save event date: %v", err)
		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR checking college on %s: %v", userDetails.Name, err)
		return false
	}
	return needToAdd
}

func IsCollegeNeeded(studentData []byte, clock Clock) (needed bool, student UserInfo, err error) {
	err = json.Unmarshal(studentData, &student)
	if err != nil {
		return false, student, err
	}
	if !student.CollegeEnd.IsZero() && clock.Now().After(student.CollegeEnd) {
		return true, student, err
	}

	return false, student, err

}

func CareerIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		if users == nil {
			return nil
		}

		studentData := users.Get([]byte(userDetails.Name))
		if studentData == nil {
			return nil
		}

		needToAdd, _, _ = IsCareerNeeded(studentData, clock)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte(KeyUsers))
		if users == nil {
			return nil
		}

		studentData := users.Get([]byte(userDetails.Name))
		if studentData == nil {
			return nil
		}

		needToAdd, student, err := IsCareerNeeded(studentData, clock)
		if err != nil {
			return err
		}

		if !needToAdd {
			return nil
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if student.College && student.CollegeEnd.IsZero() {
			student.Job = getJobIdRx(tx, KeyCollegeJobs)
			jobDetails := getJobRx(tx, KeyCollegeJobs, student.Job)
			max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
			min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
			diff := max.Sub(min)
			random := decimal.NewFromFloat32(r.Float32())
			student.Income = float32(random.Mul(diff).Add(min).Floor().InexactFloat64())
		} else {
			student.Job = getJobIdRx(tx, KeyJobs)
			jobDetails := getJobRx(tx, KeyJobs, student.Job)
			max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
			min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
			diff := max.Sub(min)
			random := decimal.NewFromFloat32(r.Float32())
			student.Income = float32(random.Mul(diff).Add(min).Floor().InexactFloat64())
		}

		student.TransitionEnd = time.Time{}
		student.CareerTransition = false
		marshal, err := json.Marshal(student)
		if err != nil {
			return err
		}

		err = users.Put([]byte(userDetails.Name), marshal)
		if err != nil {
			return fmt.Errorf("cannot save event date: %v", err)
		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR checking college on %s: %v", userDetails.Name, err)
		return false
	}
	return needToAdd
}

func IsCareerNeeded(studentData []byte, clock Clock) (needed bool, student UserInfo, err error) {
	err = json.Unmarshal(studentData, &student)
	if err != nil {
		return false, student, err
	}
	if student.CareerTransition && !student.TransitionEnd.IsZero() && clock.Now().After(student.TransitionEnd) {
		return true, student, err
	}

	return false, student, err

}

func DebtIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		student, err := getStudentBucketRx(tx, userDetails.Name)
		if err != nil {
			return nil
		}

		needToAdd, _, _, _ = IsDebtNeeded(student, clock)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		student, err := getStudentBucketRx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		needToAdd, day, balance, err := IsDebtNeeded(student, clock)

		if !needToAdd {
			return nil
		}

		days := decimal.NewFromFloat32(float32(clock.Now().Truncate(24*time.Hour).Sub(day).Hours() / 24))
		interest := decimal.NewFromFloat32(LoanRate)
		compoundedInterest := interest.Pow(days)
		compoundedInterest = compoundedInterest.Sub(decimal.NewFromInt32(1))
		change := balance.Mul(compoundedInterest)
		err = pay2StudentTx(tx, clock, userDetails, change, KeyDebt, days.String()+" days compound interest")
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		lgr.Printf("ERROR checking debt on %s: %v", userDetails.Name, err)
		return false
	}
	return needToAdd
}

func IsDebtNeeded(student *bolt.Bucket, clock Clock) (needed bool, day time.Time, balance decimal.Decimal, err error) {
	accounts := student.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return false, day, balance, err
	}

	debt := accounts.Bucket([]byte(KeyDebt))
	if debt == nil {
		return false, day, balance, err
	}

	balanceb := debt.Get([]byte(KeyBalance))
	err = balance.UnmarshalText(balanceb)
	if err != nil {
		return
	}

	if balance.IsZero() {
		return
	}

	dayB := student.Get([]byte(KeyDayPayment))

	err = day.UnmarshalText(dayB)
	if err != nil {
		return
	}

	if clock.Now().Truncate(24 * time.Hour).After(day) {
		return true, day, balance, err
	}

	return
}

func EventIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		student, _ := getStudentBucketRx(tx, userDetails.Name)
		if student == nil {
			return nil
		}

		needToAdd, _ = IsEventNeeded(student, clock, false)
		return nil
	})

	if !needToAdd {
		return false
	}
	err := db.Update(func(tx *bolt.Tx) error {
		student, err := getStudentBucketTx(tx, userDetails.Name)
		if err != nil {
			return err
		}

		needToAdd, err = IsEventNeeded(student, clock, true)
		if err != nil {
			return err
		}

		if !needToAdd {
			return nil
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		days := r.Intn(5) + 4

		eventTime := clock.Now().AddDate(0, 0, days).Truncate(24 * time.Hour)
		eventDate, err := eventTime.MarshalText()
		if err != nil {
			return err
		}
		err = student.Put([]byte(KeyDayEvent), eventDate)
		if err != nil {
			return fmt.Errorf("cannot save event date: %v", err)
		}

		students, _, err := getSchoolStudentsTx(tx, userDetails)
		if err != nil {
			return err
		}

		change, err := makeEvent(students, userDetails)
		if err != nil {
			return err
		}

		if change.IsPositive() {
			return addUbuck2StudentTx(tx, clock, userDetails, change, "Event: "+getEventIdRx(tx, KeyPEvents))
		}
		return chargeStudentUbuckTx(tx, clock, userDetails, change.Abs(), "Event: "+getEventIdRx(tx, KeyNEvents), false)
	})

	if err != nil {
		lgr.Printf("ERROR event not added to %s: %v", userDetails.Name, err)
		return false
	}
	return needToAdd
}

func getEventDescription(db *bolt.DB, typeKey string, idKey string) (description string) {
	_ = db.View(func(tx *bolt.Tx) error {
		description = getEventDescriptionRx(tx, typeKey, idKey)
		return nil
	})

	return
}

func getEventDescriptionRx(tx *bolt.Tx, typeKey string, idKey string) string {
	events := tx.Bucket([]byte(typeKey))

	var event EventRequest
	err := json.Unmarshal(events.Get([]byte(idKey)), &event)
	if err != nil {
		return "error retrieving negative event"
	}

	return event.Description
}

func getEventId(db *bolt.DB, key string) (id string) {
	_ = db.View(func(tx *bolt.Tx) error {
		id = getEventIdRx(tx, key)
		return nil
	})

	return
}

func getEventIdRx(tx *bolt.Tx, key string) string {
	events := tx.Bucket([]byte(key))
	bucketStats := events.Stats()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pick := r.Intn(bucketStats.KeyN)
	c := events.Cursor()
	i := 0
	for k, _ := c.First(); k != nil && i <= pick; k, _ = c.Next() {
		if i != pick {
			i++
			continue
		}

		i++

		return string(k)

	}

	return ""
}

func makeEvent(students []openapi.UserNoHistory, userDetails UserInfo) (change decimal.Decimal, err error) {
	multiplier := decimal.NewFromFloat(.4)
	one := decimal.NewFromInt(1)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	random := decimal.NewFromFloat32(r.Float32())
	count := len(students)
	for i := range students {
		if students[i].Id == userDetails.Name {
			var topNetWorth decimal.Decimal
			if count >= 4 {
				topNetWorth = decimal.NewFromFloat32((students[0].NetWorth + students[1].NetWorth + students[2].NetWorth + students[3].NetWorth) / 4)
			} else {
				topNetWorth = decimal.NewFromFloat32(students[0].NetWorth)
			}

			chance := decimal.NewFromFloat32(float32((i + 1)) / float32(count)) //higher chance means better good events
			max := topNetWorth.Mul(multiplier).Mul(chance)
			min := topNetWorth.Mul(multiplier).Mul(chance.Sub(one))
			change = max.Sub(min).Mul(random).Add(min).Floor()

			if change.IsNegative() && change.Abs().GreaterThanOrEqual(decimal.NewFromFloat32(students[i].NetWorth).Mul(decimal.NewFromFloat32(.5))) {
				change = decimal.NewFromFloat32(students[i].NetWorth).Mul(decimal.NewFromFloat32(-.2))
			}
			return
		}
	}

	return change, fmt.Errorf("did not find student in slice")
}

func IsEventNeeded(student *bolt.Bucket, clock Clock, tx bool) (bool, error) {
	dayB := student.Get([]byte(KeyDayEvent))
	if dayB == nil && tx {
		days := rand.Intn(5) + 4
		eventTime := clock.Now().AddDate(0, 0, days).Truncate(24 * time.Hour)
		eventDate, err := eventTime.MarshalText()
		if err != nil {
			return true, err
		}
		err = student.Put([]byte(KeyDayEvent), eventDate)
		if err != nil {
			return false, fmt.Errorf("cannot save event date: %v", err)
		}

		return false, nil
	}

	if dayB == nil && !tx {
		return true, nil
	}

	var day time.Time
	err := day.UnmarshalText(dayB)
	if err != nil {
		return true, err
	}
	if clock.Now().Truncate(24 * time.Hour).After(day) {
		return true, nil
	}
	return false, nil
}

func getCbTx(tx *bolt.Tx, schoolId string) (cb *bolt.Bucket, err error) {
	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return nil, err
	}
	return school.CreateBucketIfNotExists([]byte(KeyCB))
}

func getCbRx(tx *bolt.Tx, schoolId string) (cb *bolt.Bucket, err error) {
	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return nil, err
	}
	cb = school.Bucket([]byte(KeyCB))
	if cb == nil {
		return nil, fmt.Errorf("cannot find CB for school %s", schoolId)
	}
	return cb, nil
}

func chargeStudentUbuck(db *bolt.DB, clock Clock, userDetails UserInfo, amount decimal.Decimal, reference string, sPurchase bool) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return chargeStudentUbuckTx(tx, clock, userDetails, amount, reference, sPurchase)
	})
}
func chargeStudentUbuckTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, amount decimal.Decimal, reference string, sPurchase bool) (err error) {
	return chargeStudentTx(tx, clock, userDetails, amount, CurrencyUBuck, reference, sPurchase)

}

func pay2Student(db *bolt.DB, clock Clock, userInfo UserInfo, amount decimal.Decimal, currency string, reference string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return pay2StudentTx(tx, clock, userInfo, amount, currency, reference)
	})
}

func pay2StudentTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, currency string, reference string) error {
	if userInfo.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}
	res, err := isCurrencyTx(tx, userInfo.SchoolId, currency)
	if err != nil || !res {
		return fmt.Errorf("currency %s is not supported, %v", currency, err)
	}

	ts := clock.Now().Truncate(time.Millisecond)

	transaction := Transaction{
		Ts:             ts,
		Source:         "",
		Destination:    userInfo.Name,
		CurrencySource: currency,
		CurrencyDest:   currency,
		AmountSource:   amount,
		AmountDest:     amount,
		XRate:          decimal.NewFromFloat32(1.0),
		Reference:      reference,
	}
	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(student, currency, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	if currency != CurrencyUBuck && currency != KeyDebt {
		mma, err := addStepTx(tx, userInfo.SchoolId, currency, float32(amount.InexactFloat64()))
		if err != nil {
			return err
		}

		_, err = modifyMmaTx(tx, userInfo.SchoolId, currency, transaction.Ts, mma, clock)
		if err != nil {
			return err
		}
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(cb, currency, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func studentConvertTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, from, to, reference string, charge bool) (err error) {
	if userInfo.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)

	target := to
	toDetails := UserInfo{
		LastName: "Debt",
	}

	fromDetails := UserInfo{
		LastName: "UBuck",
	}
	if to == KeyDebt {
		target = ""
	} else if to == CurrencyUBuck {
		toDetails = UserInfo{
			LastName: "UBuck",
		}
	} else {
		toDetails, err = getUserInLocalStoreTx(tx, to)
		if err != nil {
			return err
		}

		if toDetails.Settings.CurrencyLock {
			return fmt.Errorf(toDetails.LastName + " bucks are locked by the teacher")
		}
	}

	if from != CurrencyUBuck && from != KeyDebt {
		fromDetails, err = getUserInLocalStoreTx(tx, from)
		if err != nil {
			return err
		}

		if fromDetails.Settings.CurrencyLock {
			return fmt.Errorf(fromDetails.LastName + " bucks are locked by the teacher")
		}
	}

	converted, xRate, err := convertRx(tx, userInfo.SchoolId, from, target, amount.InexactFloat64())
	if err != nil {
		return err
	}

	if charge {
		converted = converted.Mul(decimal.NewFromFloat(2 - keyCharge))
	}

	if reference == "" {
		reference = fromDetails.LastName + " to " + toDetails.LastName
	}

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Name,
		Destination:    userInfo.Name,
		CurrencySource: from,
		CurrencyDest:   to,
		AmountSource:   amount,
		AmountDest:     converted,
		XRate:          xRate,
		Reference:      reference,
		FromSource:     true,
	}

	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}
	if to != KeyDebt || charge {
		_, _, err = addToHolderTx(student, from, transaction, OperationDebit, true)
		if err != nil {
			return err
		}
	}

	transaction.FromSource = false

	if to == KeyDebt && charge {
		transaction.AmountDest = transaction.AmountDest.Neg()
	}

	_, _, err = addToHolderTx(student, to, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	_, _, err = addToHolderTx(cb, to, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, _, err = addToHolderTx(cb, from, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func chargeStudent(db *bolt.DB, clock Clock, userDetails UserInfo, amount decimal.Decimal, currency string, reference string, sPurchase bool) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return chargeStudentTx(tx, clock, userDetails, amount, currency, reference, sPurchase)
	})
}

// sPurchase is asking if the student is making this transaction or the teacher
func chargeStudentTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, amount decimal.Decimal, currency string, reference string, sPurchase bool) (err error) {
	if userDetails.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)

	transaction := Transaction{
		Ts:             ts,
		Source:         userDetails.Name,
		Destination:    "",
		CurrencySource: currency,
		CurrencyDest:   currency,
		AmountSource:   amount,
		AmountDest:     amount,
		XRate:          decimal.NewFromFloat32(1.0),
		Reference:      reference,
	}
	student, err := getStudentBucketTx(tx, userDetails.Name)
	if err != nil {
		return fmt.Errorf("can't find student, error: " + err.Error())
	}
	_, _, err = addToHolderTx(student, currency, transaction, OperationDebit, true)
	if err != nil {
		if err.Error() == "insufficient funds" && !sPurchase {
			if strings.Contains(reference, "Event: ") {
				err := studentConvertTx(tx, clock, userDetails, amount, currency, KeyDebt, reference, false)
				if err != nil {
					return err
				}
			} else {
				err := studentConvertTx(tx, clock, userDetails, amount, currency, KeyDebt, reference, false) //this was "" I may need to change it back
				if err != nil {
					return err
				}
			}

			return nil

		}

		return err

	}
	cb, err := getCbTx(tx, userDetails.SchoolId)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(cb, currency, transaction, OperationCredit, false)
	if err != nil {
		return err
	}

	return nil
}

/*
check if the currency exists in given school
*/
func isCurrencyTx(tx *bolt.Tx, schoolId string, currency string) (bool, error) {
	if currency == CurrencyUBuck || currency == KeyDebt {
		return true, nil
	}

	schools := tx.Bucket([]byte(KeySchools))
	if schools == nil {
		return false, fmt.Errorf("schools  does not exist")
	}

	school := schools.Bucket([]byte(schoolId))

	if school == nil {
		return false, fmt.Errorf("student not found")
	}
	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return false, fmt.Errorf("teachers does not exist")
	}
	teacher := teachers.Bucket([]byte(currency))
	if teacher == nil {
		return false, fmt.Errorf("teacher does not exist")
	}
	return true, nil
}

func getStudentUbuck(db *bolt.DB, userDetails UserInfo) (uBucks openapi.ResponseSearchStudentUbuck, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		uBucks, err = getStudentUbuckRx(tx, userDetails)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func getStudentUbuckRx(tx *bolt.Tx, userDetails UserInfo) (resp openapi.ResponseSearchStudentUbuck, err error) {
	student, err := getStudentBucketRx(tx, userDetails.Name)
	if student == nil {
		return resp, err
	}

	accounts := student.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return resp, fmt.Errorf("cannot find Buck Accounts")
	}

	ubuck := accounts.Bucket([]byte(CurrencyUBuck))
	if ubuck == nil {
		return resp, fmt.Errorf("cannot find ubuck")
	}

	var balance decimal.Decimal
	balanceB := ubuck.Get([]byte(KeyBalance))
	err = balance.UnmarshalText(balanceB)
	if err != nil {
		return resp, fmt.Errorf("cannot extract balance for the account %s: %v", userDetails.Name, err)
	}
	balanceF, _ := balance.Float64()
	resp = openapi.ResponseSearchStudentUbuck{
		Value: float32(balanceF),
	}

	return
}

func getStudentBuck(db *bolt.DB, userDetails UserInfo, teacherId string) (resp openapi.ResponseSearchStudentUbuck, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getStudentBuckRx(tx, userDetails, teacherId)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func getStudentBuckRx(tx *bolt.Tx, userDetails UserInfo, teacherId string) (resp openapi.ResponseSearchStudentUbuck, err error) {
	student, err := getStudentBucketRx(tx, userDetails.Name)
	if student == nil {
		return resp, err
	}

	accounts := student.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return resp, fmt.Errorf("cannot find Buck Accounts")
	}

	buck := accounts.Bucket([]byte(teacherId))
	if buck == nil {
		return resp, fmt.Errorf("cannot find ubuck")
	}

	var balance decimal.Decimal
	balanceB := buck.Get([]byte(KeyBalance))
	err = balance.UnmarshalText(balanceB)
	if err != nil {
		return resp, fmt.Errorf("cannot extract balance for the account %s: %v", userDetails.Name, err)
	}
	balanceF, _ := balance.Float64()
	resp = openapi.ResponseSearchStudentUbuck{
		Value: float32(balanceF),
	}

	return
}

func getStudentAuctions(db *bolt.DB, userDetails UserInfo) (auctions []openapi.Auction, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		auctions, err = getStudentAuctionsRx(tx, userDetails)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func getStudentAuctionsRx(tx *bolt.Tx, userDetails UserInfo) (auctions []openapi.Auction, err error) {
	auctionsBucket, err := getAuctionsTx(tx, userDetails)
	if err != nil {
		return
	}

	c := auctionsBucket.Cursor()

	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		if v == nil {
			continue
		}

		auctionData := auctionsBucket.Get(k)
		if auctionData == nil {
			return auctions, fmt.Errorf("cannot find auction bucket")
		}

		var auction openapi.Auction
		err := json.Unmarshal(auctionData, &auction)
		if err != nil {
			return nil, err
		}

		classes, err := getStudentClassesRx(tx, userDetails)
		if err != nil {
			return auctions, fmt.Errorf("cannot get student classes %s: %v", userDetails.Name, err)
		}

		for _, k := range auction.Visibility {
			for _, class := range classes {
				if class.Id == string(k) || string(k) == KeyEntireSchool {

					ownerDetails, err := getUserInLocalStoreTx(tx, auction.OwnerId.Id)
					if err != nil {
						return nil, err
					}

					auction.OwnerId = openapi.AuctionOwnerId{
						LastName: ownerDetails.LastName,
						Id:       ownerDetails.Name,
					}

					winnerDetails, err := getUserInLocalStoreTx(tx, auction.WinnerId.Id)
					if err != nil {
						auction.WinnerId = openapi.AuctionWinnerId{
							FirstName: "nil",
							LastName:  "nil",
							Id:        "nil",
						}
					} else {
						auction.WinnerId = openapi.AuctionWinnerId{
							FirstName: winnerDetails.FirstName,
							LastName:  winnerDetails.LastName,
							Id:        winnerDetails.Name,
						}
					}

					auctions = append(auctions, auction)
					break

				}
			}
		}

	}

	return auctions, nil
}

func getStudentClasses(db *bolt.DB, userDetails UserInfo) (classes []openapi.Class, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		classes, err = getStudentClassesRx(tx, userDetails)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func getStudentClassesRx(tx *bolt.Tx, userDetails UserInfo) (classes []openapi.Class, err error) {
	school, err := getSchoolBucketRx(tx, userDetails)
	if err != nil {
		return classes, fmt.Errorf("cannot get schools %s: %v", userDetails.Name, err)
	}

	schoolClasses := school.Bucket([]byte(KeyClasses))
	if schoolClasses == nil {
		return classes, fmt.Errorf("cannot get school classes %s: %v", userDetails.Name, err)
	}

	s := schoolClasses.Cursor()
	for k, v := s.First(); k != nil; k, v = s.Next() {
		if v != nil {
			continue
		}

		class := schoolClasses.Bucket(k)
		if class == nil {
			return classes, fmt.Errorf("cannot get school classes %s: %v", userDetails.Name, err)
		}

		iClass := getClassMembershipTx(k, class, userDetails)
		if iClass.AddCode == "" {
			continue
		}

		classes = append(classes, iClass)
	}

	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return classes, fmt.Errorf("cannot get teachers %s: %v", userDetails.Name, err)
	}

	t := teachers.Cursor()
	for k, v := t.First(); k != nil; k, v = t.Next() {
		if v != nil {
			continue
		}

		teacher := teachers.Bucket(k)
		if teacher == nil {
			return classes, fmt.Errorf("cannot get teacher %s: %v", userDetails.Name, err)
		}

		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			continue
		}

		c := classesBucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			class := classesBucket.Bucket(k)
			if class == nil {
				return classes, fmt.Errorf("cannot get teacher classes %s: %v", userDetails.Name, err)
			}

			iClass := getClassMembershipTx(k, class, userDetails)
			if iClass.AddCode == "" {
				continue
			}

			classes = append(classes, iClass)
		}

	}

	return

}

func getClassMembershipTx(key []byte, class *bolt.Bucket, userDetails UserInfo) (classResp openapi.Class) {
	students := class.Bucket([]byte(KeyStudents))
	if students == nil {
		return
	}

	student := students.Get([]byte(userDetails.Name))
	if student != nil {
		classResp = openapi.Class{
			Id:      string(key),
			OwnerId: string(class.Get([]byte(KeyOwnerId))),
			Period:  btoi32(class.Get([]byte(KeyPeriod))),
			Name:    string(class.Get([]byte(KeyName))),
			AddCode: string(class.Get([]byte(KeyAddCode))),
		}
	}

	return
}

func getStudentTransactionsTx(tx *bolt.Tx, bAccounts *bolt.Bucket, student UserInfo) (resp []openapi.ResponseBuckTransaction, err error) {
	trans := make([]Transaction, 0)

	c := bAccounts.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {

		if string(k) != CurrencyUBuck && string(k) != KeyDebt && !strings.Contains(string(k), "@") {
			continue
		}

		buck := bAccounts.Bucket(k)

		transactions := buck.Bucket([]byte(KeyTransactions))
		if transactions == nil {
			continue
		}

		trans, err = getBuckTransactionsTx(transactions, trans, student, string(k))
		if err != nil {
			return resp, err
		}
	}

	for _, tran := range trans {
		response, err := transactionToResponseBuckTransactionTx(tx, tran)
		if err != nil {
			return resp, err
		}

		resp = append(resp, response)
	}

	return

}

func getBuckTransactionsTx(transactions *bolt.Bucket, trans []Transaction, student UserInfo, accountId string) (resp []Transaction, err error) {
	c := transactions.Cursor()
	var balance decimal.Decimal
OUTER:
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		v := transactions.Get(k)
		newAdd := parseTransactionStudent(v, student, accountId)
		balance = balance.Add(newAdd.Net)
		newAdd.Balance = balance

		candidateTime, err := time.Parse(time.RFC3339, string(k))
		if err != nil {
			return resp, fmt.Errorf("Cannot parse time")
		}

		if len(trans) == 0 {
			trans = append(trans, newAdd)
			continue
		}

		for i, residentTime := range trans {

			if candidateTime.After(residentTime.Ts) {
				trans = insert(trans, i, newAdd)
				continue OUTER
			}

		}

		if len(trans) < 50 {
			trans = append(trans, newAdd)
		}
	}

	resp = trans

	return

}

func parseTransactionStudent(transData []byte, user UserInfo, accountId string) (trans Transaction) {
	err := json.Unmarshal(transData, &trans)
	if err != nil {
		lgr.Printf("ERROR cannot unmarshal trans")
		return
	}

	if trans.CurrencySource == user.Name { //this is a teacher trans
		if trans.Destination == "" {
			trans.Net = trans.AmountSource.Neg()
		} else {
			trans.Net = trans.AmountDest
		}
	} else if trans.Destination == user.Name && trans.Source == user.Name {
		if trans.CurrencyDest == accountId {
			trans.Net = trans.AmountDest
			trans.Ubuck = trans.AmountSource
		} else {
			trans.Net = trans.AmountSource.Neg()
			trans.Ubuck = trans.AmountDest
		}
	} else if trans.Destination == user.Name {
		trans.Net = trans.AmountDest
	} else {
		trans.Net = trans.AmountSource.Neg()
	}

	return
}

func transactionToResponseBuckTransactionTx(tx *bolt.Tx, trans Transaction) (resp openapi.ResponseBuckTransaction, err error) {
	xrate, _ := trans.XRate.Float64()
	amount, _ := trans.Net.Float64()
	balance, _ := trans.Balance.Float64()
	var buckName string

	if trans.FromSource {

		if trans.CurrencySource == CurrencyUBuck {
			buckName = "UBuck"
		} else if trans.CurrencySource == KeyDebt {
			buckName = KeyDebt
		} else if strings.Contains(trans.CurrencySource, "@") {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencySource)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"

		} else {
			buckName = trans.CurrencySource
		}
	} else {

		if trans.CurrencyDest == CurrencyUBuck {
			buckName = "UBuck"
		} else if trans.CurrencyDest == KeyDebt {
			buckName = KeyDebt
		} else if strings.Contains(trans.CurrencyDest, "@") {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencyDest)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"
		} else {
			buckName = trans.CurrencyDest
		}

	}

	resp = openapi.ResponseBuckTransaction{
		Balance:     float32(balance),
		Description: trans.Reference,
		Conversion:  float32(xrate),
		Amount:      float32(amount),
		Name:        buckName,
		CreatedAt:   trans.Ts,
	}

	if strings.Contains(resp.Description, "Event: ") {
		typeKey := KeyPEvents
		if resp.Amount <= 0 || trans.CurrencyDest == KeyDebt {
			typeKey = KeyNEvents
		}

		idKey := resp.Description[7:]
		resp.Description = getEventDescriptionRx(tx, typeKey, idKey)
	}

	return
}

func insert(a []Transaction, index int, value Transaction) []Transaction {
	if len(a) == index { // nil or empty slice or after last element
		return append(a, value)
	}
	a = append(a[:index+1], a[index:]...) // index < len(a)
	a[index] = value
	return a
}

func getStudentCryptosRx(tx *bolt.Tx, userDetails UserInfo) (resp []CryptoDecimal, err error) {
	studentBucket, err := getStudentBucketRx(tx, userDetails.Name)
	if err != nil {
		return
	}

	accountsBucket := studentBucket.Bucket([]byte(KeyAccounts))
	if accountsBucket == nil {
		return resp, fmt.Errorf("Cannot find accounts bucket")
	}

	c := accountsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {

		if string(k) == CurrencyUBuck || string(k) == KeyDebt || strings.Contains(string(k), "@") {
			continue
		}

		account := accountsBucket.Bucket(k)

		var balance decimal.Decimal
		balanceB := account.Get([]byte(KeyBalance))
		if balanceB != nil {
			err = balance.UnmarshalText(balanceB)
			if err != nil {
				return resp, err
			}
		} else {
			balance = decimal.Zero
		}

		var basis decimal.Decimal
		basisB := account.Get([]byte(KeyBasis))
		if basisB != nil {
			err = basis.UnmarshalText(basisB)
			if err != nil {
				return resp, err
			}
		} else {
			basis = decimal.Zero
		}

		if balance.IsPositive() {
			var crypto CryptoDecimal
			crypto.Name = string(k)
			usd, err := getCrypto(tx, string(k))
			if err != nil {
				return resp, err
			}
			crypto.CurrentPrice = usd.Round(4)
			crypto.Basis = basis.Round(4)
			crypto.Quantity = balance.Round(4)
			resp = append(resp, crypto)
		}
	}

	return
}

func getStudentCrypto(db *bolt.DB, userDetails UserInfo, crypto string) (resp CryptoDecimal, accountsBucket *bolt.Bucket, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, accountsBucket, err = getStudentCryptoRx(tx, userDetails, crypto)
		return err
	})

	return
}

func getStudentCryptoRx(tx *bolt.Tx, userDetails UserInfo, crypto string) (resp CryptoDecimal, accountsBucket *bolt.Bucket, err error) {
	studentBucket, err := getStudentBucketRx(tx, userDetails.Name)
	if err != nil {
		return
	}

	accountsBucket = studentBucket.Bucket([]byte(KeyAccounts))
	if accountsBucket == nil {
		return resp, accountsBucket, fmt.Errorf("Cannot find accounts bucket")
	}

	accountData := accountsBucket.Bucket([]byte(crypto))
	if accountData == nil {
		return
	}

	resp.Name = crypto
	var basis decimal.Decimal
	var balance decimal.Decimal
	basisData := accountData.Get([]byte(KeyBasis))
	err = json.Unmarshal(basisData, &basis)
	if err != nil {
		return
	}

	balanceData := accountData.Get([]byte(KeyBalance))
	err = json.Unmarshal(balanceData, &balance)
	if err != nil {
		return
	}

	resp.Basis = basis
	resp.Quantity = balance

	return
}

func getCryptoForStudentRequest(db *bolt.DB, userDetails UserInfo, crypto string) (resp openapi.ResponseCrypto, err error) {
	crypto = strings.ToLower(crypto)
	usd, err := getOrUpdateCrypto(db, crypto)
	if err != nil {
		return
	}

	ubuck, err := getStudentUbuck(db, userDetails)
	if err != nil {
		return
	}

	studentCrypto, _, err := getStudentCrypto(db, userDetails, crypto)
	if err != nil {
		return
	}

	resp = openapi.ResponseCrypto{
		Searched: crypto,
		Usd:      usd,
		Owned:    float32(studentCrypto.Quantity.InexactFloat64()),
		UBuck:    ubuck.Value,
		Basis:    float32(studentCrypto.Basis.InexactFloat64()),
	}

	return
}

func getOrUpdateCrypto(db *bolt.DB, crypto string) (usd float32, err error) {
	needToAdd := false
	var cryptoInfoOut openapi.CryptoCb
	err = db.View(func(tx *bolt.Tx) error {
		toUpdate, cryptoInfo, err := isCryptoNeeded(tx, crypto)
		if err != nil {
			return err
		}
		if toUpdate {
			needToAdd = true
			return nil
		}

		cryptoInfoOut = cryptoInfo
		return nil
	})

	if !needToAdd || err != nil {
		return cryptoInfoOut.Usd, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		cryptos, err := tx.CreateBucketIfNotExists([]byte(KeyCryptos))
		if err != nil {
			return err
		}

		var cryptoInfo openapi.CryptoCb

		usd, err := getCrypto(tx, crypto)
		if err != nil {
			return err
		}

		cryptoInfo.Usd = float32(usd.InexactFloat64())
		cryptoInfo.UpdatedAt = time.Now().Truncate(time.Second)
		marshal, err := json.Marshal(cryptoInfo)
		if err != nil {
			return err
		}
		err = cryptos.Put([]byte(crypto), marshal)
		if err != nil {
			return err
		}

		cryptoInfoOut = cryptoInfo
		return nil

	})

	return cryptoInfoOut.Usd, err
}

func isCryptoNeeded(tx *bolt.Tx, crypto string) (needToAdd bool, cryptoInfo openapi.CryptoCb, err error) {
	cryptos := tx.Bucket([]byte(KeyCryptos))
	if cryptos == nil {
		return true, cryptoInfo, nil
	}

	cryptoData := cryptos.Get([]byte(crypto))
	if cryptoData == nil {
		return true, cryptoInfo, nil
	}

	err = json.Unmarshal(cryptoData, &cryptoInfo)
	if err != nil {
		return false, cryptoInfo, err
	}

	if time.Now().Sub(cryptoInfo.UpdatedAt) > time.Minute {
		return true, cryptoInfo, nil
	}

	return
}

func getCrypto(tx *bolt.Tx, crypto string) (usd decimal.Decimal, err error) {
	crypto = strings.ToLower(crypto)
	cryptoBucket := tx.Bucket([]byte(KeyCryptos))
	cryptoData := cryptoBucket.Get([]byte(crypto))
	if cryptoData == nil {
		return usd, fmt.Errorf("can't find the Crypto")
	}

	var cryptoRecord openapi.CryptoCb
	err = json.Unmarshal(cryptoData, &cryptoRecord)
	if err != nil {
		return
	}

	return decimal.NewFromFloat32(cryptoRecord.Usd), err
}

func cryptoTransaction(db *bolt.DB, clock Clock, userDetails UserInfo, body openapi.RequestCryptoConvert) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return cryptoTransactionTx(tx, clock, userDetails, body)
	})

}

func cryptoTransactionTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, body openapi.RequestCryptoConvert) (err error) {
	buy := decimal.NewFromFloat32(body.Buy)
	sell := decimal.NewFromFloat32(body.Sell)
	if body.Buy > 0 {
		err = ubuckToCrypto(tx, clock, userDetails, buy, body.Name)
	} else {
		err = cryptoToUbuck(tx, clock, userDetails, sell, body.Name)
	}

	return
}

func ubuckToCrypto(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, to string) (err error) {
	to = strings.ToLower(to)
	if userInfo.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)

	usd, err := getCrypto(tx, to)
	if err != nil {
		return err
	}

	ubuck := amount.Mul(usd).Mul(decimal.NewFromFloat(keyCharge))

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Name,
		Destination:    userInfo.Name,
		CurrencySource: CurrencyUBuck,
		CurrencyDest:   to,
		AmountSource:   ubuck,
		AmountDest:     amount,
		XRate:          usd,
		Reference:      "Ubuck to " + to,
		FromSource:     true,
	}

	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(student, CurrencyUBuck, transaction, OperationDebit, true)
	if err != nil {
		return err
	}

	transaction.FromSource = false

	_, _, err = addToHolderTx(student, to, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	_, _, err = addToHolderTx(cb, to, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func cryptoToUbuck(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, from string) (err error) {
	from = strings.ToLower(from)
	if userInfo.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)

	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}

	usd, err := getCrypto(tx, from)
	if err != nil {
		return err
	}

	basis, err := getStudentCryptoBasisRx(student, from)
	if err != nil {
		return err
	}

	ubuck := cryptoConvert(basis, usd, amount)
	charge := 1 - (keyCharge - 1)
	ubuck = ubuck.Mul(decimal.NewFromFloat(charge))

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Name,
		Destination:    userInfo.Name,
		CurrencySource: strings.ToLower(from),
		CurrencyDest:   CurrencyUBuck,
		AmountSource:   amount,
		AmountDest:     ubuck,
		XRate:          usd,
		Reference:      from + " to Ubuck",
		FromSource:     true,
	}

	_, _, err = addToHolderTx(student, from, transaction, OperationDebit, true)
	if err != nil {
		return err
	}

	transaction.FromSource = false

	_, _, err = addToHolderTx(student, CurrencyUBuck, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, _, err = addToHolderTx(cb, from, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func getStudentCryptoBasisRx(holder *bolt.Bucket, from string) (basis decimal.Decimal, err error) {
	accounts := holder.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		basis = decimal.Zero
		return
	}

	accountBucket := accounts.Bucket([]byte(from))
	if accountBucket == nil {
		basis = decimal.Zero
		return
	}

	basisB := accountBucket.Get([]byte(KeyBasis))
	if basisB != nil {
		err = basis.UnmarshalText(basisB)
		if err != nil {
			err = fmt.Errorf("cannot extract basis for the account %s: %v", from, err)
			return
		}
	} else {
		basis = decimal.Zero
	}

	return
}

func getStudentCryptoTransactionsRx(tx *bolt.Tx, userDetails UserInfo) (resp []openapi.ResponseCryptoTransaction, err error) {
	student, err := getStudentBucketRx(tx, userDetails.Name)
	if err != nil {
		return
	}

	trans := make([]Transaction, 0)
	accounts := student.Bucket([]byte(KeyAccounts))
	c := accounts.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		if string(k) == CurrencyUBuck || string(k) == KeyDebt || strings.Contains(string(k), "@") || string(k) == KeyCertificateOfDeposit {
			continue
		}

		account := accounts.Bucket(k)
		if account == nil {
			continue
		}

		transactions := account.Bucket([]byte(KeyTransactions))
		if transactions == nil {
			return resp, fmt.Errorf("Cannot get transactions")
		}

		trans, err = getBuckTransactionsTx(transactions, trans, userDetails, string(k))
		if err != nil {
			return resp, err
		}
	}

	for _, tran := range trans {
		response, err := transactionToResponseCryptoTransactionRx(tx, tran)
		if err != nil {
			return resp, err
		}

		resp = append(resp, response)
	}

	return
}

func transactionToResponseCryptoTransactionRx(tx *bolt.Tx, trans Transaction) (resp openapi.ResponseCryptoTransaction, err error) {
	xrate, _ := trans.XRate.Float64()
	amount, _ := trans.Net.Float64()
	balance, _ := trans.Balance.Float64()
	var buckName string

	if trans.FromSource {

		if trans.CurrencySource == CurrencyUBuck {
			buckName = "UBuck"
		} else if trans.CurrencySource == KeyDebt {
			buckName = KeyDebt
		} else if strings.Contains(trans.CurrencySource, "@") {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencySource)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"

		} else {
			buckName = trans.CurrencySource
		}
	} else {

		if trans.CurrencyDest == CurrencyUBuck {
			buckName = "UBuck"
		} else if trans.CurrencyDest == KeyDebt {
			buckName = KeyDebt
		} else if strings.Contains(trans.CurrencyDest, "@") {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencyDest)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"
		} else {
			buckName = trans.CurrencyDest
		}

	}

	resp = openapi.ResponseCryptoTransaction{
		Balance:         float32(balance),
		Description:     trans.Reference,
		ConversionRatio: float32(xrate),
		Amount:          float32(amount),
		Name:            buckName,
		CreatedAt:       trans.Ts,
		UBucks:          float32(trans.Ubuck.InexactFloat64()),
	}

	return
}

func cryptoConvert(basis, usd, amount decimal.Decimal) decimal.Decimal {

	gains := usd.Sub(basis).IsPositive()
	var newPercent decimal.Decimal
	if gains {
		percentChange := usd.Div(basis).Sub(decimal.NewFromInt32(1))
		adjustedChange := percentChange.Mul(decimal.NewFromInt32(3))
		newPercent = adjustedChange.Add(decimal.NewFromInt32(1))
	} else {
		percentChange := usd.Div(basis)
		newPercent = percentChange.Pow(decimal.NewFromInt32(3))
	}

	return newPercent.Mul(basis).Mul(amount)
}

func purchaseLotto(db *bolt.DB, clock Clock, studentDetails UserInfo, tickets int32) (winner bool, err error) {

	err = db.Update(func(tx *bolt.Tx) error {

		lottery, err := getLottoLatestTx(tx, studentDetails)
		if err != nil {
			return err
		}

		if lottery.Winner != "" {
			return fmt.Errorf("last winner " + lottery.Winner + " with " + strconv.Itoa(int(lottery.Jackpot)) + "...lottery has been disabled")
		}

		if lottery.Odds == 0 && lottery.Jackpot == 0 {
			return fmt.Errorf("the lotto has not been initialized")
		}

		chargeStudentTx(tx, clock, studentDetails, decimal.NewFromInt32(tickets).Mul(decimal.NewFromInt32(KeyPricePerTicket)), CurrencyUBuck, "Lotto", true)
		if err != nil {
			return err
		}

		userUpdates := openapi.RequestUserEdit{
			LottoPlay: tickets * KeyPricePerTicket,
		}

		lgr.Printf(studentDetails.Email + " purchased " + strconv.Itoa(int(tickets)))

		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		for i := 0; i < int(tickets); i++ {
			play := r.Intn(int(lottery.Odds))
			if play == int(lottery.Number) {

				err = updateLottoLatestTx(tx, studentDetails, tickets, studentDetails.Email)
				if err != nil {
					return err
				}

				err = pay2StudentTx(tx, clock, studentDetails, decimal.NewFromInt32(lottery.Jackpot+tickets), CurrencyUBuck, "Lotto Winner "+clock.Now().Format("02/03/2006"))
				if err != nil {
					return err
				}

				lgr.Printf(studentDetails.Email + " wins " + strconv.Itoa(int(lottery.Jackpot+tickets)))

				userUpdates.LottoWin = lottery.Jackpot + tickets

				err = userEditTx(tx, clock, studentDetails, userUpdates)
				if err != nil {
					return err
				}

				settings, err := getSettingsRx(tx, studentDetails)
				if err != nil {
					return err
				}

				if settings.Lottery {
					err = initializeLotteryTx(tx, studentDetails, settings, clock)
					if err != nil {
						return err
					}
				}

				winner = true

				return nil

			}
		}

		err = userEditTx(tx, clock, studentDetails, userUpdates)
		if err != nil {
			return err
		}

		err = updateLottoLatestTx(tx, studentDetails, tickets, "")
		if err != nil {
			return err
		}

		return nil
	})

	return

}

func timeToInterest(time int32) float32 {
	if time <= 14 {
		return 1.03
	}
	if time <= 30 {
		return 1.04
	}
	if time <= 50 {
		return 1.05
	}
	if time <= 70 {
		return 1.06
	}
	return 1.07
}

func buyCD(db *bolt.DB, clock Clock, userDetails UserInfo, body openapi.RequestBuyCd) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return buyCDTx(tx, clock, userDetails, body)
	})

}

func buyCDTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, body openapi.RequestBuyCd) (err error) {

	prinInv := decimal.NewFromInt32(body.PrinInv)

	if prinInv.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)
	mature := ts.Add(time.Hour * 24 * time.Duration(body.Time))

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Email,
		Destination:    userInfo.Email,
		CurrencySource: CurrencyUBuck,
		CurrencyDest:   KeyCertificateOfDeposit,
		AmountSource:   prinInv,
		AmountDest:     prinInv,
		XRate:          decimal.NewFromInt32(1),
		Reference:      "Ubuck to " + strconv.FormatInt(int64(body.Time), 10) + " day CD",
		FromSource:     true,
	}

	CD := CertificateOfDeposit{
		Ts:           ts,
		Principal:    body.PrinInv,
		CurrentValue: prinInv,
		RefundValue:  prinInv.Mul(decimal.NewFromFloat32(.9)),
		Interest:     timeToInterest(body.Time),
		Maturity:     mature,
		Active:       true,
	}

	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(student, CurrencyUBuck, transaction, OperationDebit, true)
	if err != nil {
		return err
	}

	transaction.FromSource = false

	err = addToCDHolderTx(student, transaction, CD, "")
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	err = addToCDHolderTx(cb, transaction, CD, "")
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func getCDS(db *bolt.DB, userDetails UserInfo) (resp []openapi.ResponseCd, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getCDSRx(tx, userDetails)
		return err
	})

	return

}

func getCDSRx(tx *bolt.Tx, userInfo UserInfo) (resp []openapi.ResponseCd, err error) {

	student, err := getStudentBucketRx(tx, userInfo.Name)
	if err != nil {
		return
	}

	accountsBucket := student.Bucket([]byte(KeyAccounts))
	if accountsBucket == nil {
		return resp, fmt.Errorf("cannot find accounts bucket")
	}

	CDS_bucket := accountsBucket.Bucket([]byte(KeyCertificateOfDeposit))
	if CDS_bucket == nil {
		return resp, fmt.Errorf("cannot find CDS bucket")
	}

	c := CDS_bucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			continue
		}

		CD_data := CDS_bucket.Get(k)

		var deposit CertificateOfDeposit
		err = json.Unmarshal(CD_data, &deposit)
		if err != nil {
			return resp, err
		}

		if deposit.Active {
			item := openapi.ResponseCd{
				Ts:           deposit.Ts,
				Principal:    deposit.Principal,
				CurrentValue: float32(deposit.CurrentValue.InexactFloat64()),
				Interest:     deposit.Interest,
				Maturity:     deposit.Maturity,
				RefundValue:  float32(deposit.RefundValue.InexactFloat64()),
			}
			resp = append(resp, item)
		}

	}

	return
}

func getCDTransactions(db *bolt.DB, userDetails UserInfo) (resp []openapi.ResponseTransactions, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getCDTransactionsRx(tx, userDetails)
		return err
	})

	return

}

func getCDTransactionsRx(tx *bolt.Tx, userInfo UserInfo) (resp []openapi.ResponseTransactions, err error) {

	student, err := getStudentBucketRx(tx, userInfo.Name)
	if err != nil {
		return
	}

	accountsBucket := student.Bucket([]byte(KeyAccounts))
	if accountsBucket == nil {
		return resp, fmt.Errorf("cannot find accounts bucket")
	}

	CDS_bucket := accountsBucket.Bucket([]byte(KeyCertificateOfDeposit))
	if CDS_bucket == nil {
		return resp, fmt.Errorf("cannot find CDS bucket")
	}

	transactionsBucket := CDS_bucket.Bucket([]byte(KeyTransactions))
	if transactionsBucket == nil {
		return resp, fmt.Errorf("cannot find transactions bucket")
	}

	c := transactionsBucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			continue
		}

		var trans Transaction
		err = json.Unmarshal(v, &trans)
		if err != nil {
			return resp, err
		}

		item := openapi.ResponseTransactions{
			CreatedAt:   trans.Ts,
			Amount:      float32(trans.AmountSource.InexactFloat64()),
			Description: trans.Reference,
		}
		resp = append(resp, item)

		if len(resp) > 24 {
			break
		}
	}

	return
}

func InterestToTime(interest float32) int64 {
	if interest == 1.03 {
		return 14
	}
	if interest == 1.04 {
		return 30
	}
	if interest == 1.05 {
		return 50
	}
	if interest == 1.06 {
		return 70
	}
	return 90
}

func matureCheck(CD CertificateOfDeposit) string {
	if CD.CurrentValue.Equal(CD.RefundValue) {
		return "Fully Matured"
	}
	return "Early Refund"
}

func refundCD(db *bolt.DB, clock Clock, userDetails UserInfo, CD_id string) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return refundCDTx(tx, clock, userDetails, CD_id)
	})

}

func refundCDTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, CD_id string) (err error) {

	student, err := getStudentBucketRx(tx, userInfo.Name)
	if err != nil {
		return
	}

	accountsBucket := student.Bucket([]byte(KeyAccounts))
	if accountsBucket == nil {
		return fmt.Errorf("cannot find accounts bucket")
	}

	CDS_bucket := accountsBucket.Bucket([]byte(KeyCertificateOfDeposit))
	if CDS_bucket == nil {
		return fmt.Errorf("cannot find CDS bucket")
	}

	CD_data := CDS_bucket.Get([]byte(CD_id))
	if CD_data == nil {
		return fmt.Errorf("cannot find the CD")
	}

	var CD CertificateOfDeposit
	err = json.Unmarshal(CD_data, &CD)
	if err != nil {
		return err
	}

	if !CD.Active {
		return nil
	}

	ts := clock.Now().Truncate(time.Millisecond)

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Email,
		Destination:    userInfo.Email,
		CurrencySource: KeyCertificateOfDeposit,
		CurrencyDest:   CurrencyUBuck,
		AmountSource:   CD.RefundValue,
		AmountDest:     CD.RefundValue,
		XRate:          decimal.NewFromInt32(1),
		Reference:      strconv.FormatInt(InterestToTime(CD.Interest), 10) + " day CD to Ubuck: " + matureCheck(CD),
		FromSource:     true,
	}

	CD.Active = false

	_, _, err = addToHolderTx(student, CurrencyUBuck, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	transaction.FromSource = false

	err = addToCDHolderTx(student, transaction, CD, CD_id)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	err = addToCDHolderTx(cb, transaction, CD, CD_id)
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationCredit, false)
	if err != nil {
		return err
	}

	return nil
}

func CDSumTx(tx *bolt.Tx, account *bolt.Bucket) (sum decimal.Decimal) {
	sum = decimal.Zero
	c := account.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v == nil {
			continue
		}

		var deposit CertificateOfDeposit
		_ = json.Unmarshal(v, &deposit)

		if deposit.Active {
			sum = sum.Add(deposit.RefundValue)
		}

	}

	return
}
