package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

func deleteStudent(db *bolt.DB, studentId string) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		return deleteStudentTx(tx, studentId)
	})

	return
}

func deleteStudentTx(tx *bolt.Tx, studentId string) (err error) {
	studentInfo, err := getUserInLocalStoreTx(tx, studentId)
	if err != nil {
		return err
	}

	auctionsBucket, err := getAuctionsTx(tx, studentInfo)
	if err != nil {
		return err
	}

	auctions, err := getStudentAuctionsTx(tx, auctionsBucket, studentInfo)
	if err != nil {
		return err
	}

	for _, c := range auctions {
		if c.OwnerId.Id == studentInfo.Name {
			// repayAuctionLoser()
			err := auctionsBucket.DeleteBucket([]byte(c.Id.String()))
			if err != nil {
				return err
			}
		}
	}

	classes, err := getStudentClassesTx(tx, studentInfo)
	if err != nil {
		return err
	}

	for _, c := range classes {
		classBucket, _, err := getClassAtSchoolTx(tx, studentInfo.SchoolId, c.Id)
		if err != nil {
			return err
		}
		students := classBucket.Bucket([]byte(KeyStudents))
		students.Delete([]byte(studentInfo.Name))
		studentFound := students.Get([]byte(studentInfo.Name))
		if studentFound != nil {
			return fmt.Errorf("failed to delete student from class: %v", c.Name)
		}
	}

	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		return fmt.Errorf("cannot get users bucket")
	}

	users.Delete([]byte(studentId))
	user := users.Get([]byte(studentId))
	if user != nil {
		return fmt.Errorf("failed to delete user")
	}

	school, err := getSchoolBucketTx(tx, studentInfo)
	if err != nil {
		return fmt.Errorf("cannot get school: %v", err)
	}
	students := school.Bucket([]byte(KeyStudents))
	if students == nil {
		return fmt.Errorf("cannot get students bucket")
	}

	err = students.DeleteBucket([]byte(studentId))
	if err != nil {
		return err
	}

	return

}

func getSchoolClasses(db *bolt.DB, schoolId string) (res []openapi.Class) {
	_ = db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolId)
		if err != nil {
			return err
		}
		if school == nil {
			return nil
		}
		classes := school.Bucket([]byte(KeyClasses))
		if classes == nil {
			return nil
		}
		res = getClasses1Tx(classes, "")
		return nil
	})
	return
}

func getTeacherClasses(db *bolt.DB, schoolId, teacherId string) (res []openapi.Class) {
	_ = db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolId)
		if err != nil {
			return err
		}
		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			return nil
		}
		teacher := teachers.Bucket([]byte(teacherId))
		if teacher == nil {
			return nil
		}
		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			return err
		}
		res = getClasses1Tx(classesBucket, teacherId)
		return nil
	})
	return
}

func getTeacherBucketTx(tx *bolt.Tx, schoolId, teacherId string) (teacher *bolt.Bucket, err error) {
	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return teacher, fmt.Errorf("Cannot find school")
	}
	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return teacher, fmt.Errorf("Cannot find teachers")
	}
	teacher = teachers.Bucket([]byte(teacherId))
	if teacher == nil {
		return teacher, fmt.Errorf("Cannot find teacher")
	}
	return
}

func getClassesTx(classesBucket *bolt.Bucket) []openapi.Class {
	classes := make([]openapi.Class, 0)

	c := classesBucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}
		classBucket := classesBucket.Bucket(k)
		iClass := openapi.Class{
			Id:      string(k),
			Name:    string(classBucket.Get([]byte(KeyName))),
			OwnerId: "",
			Period:  btoi32(classBucket.Get([]byte(KeyPeriod))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: make([]string, 0),
		}
		classes = append(classes, iClass)
	}
	return classes
}

func getAuctionsTx(tx *bolt.Tx, userDetails UserInfo) (auctionsBucket *bolt.Bucket, err error) {
	school, err := SchoolByIdTx(tx, userDetails.SchoolId)
	if err != nil {
		return auctionsBucket, err
	}
	auctionsBucket = school.Bucket([]byte(KeyAuctions))
	if err != nil {
		return auctionsBucket, nil
	}

	return
}

func getAuctionBucket(db *bolt.DB, schoolBucket *bolt.Bucket, auctionId string) (auctionsBucket *bolt.Bucket, auctionBucket []byte, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		auctionsBucket, auctionBucket, err = getAuctionBucketTx(tx, schoolBucket, auctionId)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return
}

func getAuctionBucketTx(tx *bolt.Tx, schoolBucket *bolt.Bucket, auctionId string) (auctionsBucket *bolt.Bucket, auctionByte []byte, err error) {
	auctionsBucket = schoolBucket.Bucket([]byte(KeyAuctions))
	if auctionsBucket == nil {
		return auctionsBucket, auctionByte, fmt.Errorf("cannot find auctions bucket")
	}

	auctionByte = auctionsBucket.Get([]byte(auctionId))
	if auctionByte == nil {
		return auctionsBucket, auctionByte, fmt.Errorf("cannot find auction bucket, auction: " + auctionId)
	}

	return
}

