package main

import (
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

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

func getAuctions(db *bolt.DB, userDetails UserInfo) (resp []openapi.Auction, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}
		auctionsBucket := school.Bucket([]byte(KeyAuctions))
		if auctionsBucket == nil {
			return fmt.Errorf("Cannot find auctions")
		}

		resp, err = getAuctionsTx(tx, auctionsBucket, userDetails.Name)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return
	}

	return
}

func getAuctionBucket(db *bolt.DB, schoolBucket *bolt.Bucket, auctionId string) (auctionsBucket, auctionBucket *bolt.Bucket, err error) {
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

func getAuctionBucketTx(tx *bolt.Tx, schoolBucket *bolt.Bucket, auctionId string) (auctionsBucket, auctionBucket *bolt.Bucket, err error) {
	auctionsBucket = schoolBucket.Bucket([]byte(KeyAuctions))
	if auctionsBucket == nil {
		return auctionsBucket, auctionBucket, fmt.Errorf("cannot find auctions bucket")
	}

	auctionBucket = auctionsBucket.Bucket([]byte(auctionId))
	if auctionBucket == nil {
		return auctionsBucket, auctionBucket, fmt.Errorf("cannot find auction bucket")
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

func getAuctionsTx(tx *bolt.Tx, auctionsBucket *bolt.Bucket, ownerId string) ([]openapi.Auction, error) {
	auctions := make([]openapi.Auction, 0)

	c := auctionsBucket.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}

		auctionBucket := auctionsBucket.Bucket(k)
		if string(auctionBucket.Get([]byte(KeyOwnerId))) == ownerId {
			iAuction := openapi.Auction{
				Id:          string(k),
				StartDate:   string(auctionBucket.Get([]byte(KeyStartDate))),
				EndDate:     string(auctionBucket.Get([]byte(KeyEndDate))),
				ItemNumber:  RandomString(3),
				Bid:         btoi32(auctionBucket.Get([]byte(KeyBid))),
				MaxBid:      btoi32(auctionBucket.Get([]byte(KeyMaxBid))),
				Description: string(auctionBucket.Get([]byte(KeyDescription))),
				Visibility:  visibilityToSlice(auctionBucket.Bucket([]byte(KeyVisibility))),
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
					FirstName: winnerDetails.FirstName,
					LastName:  winnerDetails.LastName,
					Id:        winnerDetails.Name,
				}
			}

			auctions = append(auctions, iAuction)
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

func (s *StaffApiServiceImpl) MakeAuctionImpl(userDetails UserInfo, request openapi.RequestMakeAuction) (auctions []openapi.Auction, err error) {
	_, auctions, err = CreateAuction(s.db, userDetails, request)

	return auctions, err
}

func CreateAuction(db *bolt.DB, userDetails UserInfo, request openapi.RequestMakeAuction) (auctionId string, auctions []openapi.Auction, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}
		auctionsBucket, err := school.CreateBucketIfNotExists([]byte(KeyAuctions))
		if err != nil {
			return err
		}

		auctionId, err = addAuctionDetailsTx(auctionsBucket, request)
		if err != nil {
			return err
		}

		auctions, err = getAuctionsTx(tx, auctionsBucket, userDetails.Name)
		if err != nil {
			return err
		}

		return nil
	})

	return
}

func addAuctionDetailsTx(bucket *bolt.Bucket, request openapi.RequestMakeAuction) (auctionId string, err error) {
	// auctionId = RandomString(15)
	auctionId = request.EndDate.String()
	auction, err1 := bucket.CreateBucket([]byte(auctionId)) //what happens if 2 actions are made to end at the same time?
	if err1 != nil {
		return "", err1
	}

	err = auction.Put([]byte(KeyBid), itob32(int32(request.Bid)))
	if err != nil {
		return "", err
	}
	err = auction.Put([]byte(KeyMaxBid), itob32(int32(request.MaxBid)))
	if err != nil {
		return "", err
	}
	err = auction.Put([]byte(KeyDescription), []byte(request.Description))
	if err != nil {
		return "", err
	}
	err = auction.Put([]byte(KeyEndDate), []byte(request.EndDate.String()))
	if err != nil {
		return "", err
	}
	err = auction.Put([]byte(KeyStartDate), []byte(request.StartDate.String()))
	if err != nil {
		return "", err
	}
	err = auction.Put([]byte(KeyOwnerId), []byte(request.OwnerId))
	if err != nil {
		return "", err
	}
	visibility, err := auction.CreateBucket([]byte(KeyVisibility))
	if err != nil {
		return "", err
	}

	for _, s := range request.Visibility {
		visibility.CreateBucket([]byte(s))
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
