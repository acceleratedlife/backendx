package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
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

	if balance.Sign() < 0 && negBlock {
		errR = fmt.Errorf("Insufficient funds")
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
	if oldTransaction != nil {
		transaction.Ts.Add(time.Millisecond * 1)
		tsB, err = transaction.Ts.MarshalText()
		if err != nil {
			errR = err
			return
		}
	}

	transactionB, err := json.Marshal(transaction)
	if err != nil {
		errR = err
		return
	}

	errR = transactions.Put(tsB, transactionB)
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
			garnish := pay.Mul(decimal.NewFromFloat32(.3))
			pay = pay.Mul(decimal.NewFromFloat32(.7))
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
		random := decimal.NewFromFloat32(rand.Float32())

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

		if student.College && student.CollegeEnd.IsZero() {
			student.Job = getJobIdRx(tx, KeyCollegeJobs)
			jobDetails := getJobRx(tx, KeyCollegeJobs, student.Job)
			max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
			min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
			diff := max.Sub(min)
			random := decimal.NewFromFloat32(rand.Float32())
			student.Income = float32(random.Mul(diff).Add(min).Floor().InexactFloat64())
		} else {
			student.Job = getJobIdRx(tx, KeyJobs)
			jobDetails := getJobRx(tx, KeyJobs, student.Job)
			max := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(192))
			min := decimal.NewFromInt32(jobDetails.Pay).Div(decimal.NewFromInt32(250))
			diff := max.Sub(min)
			random := decimal.NewFromFloat32(rand.Float32())
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
		interest := decimal.NewFromFloat32(1.06)
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

		days := rand.Intn(5) + 4

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

	var event eventRequest
	err := json.Unmarshal(events.Get([]byte(idKey)), &event)
	if err != nil {
		return ""
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
	pick := rand.Intn(bucketStats.KeyN)
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
	multiplier := decimal.NewFromFloat(.3)
	one := decimal.NewFromInt(1)
	random := decimal.NewFromFloat32(rand.Float32())
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

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(cb, currency, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	if currency != CurrencyUBuck && currency != KeyDebt {
		_, err = addStepTx(tx, userInfo.SchoolId, currency, float32(amount.InexactFloat64()))
		if err != nil {
			return err
		}
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
	}

	if from != CurrencyUBuck && from != KeyDebt {
		fromDetails, err = getUserInLocalStoreTx(tx, from)
		if err != nil {
			return err
		}
	}

	converted, xRate, err := convertRx(tx, userInfo.SchoolId, from, target, amount.InexactFloat64())
	if err != nil {
		return err
	}

	if charge {
		amount = amount.Mul(decimal.NewFromFloat32(keyCharge))
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
		if err.Error() == "Insufficient funds" && !sPurchase {
			if strings.Contains(reference, "Event: ") {
				err := studentConvertTx(tx, clock, userDetails, amount, currency, KeyDebt, reference, false)
				if err != nil {
					return err
				}
			}
			err := studentConvertTx(tx, clock, userDetails, amount, currency, KeyDebt, "", false)
			if err != nil {
				return err
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

	for k, v := c.First(); k != nil; k, v = c.Next() {
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
	school, err := getSchoolBucketTx(tx, userDetails)
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

	if strings.Contains(resp.Description, "Event:") {
		typeKey := KeyPEvents
		if resp.Amount <= 0 {
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
	// crypto = strings.ToLower(crypto)
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
		if string(k) == CurrencyUBuck || string(k) == KeyDebt || strings.Contains(string(k), "@") {
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
