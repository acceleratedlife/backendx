package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
	"time"
)

type Order struct {
	Source      string // source student/cb
	Destination string // dest student
	Currency    string
	Amount      decimal.Decimal
	Reference   string // reason
}

type Transaction struct {
	Source         string
	Destination    string
	CurrencySource string
	CurrencyDest   string
	Amount         decimal.Decimal
	XRate          decimal.Decimal
	Reference      string
}

// creates order, transaction into student account, update account balance, update ubuck balance
func addUbuck2Student(db *bolt.DB, clock AppClock, studentId string, amount decimal.Decimal, reference string) error {
	return db.Update(func(tx *bolt.Tx) error {
		ts := clock.Now()
		tsk, err := ts.MarshalText()

		if err != nil {
			return err
		}

		order := Order{
			Source:      KeyCB,
			Destination: studentId,
			Currency:    CurrencyUBuck,
			Amount:      amount,
			Reference:   reference,
		}

		orders, err := tx.CreateBucketIfNotExists([]byte(KeyOrders))
		if err != nil {
			return fmt.Errorf("no bucket orders: %v", err)
		}

		orderV, err := json.Marshal(order)
		if err != nil {
			return err
		}

		err = orders.Put(tsk, orderV)
		if err != nil {
			return err
		}

		cb, err := tx.CreateBucketIfNotExists([]byte(KeyCB))
		if err != nil {
			return err
		}

		err = addBalanceTx(cb, CurrencyUBuck, amount)
		if err != nil {
			return err
		}

		user, err := getStudentBucketTx(tx, studentId)
		if err != nil {
			return err
		}

		accounts, err := user.CreateBucketIfNotExists([]byte(KeyAccounts))
		if err != nil {
			return err
		}

		ubuckBucket, err := accounts.CreateBucketIfNotExists([]byte(CurrencyUBuck))
		if err != nil {
			return err
		}
		err = addBalanceTx(ubuckBucket, KeyBalance, amount)
		if err != nil {
			return err
		}

		transactions, err := ubuckBucket.CreateBucketIfNotExists([]byte(KeyTransactions))
		if err != nil {
			return err
		}

		transaction := Transaction{
			Source:         KeyCB,
			Destination:    studentId,
			CurrencySource: CurrencyUBuck,
			CurrencyDest:   CurrencyUBuck,
			Amount:         amount,
			XRate:          decimal.NewFromInt(1),
			Reference:      reference,
		}

		transactionB, err := json.Marshal(transaction)
		if err != nil {
			return err
		}
		return transactions.Put(tsk, transactionB)

	})
}

func addBalanceTx(bucket *bolt.Bucket, key string, amount decimal.Decimal) error {
	var ubucks decimal.Decimal
	ubucksB := bucket.Get([]byte(key))
	if ubucksB != nil {
		err := ubucks.UnmarshalText(ubucksB)
		if err != nil {
			return err
		}
	}

	ubucks = ubucks.Add(amount)
	ubucksB, err := ubucks.MarshalText()
	if err != nil {
		return err
	}
	err = bucket.Put([]byte(key), ubucksB)
	if err != nil {
		return err
	}
	return nil
}

func StudentNetWorthTx(tx *bolt.Tx, userName string) (res decimal.Decimal) {
	res = decimal.Zero
	student, err := getStudentBucketTx(tx, userName)
	if err != nil {
		return
	}
	accounts := student.Bucket([]byte(KeyAccounts))
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

		pay := decimal.NewFromFloat32(121.32)
		return addUbuck2StudentTx(tx, clock, userDetails.Name, pay, "daily payment")
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

func addUbuck2StudentTx(tx *bolt.Tx, clock Clock, studentId string, amount decimal.Decimal, reference string) error {
	ts := clock.Now()
	tsk, err := ts.MarshalText()

	if err != nil {
		return err
	}

	order := Order{
		Source:      KeyCB,
		Destination: studentId,
		Currency:    CurrencyUBuck,
		Amount:      amount,
		Reference:   reference,
	}

	orders, err := tx.CreateBucketIfNotExists([]byte(KeyOrders))
	if err != nil {
		return fmt.Errorf("no bucket orders: %v", err)
	}

	orderV, err := json.Marshal(order)
	if err != nil {
		return err
	}

	err = orders.Put(tsk, orderV)
	if err != nil {
		return err
	}

	cb, err := tx.CreateBucketIfNotExists([]byte(KeyCB))
	if err != nil {
		return err
	}

	err = addBalanceTx(cb, CurrencyUBuck, amount)
	if err != nil {
		return err
	}

	user, err := getStudentBucketTx(tx, studentId)
	if err != nil {
		return err
	}

	accounts, err := user.CreateBucketIfNotExists([]byte(KeyAccounts))
	if err != nil {
		return err
	}

	ubuckBucket, err := accounts.CreateBucketIfNotExists([]byte(CurrencyUBuck))
	if err != nil {
		return err
	}
	err = addBalanceTx(ubuckBucket, KeyBalance, amount)
	if err != nil {
		return err
	}

	transactions, err := ubuckBucket.CreateBucketIfNotExists([]byte(KeyTransactions))
	if err != nil {
		return err
	}

	transaction := Transaction{
		Source:         KeyCB,
		Destination:    studentId,
		CurrencySource: CurrencyUBuck,
		CurrencyDest:   CurrencyUBuck,
		Amount:         amount,
		XRate:          decimal.NewFromInt(1),
		Reference:      reference,
	}

	transactionB, err := json.Marshal(transaction)
	if err != nil {
		return err
	}
	return transactions.Put(tsk, transactionB)
}

//How to calculate netWorth
//givens:
//student accounts:
//uBucks 0, Kirill Bucks 5
//
//uBuck total currency: 1000
//Kirill Bucks total currency: 100
//conversion ratio 1000/100 = 10
//10 ubucks = 1 kirill buck
//networth = uBucks + Kirill bucks *10
//50 = 0 + (5*10)
//networth = 50
