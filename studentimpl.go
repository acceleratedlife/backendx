package main

import (
	"encoding/json"
	"fmt"
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
func addToHolderTx(holder *bolt.Bucket, account string, transaction Transaction, direction int, negBlock bool) (balance decimal.Decimal, errR error) {
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

	if direction == OperationCredit {
		balance = balance.Add(transaction.AmountDest)
	} else {
		balance = balance.Sub(transaction.AmountSource)
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
	student, err := getStudentBucketRoTx(tx, userName)
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

		ubuck, _, err := convertRx(tx, userData.SchoolId, string(k), "", value.InexactFloat64())
		if err != nil {
			return
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
		pay := decimal.NewFromFloat32(userDetails.Income)
		return addUbuck2StudentTx(tx, clock, userDetails, pay, "daily payment")
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
	_, err = addToHolderTx(student, currency, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}
	_, err = addToHolderTx(cb, currency, transaction, OperationDebit, false)
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

func studentConvertTx(tx *bolt.Tx, clock Clock, userInfo UserInfo, amount decimal.Decimal, from string, to string, charge bool) (err error) {
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

	if from != CurrencyUBuck {
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

	transaction := Transaction{
		Ts:             ts,
		Source:         userInfo.Name,
		Destination:    userInfo.Name,
		CurrencySource: from,
		CurrencyDest:   to,
		AmountSource:   amount,
		AmountDest:     converted,
		XRate:          xRate,
		Reference:      fromDetails.LastName + " to " + toDetails.LastName,
		FromSource:     true,
	}

	student, err := getStudentBucketTx(tx, userInfo.Name)
	if err != nil {
		return err
	}
	if to != KeyDebt || charge {
		_, err = addToHolderTx(student, from, transaction, OperationDebit, true)
		if err != nil {
			return err
		}
	}

	transaction.FromSource = false

	if to == KeyDebt && charge {
		transaction.AmountDest = transaction.AmountDest.Neg()
	}

	_, err = addToHolderTx(student, to, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, userInfo.SchoolId)
	if err != nil {
		return err
	}

	_, err = addToHolderTx(cb, to, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	transaction.FromSource = true

	_, err = addToHolderTx(cb, from, transaction, OperationDebit, false)
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
		return err
	}

	_, err = addToHolderTx(student, currency, transaction, OperationDebit, true)
	if err != nil {
		if err.Error() == "Insufficient funds" {
			err := studentConvertTx(tx, clock, userDetails, amount, currency, KeyDebt, false)
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
	_, err = addToHolderTx(cb, currency, transaction, OperationCredit, false)
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
		uBucks, err = getStudentUbuckTx(tx, userDetails)
		if err != nil {
			return err
		}
		return nil
	})

	return
}

func getStudentUbuckTx(tx *bolt.Tx, userDetails UserInfo) (resp openapi.ResponseSearchStudentUbuck, err error) {
	student, err := getStudentBucketRoTx(tx, userDetails.Name)
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

func getStudentAuctionsTx(tx *bolt.Tx, auctionsBucket *bolt.Bucket, userDetails UserInfo) (auctions []openapi.Auction, err error) {

	c := auctionsBucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}

		auctionBucket := auctionsBucket.Bucket(k)
		if auctionBucket == nil {
			return auctions, fmt.Errorf("cannot find auction bucket")
		}

		visibilityBucket := auctionBucket.Bucket([]byte(KeyVisibility))
		if visibilityBucket == nil {
			return auctions, fmt.Errorf("cannot find visibility bucket")
		}

		classes, err := getStudentClassesTx(tx, userDetails)
		if err != nil {
			return auctions, fmt.Errorf("cannot get student classes %s: %v", userDetails.Name, err)
		}

		x := visibilityBucket.Cursor()
		for k, _ := x.First(); k != nil; k, _ = x.Next() {
			for _, class := range classes {
				if class.Id == string(k) || string(k) == KeyEntireSchool {

					iAuction := openapi.Auction{
						Id:          string(k),
						StartDate:   string(auctionBucket.Get([]byte(KeyStartDate))),
						EndDate:     string(auctionBucket.Get([]byte(KeyEndDate))),
						Bid:         btoi32(auctionBucket.Get([]byte(KeyBid))),
						MaxBid:      btoi32(auctionBucket.Get([]byte(KeyMaxBid))),
						Description: string(auctionBucket.Get([]byte(KeyDescription))),
					}

					ownerDetails, err := getUserInLocalStoreTx(tx, string(auctionBucket.Get([]byte(KeyOwnerId))))
					if err != nil {
						return nil, err
					}

					iAuction.OwnerId = openapi.AuctionOwnerId{
						LastName: ownerDetails.LastName,
						Id:       ownerDetails.Name,
					}

					winnerDetails, err := getUserInLocalStoreTx(tx, string(auctionBucket.Get([]byte(KeyWinnerId))))
					if err != nil {
						iAuction.WinnerId = openapi.AuctionWinnerId{
							FirstName: "nil",
							LastName:  "nil",
							Id:        "nil",
						}
					} else {
						iAuction.WinnerId = openapi.AuctionWinnerId{
							FirstName: winnerDetails.FirstName,
							LastName:  winnerDetails.LastName,
							Id:        winnerDetails.Name,
						}
					}

					auctions = append(auctions, iAuction)
					break

				}
			}
		}

	}

	return auctions, nil
}

func getStudentClassesTx(tx *bolt.Tx, userDetails UserInfo) (classes []openapi.Class, err error) {
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

		buck := bAccounts.Bucket(k)

		transactions := buck.Bucket([]byte(KeyTransactions))
		if transactions == nil {
			return resp, fmt.Errorf("Cannot get transactions")
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
		} else {
			trans.Net = trans.AmountSource.Neg()
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
		} else {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencySource)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"

		}
	} else {

		if trans.CurrencyDest == CurrencyUBuck {
			buckName = "UBuck"
		} else if trans.CurrencyDest == KeyDebt {
			buckName = KeyDebt
		} else {
			user, err := getUserInLocalStoreTx(tx, trans.CurrencyDest)
			if err != nil {
				return resp, err
			}

			buckName = user.LastName + " Buck"
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