func getSchoolBucket(db *bolt.DB, userDetails UserInfo) (school *bolt.Bucket, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		school, err = getSchoolBucketTx(tx, userDetails)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return
}

func getSchoolBucketTx(tx *bolt.Tx, userDetails UserInfo) (school *bolt.Bucket, err error) {
	schools := tx.Bucket([]byte(KeySchools))
	if schools == nil {
		return nil, fmt.Errorf("cannot find schools bucket")
	}

	school = schools.Bucket([]byte(userDetails.SchoolId))
	if school == nil {
		return nil, fmt.Errorf("cannot find school bucket")
	}

	return
}

func getTeacherAuctionsTx(tx *bolt.Tx, auctionsBucket *bolt.Bucket, userDetails UserInfo) ([]openapi.Auction, error) {
	auctions := make([]openapi.Auction, 0)

	c := auctionsBucket.Cursor()

	for k, _ := c.First(); k != nil; k, _ = c.Next() {

		auctionBucket := auctionsBucket.Get(k)
		var auction openapi.Auction
		err := json.Unmarshal(auctionBucket, &auction)
		if err != nil {
			return nil, err
		}

		if auction.OwnerId.Id == userDetails.Name {
			auction.Visibility = visibilityToSlice(tx, userDetails, auction.Visibility)

			ownerDetails, err := getUserInLocalStoreTx(tx, auction.OwnerId.Id)
			if err != nil {
				return nil, err
			}

			auction.OwnerId.LastName = ownerDetails.LastName

			winnerDetails, err := getUserInLocalStoreTx(tx, auction.WinnerId.Id)
			if err != nil {
				auction.WinnerId = openapi.AuctionWinnerId{
					FirstName: "nil",
					LastName:  "nil",
					Id:        "nil",
				}
			} else {
				auction.WinnerId.FirstName = winnerDetails.FirstName
				auction.WinnerId.LastName = winnerDetails.LastName
			}

			auctions = append(auctions, auction)
			if len(auctions) > 49 {
				break
			}
		}

	}

	return auctions, nil
}

func getClasses1Tx(classesBucket *bolt.Bucket, ownerId string) []openapi.Class {
	data := make([]openapi.Class, 0)
	c := classesBucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}
		classBucket := classesBucket.Bucket(k)
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		members := make([]string, 0)
		if studentsBucket != nil {
			innerMembers, err := studentsToSlice(studentsBucket)
			if err != nil {
				lgr.Printf("ERROR cannot turn students to slice: %s %v", ownerId, err)
			} else {
				members = innerMembers
			}
		}

		iClass := openapi.Class{
			Id:      string(k),
			OwnerId: ownerId,
			Period:  btoi32(classBucket.Get([]byte(KeyPeriod))),
			Name:    string(classBucket.Get([]byte(KeyName))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: members,
		}
		data = append(data, iClass)
	}
	return data
}

func (s *StaffApiServiceImpl) MakeClassImpl(userDetails UserInfo, request openapi.RequestMakeClass) (classes []openapi.Class, err error) {
	schoolId := userDetails.SchoolId
	teacherId := userDetails.Name
	className := request.Name
	period := request.Period

	_, classes, err = CreateClass(s.db, schoolId, teacherId, className, int(period))

	return classes, err
}

func CreateClass(db *bolt.DB, schoolId, teacherId, className string, period int) (classId string, classes []openapi.Class, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolId)
		if err != nil {
			return err
		}
		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			return fmt.Errorf("user does not exist")
		}
		teacher := teachers.Bucket([]byte(teacherId))
		if teacher == nil {
			return fmt.Errorf("user does not exist")
		}

		classesBucket := teacher.Bucket([]byte(KeyClasses))
		if classesBucket == nil {
			return fmt.Errorf("Problem finding classesBucket")
		}

		classId, err = addClassDetailsTx(classesBucket, className, period, false)
		if err != nil {
			return err
		}

		classes = getClassesTx(classesBucket)

		return nil
	})
	return
}

func MakeAuctionImpl(db *bolt.DB, userDetails UserInfo, request openapi.RequestMakeAuction) (auctions []openapi.Auction, err error) {
	_, auctions, err = CreateAuction(db, userDetails, request)
	if err != nil {
		return nil, fmt.Errorf("cannot create auction")
	}

	return auctions, err
}

func CreateAuction(db *bolt.DB, userDetails UserInfo, request openapi.RequestMakeAuction) (auctionId time.Time, auctions []openapi.Auction, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, userDetails.SchoolId)
		if err != nil {
			return fmt.Errorf("Problem finding auctions bucket: %v", err)
		}
		auctionsBucket := school.Bucket([]byte(KeyAuctions))
		if auctionsBucket == nil {
			return fmt.Errorf("Problem finding auctions bucket")
		}

		auctionId, err = addAuctionDetailsTx(auctionsBucket, request)
		if err != nil {
			return fmt.Errorf("Problem adding auctions details: %v", err)
		}

		auctions, err = getTeacherAuctionsTx(tx, auctionsBucket, userDetails)
		if err != nil {
			return fmt.Errorf("Problem finding teacher auctions: %v", err)
		}

		return nil
	})

	if err != nil {
		return auctionId, nil, err
	}

	return
}

