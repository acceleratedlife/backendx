package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

func getClassAtSchoolTx(tx *bolt.Tx, schoolId, classId string) (classBucket *bolt.Bucket, parentBucket *bolt.Bucket, err error) {

	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return nil, nil, err
	}

	classes := school.Bucket([]byte(KeyClasses))
	if classes != nil {
		classBucket := classes.Bucket([]byte(classId))
		if classBucket != nil {
			return classBucket, classes, nil
		}
	}

	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return nil, nil, fmt.Errorf("no teachers at school")
	}
	cTeachers := teachers.Cursor()
	for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
		if v != nil {
			continue
		}
		teacher := teachers.Bucket(k)
		if teacher == nil {
			continue
		}
		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			continue
		}
		classBucket = classesBucket.Bucket([]byte(classId)) //found the class
		if classBucket == nil {
			continue
		}
		return classBucket, classesBucket, nil
	}
	return nil, nil, fmt.Errorf("class not found")
}

func PopulateClassMembers(tx *bolt.Tx, classBucket *bolt.Bucket) (Members []openapi.ClassWithMembersMembers, err error) {
	Members = make([]openapi.ClassWithMembersMembers, 0)
	students := classBucket.Bucket([]byte(KeyStudents))
	if students == nil {
		return Members, nil
	}
	cStudents := students.Cursor()
	for k, _ := cStudents.First(); k != nil; k, _ = cStudents.Next() { //iterate students bucket
		user, err := getUserInLocalStoreTx(tx, string(k))
		if err != nil {
			return nil, err
		}

		nUser := openapi.ClassWithMembersMembers{
			Id:        user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Rank:      user.Rank,
			NetWorth:  user.NetWorth,
		}
		Members = append(Members, nUser)
	}
	return Members, nil
}

func getStudentHistory(db *bolt.DB, userName string, schoolId string) (history []openapi.History, err error) {
	_ = db.View(func(tx *bolt.Tx) error {
		history, err = getStudentHistoryTX(tx, userName)
		return nil
	})
	return
}
func getStudentHistoryTX(tx *bolt.Tx, userName string) (history []openapi.History, err error) {
	student, err := getStudentBucketRx(tx, userName)
	if err == nil {
		return nil, err
	}

	historyData := student.Get([]byte(KeyHistory))
	if historyData == nil {
		return nil, fmt.Errorf("Failed to get history")
	}
	err = json.Unmarshal(historyData, &history)
	if err != nil {
		return nil, fmt.Errorf("ERROR cannot unmarshal History")
	}
	return
}

func getStudentAccount(db *bolt.DB, bAccount *bolt.Bucket, accountId string) (resp openapi.ResponseCurrencyExchange, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getStudentAccountRx(tx, bAccount, accountId)
		return err
	})

	return
}

func getStudentAccountRx(tx *bolt.Tx, bAccount *bolt.Bucket, accountId string) (resp openapi.ResponseCurrencyExchange, err error) {

	balanceData := bAccount.Get([]byte(KeyBalance))
	err = json.Unmarshal(balanceData, &resp.Balance)
	resp.Id = accountId

	return
}

func getCBaccountDetailsRx(tx *bolt.Tx, userDetails UserInfo, account openapi.ResponseCurrencyExchange) (finalAccount openapi.ResponseCurrencyExchange, err error) {
	cb, err := getCbRx(tx, userDetails.SchoolId)
	if err != nil {
		return
	}

	accounts := cb.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return finalAccount, fmt.Errorf("cannot find cb buck accounts")
	}

	bAccount := accounts.Bucket([]byte(account.Id))
	if bAccount == nil {
		return finalAccount, fmt.Errorf("cannot find cb buck account")
	}

	finalAccount = account

	return
}

func getBuckNameTx(tx *bolt.Tx, id string) (string, error) {
	if id == KeyDebt {
		return "Debt", nil
	}

	if id == CurrencyUBuck {
		return "UBuck", nil
	}

	user, err := getUserInLocalStoreTx(tx, id)
	if err != nil {
		return "", err
	}

	return user.LastName + " Buck", nil
}

