package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	"github.com/shopspring/decimal"
	bolt "go.etcd.io/bbolt"
)

func getClassAtSchoolTx(tx *bolt.Tx, schoolId, classId string) (classBucket *bolt.Bucket, parentBucket *bolt.Bucket, err error) {

	school, err := schoolByIdTx(tx, schoolId)
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
		return nil, fmt.Errorf("failed to get history")
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

	resp.Id = accountId
	balanceData := bAccount.Get([]byte(KeyBalance))
	if balanceData == nil {
		resp.Balance = 0
		return
	}

	err = json.Unmarshal(balanceData, &resp.Balance)
	return
}

// looks to see if a cb account for this account exists and if it does then return the account given to it
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

// get all the students from a school, update the rank to the top students
func getSchoolStudents(db *bolt.DB, userDetails UserInfo) (resp []openapi.UserNoHistory, ranked int, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		resp, ranked, err = getSchoolStudentsRx(tx, userDetails)
		return err
	})

	return
}

func getSchoolStudentsRx(tx *bolt.Tx, userDetails UserInfo) (resp []openapi.UserNoHistory, ranked int, err error) {
	school, err := schoolByIdTx(tx, userDetails.SchoolId)
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

		nUser := openapi.UserNoHistory{
			Id:        student.Email,
			FirstName: student.FirstName,
			LastName:  student.LastName,
			Rank:      student.Rank,
			NetWorth:  student.NetWorth,
		}

		resp = append(resp, nUser)

	}

	sort.SliceStable(resp, func(i, j int) bool {
		return resp[i].NetWorth > resp[j].NetWorth
	})

	for i := 0; i < len(resp); i++ {
		resp[i].Rank = int32(i + 1)
	}

	ranked, err = limitRanks(resp)
	if err != nil {
		return resp, ranked, fmt.Errorf("ERROR saving students ranks: %s %v", userDetails.SchoolId, err)
	}

	return
}

func limitRanks(students []openapi.UserNoHistory) (ranked int, err error) {

	length := float32(len(students))

	if length <= 60 {
		return int(.33334 * length), err
	} else if length <= 150 {
		return int(.2 * length), err
	} else if length <= 500 {
		return int(.15 * length), err
	} else if length <= 1000 {
		return int(.1 * length), err
	} else if length <= 2000 {
		return int(.0625 * length), err
	} else if length <= 4000 {
		return int(.05 * length), err
	} else if length > 4000 {
		return int(200), err
	}

	return ranked, err

}

func userEdit(db *bolt.DB, clock Clock, userDetails UserInfo, body openapi.RequestUserEdit) error {
	return db.Update(func(tx *bolt.Tx) error {
		return userEditTx(tx, clock, userDetails, body)
	})
}

func userEditTx(tx *bolt.Tx, clock Clock, userDetails UserInfo, body openapi.RequestUserEdit) error {
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
		if userDetails.Role != UserRoleStudent {
			userDetails.LastName = body.LastName
		} else {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			userDetails.LastName = string(body.LastName[0]) + strconv.Itoa(r.Intn(10000))
		}
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
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		random := decimal.NewFromFloat32(r.Float32())

		cost := random.Mul(diff).Add(low).Floor()
		err := chargeStudentUbuckTx(tx, clock, userDetails, cost, "Paying for College", false)
		if err != nil {
			return fmt.Errorf("failed to chargeStudentUbuckTx: %v", err)
		}
		userDetails.College = true
		userDetails.CollegeEnd = clock.Now().AddDate(0, 0, 14) //14 days
		userDetails.Income = userDetails.Income / 2
	}

	if body.LottoPlay > 0 {
		if userDetails.LottoPlay == 0 {
			userDetails.LottoPlay = body.LottoPlay
		} else {
			userDetails.LottoPlay += body.LottoPlay
		}
	}

	if body.LottoWin > 0 {
		if userDetails.LottoWin == 0 {
			userDetails.LottoWin = body.LottoWin
		} else {
			userDetails.LottoWin += body.LottoWin
		}
	}

	marshal, err := json.Marshal(userDetails)
	if err != nil {
		return fmt.Errorf("failed to marshal userDetails")
	}
	err = users.Put([]byte(userDetails.Name), marshal)
	if err != nil {
		return fmt.Errorf("failed to put userDetails")
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

func garnishHelper(db *bolt.DB, clock Clock, body openapi.RequestPayTransaction, isStaff bool) (garnish decimal.Decimal, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		garnish, err = garnishHelperTx(tx, clock, body, isStaff)
		return err
	})

	return
}

func garnishHelperTx(tx *bolt.Tx, clock Clock, body openapi.RequestPayTransaction, isStaff bool) (garnish decimal.Decimal, err error) {
	student, err := getStudentBucketTx(tx, body.Student)
	if err != nil {
		return
	}

	_, _, balance, err := IsDebtNeededRx(student, clock)
	if err != nil {
		return
	}

	if !balance.GreaterThan(decimal.Zero) {
		return
	}

	garnish = decimal.NewFromFloat32(body.Amount).Mul(decimal.NewFromFloat32(KeyGarnish))
	userDetails, err := getUserInLocalStoreTx(tx, body.Student)
	if err != nil {
		return
	}

	if garnish.GreaterThan(balance) {
		garnish = balance
	}

	if isStaff {
		err = studentConvertTx(tx, clock, userDetails, garnish, body.OwnerId, KeyDebt, "Garneshed "+strconv.FormatFloat(KeyGarnish*100, 'f', 0, 32)+"% of previous payment", true)
	} else {
		err = studentConvertTx(tx, clock, userDetails, garnish, CurrencyUBuck, KeyDebt, "Garneshed "+strconv.FormatFloat(KeyGarnish*100, 'f', 0, 32)+"% of previous payment", true)
	}

	return
}

