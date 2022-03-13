package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
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

type StudentApiServiceImpl struct {
	db *bolt.DB
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
func addUbuck2Student(db *bolt.DB, clock Clock, studentId string, amount decimal.Decimal, reference string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return addUbuck2StudentTx(tx, clock, studentId, amount, reference)
	})
}

// register order, transactions in students account, transactions i CB
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

	transaction := Transaction{
		Ts:             ts,
		Source:         "",
		Destination:    studentId,
		CurrencySource: CurrencyUBuck,
		CurrencyDest:   CurrencyUBuck,
		AmountSource:   amount,
		AmountDest:     amount,
		XRate:          decimal.NewFromFloat32(1.0),
		Reference:      reference,
	}
	student, err := getStudentBucketTx(tx, studentId)
	if err != nil {
		return err
	}
	_, err = addToHolderTx(student, CurrencyUBuck, transaction, OperationCredit)
	if err != nil {
		return err
	}

	cb, err := tx.CreateBucketIfNotExists([]byte(KeyCB))
	if err != nil {
		return err
	}
	_, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationDebit)
	if err != nil {
		return err
	}

	return nil

}

// addToHolderTx updates balance and adds transaction
// debit means to remove money
func addToHolderTx(holder *bolt.Bucket, account string, transaction Transaction, direction int) (balance decimal.Decimal, errR error) {
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

		payDate, err := clock.Now().Truncate(24 * time.Hour).MarshalText()
		if err != nil {
			return err
		}
		err = student.Put([]byte(KeyDayPayment), payDate)
		if err != nil {
			return fmt.Errorf("cannot save gaily payment date: %v", err)
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

func (a *StudentApiServiceImpl) StudentAddClass(ctx context.Context, body openapi.RequestAddClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role != UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	_, pathId, err := RoleByAddCode(a.db, body.AddCode)
	if err != nil {
		return openapi.Response(404,
			openapi.ResponseRegister4{
				Message: err.Error(),
			}), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		// classBucket, err := getClassbyAddCodeTx(tx, userDetails.SchoolId, body.AddCode)
		school := tx.Bucket([]byte(userDetails.SchoolId))
		if school == nil {
			return fmt.Errorf("can't find school")
		}

		schoolClasses := school.Bucket([]byte(KeyClasses))
		if schoolClasses == nil {
			return fmt.Errorf("can't find school classes")
		}
		class := schoolClasses.Bucket([]byte(pathId.classId))
		if class != nil {
			studentsBucket := class.Bucket([]byte(KeyStudents))
			if studentsBucket != nil {
				err = studentsBucket.Put([]byte(userDetails.Email), nil)
				if err != nil {
					return err
				}
				return nil
			}
		}

		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			return fmt.Errorf("can't find teachers")
		}
		teacher := teachers.Bucket([]byte(pathId.teacherId))
		if teacher == nil {
			return fmt.Errorf("can't find teacher")
		}
		classes := teacher.Bucket([]byte(KeyClasses))
		if classes == nil {
			return fmt.Errorf("can't find class")
		}
		classBucket := classes.Bucket([]byte(pathId.classId))
		if classBucket == nil {
			return fmt.Errorf("can't find class")
		}
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		if studentsBucket == nil {
			return fmt.Errorf("can't find students")
		}
		err = studentsBucket.Put([]byte(userDetails.Email), nil)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		lgr.Printf("ERROR cannot edit class with Id: %s %v", body.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	resp, err := classesWithOwnerDetails(a.db, userDetails.SchoolId, userDetails.Email)
	if err != nil {
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil

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