func saveRanksTx(tx *bolt.Tx, students []openapi.UserNoHistory) (ranked int, err error) {
	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		return ranked, fmt.Errorf("users not found")
	}

	length := float32(len(students))
	for i, student := range students {
		user := users.Get([]byte(student.Id))
		if user == nil {
			return ranked, fmt.Errorf("user not found")
		}

		var userDetails UserInfo
		err = json.Unmarshal(user, &userDetails)
		if err != nil {
			return ranked, err
		}

		num := float32(i+1) / length
		userDetails.Rank = 0

		if length <= 60 {
			if num <= .33334 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length <= 150 {
			if num <= .2 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length <= 500 {
			if num <= .15 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length <= 1000 {
			if num <= .1 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length <= 2000 {
			if num <= .0625 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length <= 4000 {
			if num <= .05 {
				userDetails.Rank = student.Rank
				ranked++
			}
		} else if length > 4000 {
			if i < 200 {
				userDetails.Rank = student.Rank
				ranked++
			}
		}

		userDetails.NetWorth = student.NetWorth

		marshal, err := json.Marshal(userDetails)
		if err != nil {
			return ranked, err
		}

		err = users.Put([]byte(student.Id), marshal)
		if err != nil {
			return ranked, err
		}

	}

	return ranked, err

}

func saveRanks(db *bolt.DB, students []openapi.UserNoHistory) (ranked int, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		ranked, err = saveRanksTx(tx, students)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func getSchoolStudentsTx(tx *bolt.Tx, userDetails UserInfo) (resp []openapi.UserNoHistory, ranked int, err error) {
	school, err := SchoolByIdTx(tx, userDetails.SchoolId)
	if err != nil {
		return
	}

	students := school.Bucket([]byte(KeyStudents))
	if students == nil {
		return resp, ranked, fmt.Errorf("cannot find students bucket")
	}

	c := students.Cursor()

	users := tx.Bucket([]byte(KeyUsers))

	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		studentData := users.Get([]byte(k))
		var student UserInfo
		err = json.Unmarshal(studentData, &student)
		if err != nil {
			lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
			continue
		}
		if student.Role != UserRoleStudent {
			lgr.Printf("ERROR student %s has role %d", k, student.Role)
			continue
		}

		nWorth, _ := StudentNetWorthTx(tx, student.Name).Float64()
		nUser := openapi.UserNoHistory{
			Id:        student.Email,
			FirstName: student.FirstName,
			LastName:  student.LastName,
			Rank:      student.Rank,
			NetWorth:  float32(nWorth),
		}

		resp = append(resp, nUser)

	}

	sort.SliceStable(resp, func(i, j int) bool {
		return resp[i].NetWorth > resp[j].NetWorth
	})

	for i := 0; i < len(resp); i++ {
		resp[i].Rank = int32(i + 1)
	}

	ranked, err = saveRanksTx(tx, resp)
	if err != nil {
		return resp, ranked, fmt.Errorf("ERROR saving students ranks: %s %v", userDetails.SchoolId, err)
	}

	return
}

func userEdit(db *bolt.DB, clock Clock, userDetails UserInfo, body openapi.UsersUserBody) error {
	return db.Update(func(tx *bolt.Tx) error {
		return userEditTx(tx, clock, userDetails, body)
	})
}

func userEditTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, body openapi.UsersUserBody) error {
	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		return fmt.Errorf("users do not exist")
	}

	user := users.Get([]byte(userDetails.Name))

	if user == nil {
		return fmt.Errorf("user does not exist")
	}

	if body.FirstName != "" {
		userDetails.FirstName = body.FirstName
	}
	if body.LastName != "" {
		userDetails.LastName = body.LastName
	}
	if len(body.Password) > 5 {
		userDetails.PasswordSha = EncodePassword(body.Password)
	}
	if body.CareerTransition && !userDetails.CareerTransition && userDetails.CollegeEnd.IsZero() {
		userDetails.CareerTransition = true
		userDetails.TransitionEnd = clock.Now().AddDate(0, 0, 4) //4 days
		userDetails.Income = userDetails.Income / 2
	}
	if body.College && !userDetails.College {
		diff := decimal.NewFromInt32(8000 - 5000)
		low := decimal.NewFromInt32(5000)
		random := decimal.NewFromFloat32(rand.Float32())

		cost := random.Mul(diff).Add(low).Floor()
		err := chargeStudentUbuckTx(tx, clock, userDetails, cost, "Paying for College", false)
		if err != nil {
			return fmt.Errorf("Failed to chargeStudentUbuckTx: %v", err)
		}
		userDetails.College = true
		userDetails.CollegeEnd = clock.Now().AddDate(0, 0, 14) //14 days
		userDetails.Income = userDetails.Income / 2
	}

	marshal, err := json.Marshal(userDetails)
	if err != nil {
		return fmt.Errorf("Failed to Marshal userDetails")
	}
	err = users.Put([]byte(userDetails.Name), marshal)
	if err != nil {
		return fmt.Errorf("Failed to Put userDetails")
	}
	return nil
}

func executeTransaction(db *bolt.DB, clock Clock, value float32, student, owner, description string) error {
	amount := decimal.NewFromFloat32(value)
	studentDetails, err := getUserInLocalStore(db, student)
	if err != nil {
		return fmt.Errorf("error finding student: %v", err)
	}

	if amount.Sign() > 0 {
		err = addBuck2Student(db, clock, studentDetails, amount, owner, description)
		if err != nil {
			return fmt.Errorf("error paying student: %v", err)
		}
	} else if amount.Sign() < 0 {
		err = chargeStudent(db, clock, studentDetails, amount.Abs(), owner, description, false)
		if err != nil {
			return fmt.Errorf("error debting student: %v", err)
		}
	}

	return nil
}

func executeStudentTransaction(db *bolt.DB, clock Clock, value float32, student string, owner UserInfo, description string) error {
	if student == owner.Name {
		return fmt.Errorf("You can't pay yourself")
	}

	amount := decimal.NewFromFloat32(value)

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("Students can't deduct Ubucks")
	}

	ubucks, err := getStudentUbuck(db, owner)
	if err != nil {
		return err
	}

	if decimal.NewFromFloat32(ubucks.Value).LessThan(amount.Mul(decimal.NewFromFloat32(1.01))) {
		return fmt.Errorf("Hey kid you don't have that much Ubucks, don't forget about 1%% charge")
	}

	studentDetails, err := getUserInLocalStore(db, student)
	if err != nil {
		return fmt.Errorf("error finding student: %v", err)
	}

	err = studentPayStudent(db, clock, amount, studentDetails, owner, description)
	if err != nil {
		return fmt.Errorf("error paying student: %v", err)
	}

	return nil
}

func studentPayStudent(db *bolt.DB, clock Clock, amount decimal.Decimal, reciever, sender UserInfo, description string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return studentPayStudentTx(tx, clock, amount, reciever, sender, description)
	})
}

