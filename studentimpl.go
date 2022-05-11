package main

import (
	"encoding/json"
	"fmt"
	"time"

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
}

// func getClassbyAddCodeTx(tx *bolt.Tx, schoolId, addCode string) (classBucket *bolt.Bucket, err error) {
// 	schools := tx.Bucket([]byte(KeySchools))
// 	if schools == nil {
// 		return nil, fmt.Errorf("schools not found")
// 	}
// 	school := schools.Bucket([]byte(schoolId))
// 	if school == nil {
// 		return nil, fmt.Errorf("school not found")
// 	}
// 	schoolClass :=
// }

// adds ubucks from CB
//creates order, transaction into student account, update account balance, update ubuck balance
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

// addToHolderTx updates balance and adds transaction
// debit means to remove money
func addToHolderTx(holder *bolt.Bucket, account string, transaction Transaction, direction int) (balance decimal.Decimal, errR error) {
	accounts, err := holder.CreateBucketIfNotExists([]byte(KeybAccounts))
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

	if direction == OperationCredit {
		balance = balance.Add(transaction.AmountDest)
	} else {
		balance = balance.Sub(transaction.AmountSource)
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
	student, err := getStudentBucketRoTx(tx, userName)
	if err != nil {
		return
	}
	accounts := student.Bucket([]byte(KeybAccounts))
	if accounts == nil {
		return
	}
	ubuck := accounts.Bucket([]byte(CurrencyUBuck))
	if ubuck == nil {
		return
	}
	ubuckB := ubuck.Get([]byte(KeyBalance))
	if ubuckB == nil {
		return
	}
	err = res.UnmarshalText(ubuckB)
	if err != nil {
		lgr.Printf("ERROR cannot unmarshal balance for ubuck of %s: %v", userName, err)
		return decimal.Zero
	}
	return
}

func DailyPayIfNeeded(db *bolt.DB, clock Clock, userDetails UserInfo) bool {
	if userDetails.Role != UserRoleStudent {
		return false
	}

	needToAdd := false
	_ = db.View(func(tx *bolt.Tx) error {
		student, _ := getStudentBucketRoTx(tx, userDetails.Name)
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

		payDate, err := clock.Now().Truncate(24 * time.Hour).MarshalText()
		if err != nil {
			return err
		}
		err = student.Put([]byte(KeyDayPayment), payDate)
		if err != nil {
			return fmt.Errorf("cannot save daily payment date: %v", err)
		}
<<<<<<< HEAD
		pay := decimal.NewFromFloat32(float32(userDetails.Income))
		return addUbuck2StudentTx(tx, clock, userDetails.Name, pay, "daily payment")
=======
		pay := decimal.NewFromFloat32(121.32)
		return addUbuck2StudentTx(tx, clock, userDetails, pay, "daily payment")
>>>>>>> origin/main
	})

	if err != nil {
		lgr.Printf("ERROR daily payment not added to %s: %v", userDetails.Name, err)
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

func chargeStudentUbuck(db *bolt.DB, clock Clock, userDetails UserInfo, amount decimal.Decimal, reference string) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return chargeStudentUbuckTx(tx, clock, userDetails, amount, reference)
	})
}
func chargeStudentUbuckTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, amount decimal.Decimal, reference string) (err error) {
	return chargeStudentTx(tx, clock, userDetails, amount, CurrencyUBuck, reference)

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
	
	ts := clock.Now()

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
	_, err = addToHolderTx(student, currency, transaction, OperationCredit)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}
	_, err = addToHolderTx(cb, currency, transaction, OperationDebit)
	if err != nil {
		return err
	}

	return nil
}

func chargeStudent(db *bolt.DB, clock Clock, userDetails UserInfo, amount decimal.Decimal, currency string, reference string) (err error) {
	return db.Update(func(tx *bolt.Tx) error {
		return chargeStudentTx(tx, clock, userDetails, amount, currency, reference)
	})
}

func chargeStudentTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, amount decimal.Decimal, currency string, reference string) (err error) {
	if userDetails.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now()

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
		return err
	}

	newBalance, err := addToHolderTx(student, currency, transaction, OperationDebit)
	if err != nil {
		return err
	}
	if newBalance.Sign() < 0 {
		return fmt.Errorf("ubuck balance for %s is negative", userDetails.Name)
	}

	cb, err := getCbTx(tx, userDetails.SchoolId)
	if err != nil {
		return err
	}
	_, err = addToHolderTx(cb, currency, transaction, OperationCredit)
	if err != nil {
		return err
	}

	return nil
}

/*
check if the currency exists in given school
*/
func isCurrencyTx(tx *bolt.Tx, schoolId string, currency string) (bool, error) {
	if currency == CurrencyUBuck {
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