func executeStudentTransaction(db *bolt.DB, clock Clock, value float32, student string, owner UserInfo, description string) error {
	if student == owner.Name {
		return fmt.Errorf("you can't pay yourself")
	}

	amount := decimal.NewFromFloat32(value)

	if amount.LessThanOrEqual(decimal.Zero) {
		return fmt.Errorf("students can't deduct Ubucks")
	}

	ubucks, err := getStudentUbuck(db, owner)
	if err != nil {
		return err
	}

	if decimal.NewFromFloat32(ubucks.Value).LessThan(amount.Mul(decimal.NewFromFloat32(keyCharge))) {
		return fmt.Errorf("hey kid you don't have that much Ubucks, don't forget about 1%% charge")
	}

	studentDetails, err := getUserInLocalStore(db, student)
	if err != nil {
		return fmt.Errorf("error finding student: %v", err)
	}

	err = studentPayStudent(db, clock, amount, studentDetails, owner, "from "+owner.FirstName+" "+owner.LastName+" to "+studentDetails.FirstName+" "+studentDetails.LastName)
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
		AmountSource:   amount.Mul(decimal.NewFromFloat32(keyCharge)),
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

func getJob(db *bolt.DB, key string, jobId string) (job openapi.UserNoHistoryJob) {
	_ = db.View(func(tx *bolt.Tx) error {
		job = getJobRx(tx, key, jobId)
		return nil
	})

	return
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

func deleteAuction(db *bolt.DB, userDetails UserInfo, clock Clock, Id time.Time) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		return deleteAuctionTx(tx, userDetails, clock, Id)
	})

	return
}

func deleteAuctionTx(tx *bolt.Tx, userDetails UserInfo, clock Clock, Id time.Time) (err error) {

	newId := Id.Truncate(time.Millisecond).String()

	schoolBucket, err := getSchoolBucketRx(tx, userDetails)
	if err != nil {
		return err
	}

	auctionsBucket, auctionData, err := getAuctionBucketTx(tx, schoolBucket, newId)
	if err != nil {
		return err
	}

	var auction openapi.Auction
	err = json.Unmarshal(auctionData, &auction)
	if err != nil {
		return err
	}

	ownerDetails, err := getUserInLocalStoreTx(tx, auction.OwnerId.Id)

	if userDetails.Role != UserRoleStudent && ownerDetails.Role != UserRoleStudent && userDetails.Email != ownerDetails.Email {
		return fmt.Errorf("you can't cancel the auction of another staff")
	}

	if clock.Now().Before(auction.EndDate) { //auction is not over
		if auction.WinnerId.Id != "" { //currently a winner
			err = repayLosertx(tx, clock, auction.WinnerId.Id, auction.MaxBid, "Canceled Auction: "+strconv.Itoa(auction.EndDate.Second()))
			if err != nil {
				return err
			}
		}

		err = auctionsBucket.Delete([]byte(newId))
		if err != nil {
			return err
		}

	} else { // auction is over
		if auction.WinnerId.Id != "" { //currently a winner
			if auction.MaxBid > auction.Bid {
				err = repayLosertx(tx, clock, auction.WinnerId.Id, auction.MaxBid-auction.Bid, "Won auction return: "+strconv.Itoa(auction.EndDate.Minute()))
				if err != nil {
					return err
				}
			}

			auction.Active = false
			marshal, err := json.Marshal(auction)
			if err != nil {
				return err
			}

			err = auctionsBucket.Put([]byte(newId), marshal)
			if err != nil {
				return err
			}

			sellerDetails := userDetails
			if userDetails.Role != UserRoleStudent {
				sellerDetails, err = getUserInLocalStoreTx(tx, auction.OwnerId.Id)
				if err != nil {
					return err
				}
			}

			if sellerDetails.Role == UserRoleStudent {
				err = addUbuck2StudentTx(tx, clock, sellerDetails, decimal.NewFromInt32(auction.Bid).Mul(decimal.NewFromFloat32(.99)), "Auction sold: "+strconv.Itoa(auction.EndDate.Minute()))

				if err != nil {
					return err
				}

				_, err := garnishHelperTx(tx, clock, openapi.RequestPayTransaction{
					Amount:  float32(auction.Bid) * .99,
					Student: sellerDetails.Name,
				}, false)

				if err != nil {
					return err
				}

			}

		} else { // over and has no winner
			err = auctionsBucket.Delete([]byte(newId))
			if err != nil {
				return err
			}
		}
	}

	return

}

func getTeachers(db *bolt.DB, userDetails UserInfo) (teachers []openapi.ResponseTeachers, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		school, err := getSchoolBucketRx(tx, userDetails)
		if err != nil {
			return err
		}

		teachersBucket := school.Bucket([]byte(KeyTeachers))
		if teachersBucket == nil {
			return fmt.Errorf("cannot find teachers bucket")
		}

		c := teachersBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			teacher, err := getUserInLocalStoreTx(tx, string(k))
			if err != nil {
				return err
			}

			teachers = append(teachers, openapi.ResponseTeachers{
				Id:   teacher.Email,
				Name: teacher.LastName + " " + teacher.FirstName[0:1],
			})

		}

		return nil

	})

	return

}
