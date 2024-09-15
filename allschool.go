package main

import (
	"context"
	"fmt"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type AllSchoolApiServiceImpl struct {
	db    *bolt.DB
	clock Clock
}

func (a *AllSchoolApiServiceImpl) AddCodeClass(ctx context.Context, body openapi.RequestUser) (openapi.ImplResponse, error) {
	v := ctx.Value("user")
	if v == nil {
		return openapi.Response(403, nil), nil
	}

	userData := v.(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, nil), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, nil), nil
	}

	addCodes, _ := constSlice()
	addCode := randomWords(2, 100, addCodes)
	free, err := freeAddCode(a.db, addCode)
	if err != nil {
		return openapi.Response(403, nil), err
	}
	for !free {
		addCode = randomWords(2, 100, addCodes)
		free, err = freeAddCode(a.db, addCode)
		if err != nil {
			return openapi.Response(403, nil), err
		}
	}

	if userDetails.Role == UserRoleTeacher {
		newCode := addCode

		class := openapi.ClassWithMembers{
			Id:      userDetails.SchoolId,
			OwnerId: userDetails.Email,
			Period:  0,
			Name:    "Teachers Class",
			AddCode: newCode,
			Members: nil,
		}
		err = a.db.Update(func(tx *bolt.Tx) error {
			//populate class
			teachersClass, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
			if err != nil {
				return err
			}
			err = teachersClass.Put([]byte(KeyAddCode), []byte(newCode))
			if err != nil {
				return err
			}

			endTime := a.clock.Now().Add(time.Minute * 10).Truncate(time.Second).Format(time.RFC3339)

			err = teachersClass.Put([]byte(KeyRegEnd), []byte(endTime))
			if err != nil {
				return err
			}
			return nil
		})

		if err != nil {
			return openapi.Response(500, nil), err
		}
		return openapi.Response(200, class), nil
	}

	if userDetails.Role == UserRoleAdmin {
		newCode := addCode

		class := openapi.ClassWithMembers{
			Id:      userDetails.SchoolId,
			OwnerId: userDetails.Email,
			Period:  -1,
			Name:    "My Teachers",
			AddCode: newCode,
			Members: nil,
		}
		err = a.db.Update(func(tx *bolt.Tx) error {
			if body.Id == userDetails.SchoolId { //school class
				schools, err := tx.CreateBucketIfNotExists([]byte("schools"))
				if err != nil {
					return err
				}
				school, err := schools.CreateBucketIfNotExists([]byte(userDetails.SchoolId))
				if err != nil {
					return err
				}
				err = school.Put([]byte(KeyAddCode), []byte(newCode))
				if err != nil {
					return err
				}

				endTime := a.clock.Now().Add(time.Hour * 72).Truncate(time.Second)

				err = school.Put([]byte(KeyRegEnd), []byte(endTime.Format(time.RFC3339)))
				if err != nil {
					return err
				}
			} else { // admin class
				classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
				if err != nil {
					return err
				}
				err = classBucket.Put([]byte(KeyAddCode), []byte(newCode))
				if err != nil {
					return err
				}

				endTime := a.clock.Now().Add(time.Minute * 10).Truncate(time.Second)

				err = classBucket.Put([]byte(KeyRegEnd), []byte(endTime.Format(time.RFC3339)))
				if err != nil {
					return err
				}
			}
			return nil
		})

		if err != nil {
			return openapi.Response(500, nil), err
		}
		return openapi.Response(200, class), nil
	}
	return openapi.Response(500, nil), fmt.Errorf("user's role is not defined")
}

func (a *AllSchoolApiServiceImpl) RemoveClass(ctx context.Context, body openapi.RequestKickClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role != UserRoleStudent {
		return openapi.Response(200, nil), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
		if err != nil {
			return err
		}
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		err = studentsBucket.Delete([]byte(userDetails.Name))
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

func (a AllSchoolApiServiceImpl) SearchAuctions(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a *AllSchoolApiServiceImpl) SearchMyClasses(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	resp := make([]openapi.ResponseMemberClass, 0)
	if Id == "" {
		resp, err = classesWithOwnerDetails(a.db, userDetails.SchoolId, userDetails.Email)
	} else {
		resp, err = classesWithOwnerDetails(a.db, userDetails.SchoolId, Id)
	}

	if err != nil {
		return openapi.Response(500, "{}"), nil
	}

	return openapi.Response(200, resp), nil
}

// NewAllSchoolApiServiceImpl creates a default api service
func NewAllSchoolApiServiceImpl(db *bolt.DB, clock Clock) openapi.AllSchoolApiServicer {
	return &AllSchoolApiServiceImpl{
		db:    db,
		clock: clock,
	}
}

func classesWithOwnerDetails(db *bolt.DB, schoolID, userId string) ([]openapi.ResponseMemberClass, error) {
	classes := make([]openapi.ResponseMemberClass, 0)
	err := db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolID)
		if err != nil {
			return err
		}
		teachers := school.Bucket([]byte(KeyTeachers))
		if teachers == nil {
			return fmt.Errorf("no teachers at school")
		}
		iterateBuckets(teachers, func(teacher *bolt.Bucket, teacherId []byte) {
			iterateBuckets(teacher, func(class *bolt.Bucket, classId []byte) {
				students := class.Bucket([]byte(KeyStudents))
				if students == nil {
					return
				}
				student := students.Get([]byte(userId))
				if student == nil {
					return
				}
				teacherDetails, err := getUserInLocalStoreTx(tx, string(teacherId))
				if err != nil {
					lgr.Printf("ERROR cannot find details for teacher %x", teacherId)
					return
				}
				period := class.Get([]byte(KeyPeriod))

				classes = append(classes, openapi.ResponseMemberClass{
					OwnerId: openapi.ResponseMemberClassOwner{
						FirstName: teacherDetails.FirstName,
						LastName:  teacherDetails.LastName,
						Id:        teacherDetails.Name,
					},
					Period: btoi32(period),
					Id:     string(classId),
				})
			})
		})
		schoolClasses := school.Bucket([]byte(KeyClasses))
		if schoolClasses == nil {
			return fmt.Errorf("no school classes at school")
		}
		iterateBuckets(schoolClasses, func(class *bolt.Bucket, classId []byte) {
			students := class.Bucket([]byte(KeyStudents))
			if students == nil {
				return
			}
			student := students.Get([]byte(userId))
			if student == nil {
				return
			}
			adminBucket := school.Bucket([]byte(KeyAdmins))
			adminId, _ := adminBucket.Cursor().First()
			if adminId == nil {
				lgr.Printf("ERROR cannot find bucket for admin %x", adminId)
				return
			}
			adminDetails, err := getUserInLocalStoreTx(tx, string(adminId))
			if err != nil {
				lgr.Printf("ERROR cannot find details for teacher %x", adminId)
				return
			}
			period := class.Get([]byte(KeyPeriod))

			classes = append(classes, openapi.ResponseMemberClass{
				OwnerId: openapi.ResponseMemberClassOwner{
					FirstName: adminDetails.FirstName,
					LastName:  adminDetails.LastName,
					Id:        adminDetails.Name,
				},
				Period: btoi32(period),
				Id:     string(classId),
			})
		})
		return nil
	})

	return classes, err
}