func studentPayStudentTx(tx *bolt.Tx, clock Clock, amount decimal.Decimal, reciever, sender UserInfo, description string) error {
	if reciever.Role != UserRoleStudent || sender.Role != UserRoleStudent {
		return fmt.Errorf("user is not a student")
	}
	if amount.Sign() <= 0 {
		return fmt.Errorf("amount must be positive")
	}

	ts := clock.Now().Truncate(time.Millisecond)

	transaction := Transaction{
		Ts:             ts,
		Source:         sender.Name,
		Destination:    reciever.Name,
		CurrencySource: CurrencyUBuck,
		CurrencyDest:   CurrencyUBuck,
		AmountSource:   amount.Mul(decimal.NewFromFloat32(1.01)),
		AmountDest:     amount,
		XRate:          decimal.NewFromFloat32(1.0),
		Reference:      description,
		FromSource:     true,
	}

	studentSender, err := getStudentBucketTx(tx, sender.Name)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(studentSender, CurrencyUBuck, transaction, OperationDebit, true)
	if err != nil {
		return err
	}

	cb, err := getCbTx(tx, reciever.SchoolId)
	if err != nil {
		return err
	}

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationCredit, false)
	if err != nil {
		return err
	}

	transaction.FromSource = false

	studentReciever, err := getStudentBucketTx(tx, reciever.Name)
	if err != nil {
		return err
	}
	_, _, err = addToHolderTx(studentReciever, CurrencyUBuck, transaction, OperationCredit, true)
	if err != nil {
		return err
	}

	_, _, err = addToHolderTx(cb, CurrencyUBuck, transaction, OperationDebit, false)
	if err != nil {
		return err
	}

	return nil
}

func getJobRx(tx *bolt.Tx, key string, jobId string) (job openapi.UserNoHistoryJob) {
	jobs := tx.Bucket([]byte(key))
	jobData := jobs.Get([]byte(jobId))

	err := json.Unmarshal(jobData, &job)
	if err != nil {
		return
	}

	job.Title = jobId

	return job
}