func addAuctionDetailsTx(bucket *bolt.Bucket, request openapi.RequestMakeAuction) (auctionId time.Time, err error) {
	auctionId = request.EndDate.Truncate(time.Millisecond)
	found := bucket.Get([]byte(auctionId.String()))
	for found != nil {
		auctionId = auctionId.Add(time.Millisecond * 1)
		found = bucket.Get([]byte(auctionId.String()))
	}

	auction := openapi.Auction{
		Id:          auctionId,
		Active:      true,
		StartDate:   request.StartDate.Truncate(time.Millisecond),
		EndDate:     auctionId,
		Bid:         int32(request.MaxBid),
		MaxBid:      int32(request.MaxBid),
		Description: request.Description,
		Visibility:  request.Visibility,
		OwnerId: openapi.AuctionOwnerId{
			Id: request.OwnerId,
		},
	}

	marshal, err := json.Marshal(auction)
	if err != nil {
		return
	}

	err = bucket.Put([]byte(auction.Id.String()), marshal)
	if err != nil {
		return
	}

	return
}

func addClassDetailsTx(bucket *bolt.Bucket, className string, period int, adminClass bool) (classId string, err error) {
	if adminClass {
		classId = className
	} else {
		classId = RandomString(15)
	}
	class, err1 := bucket.CreateBucket([]byte(classId))
	if err1 != nil {
		return "", err1
	}

	err = class.Put([]byte(KeyName), []byte(className))
	if err != nil {
		return "", err
	}
	err = class.Put([]byte(KeyPeriod), itob32(int32(period)))
	if err != nil {
		return "", err
	}
	addCode := RandomString(6)
	err = class.Put([]byte(KeyAddCode), []byte(addCode))
	return
}

func getTeacherTransactionsTx(tx *bolt.Tx, teacher UserInfo) (resp []openapi.ResponseTransactions, err error) {
	CB, err := getCbRx(tx, teacher.SchoolId)
	if err != nil {
		return
	}

	accounts := CB.Bucket([]byte(KeyAccounts))
	if accounts == nil {
		return resp, fmt.Errorf("Cannot find buck accounts bucket")
	}

	buck := accounts.Bucket([]byte(teacher.Name))
	if buck == nil {
		lgr.Printf("Cannot find " + teacher.LastName + " buck bucket")
		return resp, nil
	}

	transactions := buck.Bucket([]byte(KeyTransactions))
	if transactions == nil {
		return resp, fmt.Errorf("Cannot find transactions bucket")
	}

	c := transactions.Cursor()
	for k, v := c.Last(); k != nil; k, v = c.Prev() {
		if v == nil {
			continue
		}

		trans := parseTransactionStudent(v, teacher, teacher.Name)

		studentId := trans.Source
		if trans.Destination != "" {
			studentId = trans.Destination
		}

		student, err := getUserInLocalStoreTx(tx, studentId)
		if err != nil {
			student.FirstName = "Deleted"
			student.LastName = "Student"
		}

		slice := openapi.ResponseTransactions{
			Amount:      float32(trans.Net.InexactFloat64()),
			CreatedAt:   trans.Ts,
			Description: trans.Reference,
			Student:     student.FirstName + " " + student.LastName,
		}

		resp = append(resp, slice)
		if len(resp) >= 60 {
			break
		}
	}

	return
}

func getEventsTeacher(db *bolt.DB, clock Clock, userDetails UserInfo) (resp []openapi.ResponseEvents, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		cb, err := getCbRx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}

		accounts := cb.Bucket([]byte(KeyAccounts))

		ubuck := accounts.Bucket([]byte(CurrencyUBuck))
		transactions := ubuck.Bucket([]byte(KeyTransactions))

		c := transactions.Cursor()
		var trans Transaction
		for k, _ := c.First(); k != nil; k, _ = c.Next() {

			transTime, err := time.Parse(time.RFC3339, string(k))
			if err != nil {
				return err
			}

			if transTime.Before(clock.Now().Truncate(time.Hour * 24)) {
				break
			}

			transData := transactions.Get(k)
			err = json.Unmarshal(transData, &trans)
			if err != nil {
				return err
			}

			if !strings.Contains(trans.Reference, "Event: ") {
				continue
			}

			var student UserInfo
			if trans.Source != "" {
				student, err = getUserInLocalStoreTx(tx, trans.Source)
			} else {
				student, err = getUserInLocalStoreTx(tx, trans.Destination)
			}

			resp = append(resp, openapi.ResponseEvents{
				Value:       int32(trans.AmountSource.IntPart()),
				Description: trans.Reference,
				FirstName:   student.FirstName,
				LastName:    student.LastName,
			})
		}

		return nil
	})

	return
}
