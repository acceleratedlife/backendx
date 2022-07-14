package main

import (
	"context"
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type StaffApiServiceImpl struct {
	db    *bolt.DB
	clock Clock
}

func (s *StaffApiServiceImpl) SearchEvents(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	var resp []openapi.ResponseEvents
	resp, err = getEventsTeacher(s.db, s.clock, userDetails)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil
}

func (a *StaffApiServiceImpl) Deleteclass(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		_, parentBucket, err := getClassAtSchoolTx(tx, userDetails.SchoolId, Id)
		if err != nil {
			return err
		}
		err = parentBucket.DeleteBucket([]byte(Id))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot collect classes from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, nil), nil
}

func (s *StaffApiServiceImpl) DeleteStudent(ctx context.Context, query string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = deleteStudent(s.db, query)

	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, nil), nil
}

func (a *StaffApiServiceImpl) EditClass(ctx context.Context, body openapi.RequestEditClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}
	var class openapi.Class
	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
		if err != nil {
			return err
		}
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		members := make([]string, 0)
		if studentsBucket != nil {
			members, err = studentsToSlice(studentsBucket)
			if err != nil {
				return err
			}
		}

		err = classBucket.Put([]byte(KeyName), []byte(body.Name))
		if err != nil {
			return err
		}
		err = classBucket.Put([]byte(KeyPeriod), itob32(int32(body.Period)))
		if err != nil {
			return err
		}
		class = openapi.Class{
			Id:      body.Id,
			OwnerId: string(classBucket.Get([]byte("OwnerId"))), //OnwerId is legacy, not sure if it is needed
			Period:  btoi32(classBucket.Get([]byte(KeyPeriod))),
			Name:    string(classBucket.Get([]byte(KeyName))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: members,
		}
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot edit class with Id: %s %v", body.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, class), nil
}

func (a *StaffApiServiceImpl) KickClass(ctx context.Context, body openapi.RequestKickClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(a.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, body.Id)
		if err != nil {
			return err
		}

		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		if studentsBucket == nil {
			return fmt.Errorf("can't find students bucket")
		}

		err = studentsBucket.Delete([]byte(body.KickId))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot collect classes from the school: %s %v", userDetails.SchoolId, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, nil), nil
}

func (s *StaffApiServiceImpl) MakeClass(ctx context.Context, request openapi.RequestMakeClass) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	classes, err := s.MakeClassImpl(userDetails, request)

	if err != nil {
		lgr.Printf("ERROR class not created: %v", err)
		return openapi.Response(500, "{}"), nil
	}

	return openapi.Response(200, classes), nil
}

func (s *StaffApiServiceImpl) PayTransactions(ctx context.Context, body openapi.RequestPayTransactions) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	errors := make([]openapi.ResponsePayTransactions, 0)
	if userDetails.Role == UserRoleAdmin {
		body.Owner = body.Owner[:1] + "." + body.Owner[1:]
	}
	for _, student := range body.Students {
		err = executeTransaction(s.db, s.clock, body.Amount, student, body.Owner, body.Description)
		if err != nil {
			errors = append(errors, openapi.ResponsePayTransactions{
				Message: student + " was not paid, error: " + err.Error(),
			})
		}
	}
	if len(errors) != 0 {
		return openapi.Response(400, errors), nil
	}
	return openapi.Response(200, ""), nil
}

func (s *StaffApiServiceImpl) ResetPassword(ctx context.Context, body openapi.RequestUser) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	studentDetails, err := getUserInLocalStore(s.db, body.Id)
	if err != nil {
		return openapi.Response(400, ""), err
	}

	var resp openapi.ResponseResetPassword
	err = s.db.Update(func(tx *bolt.Tx) error {
		resp, err = resetPasswordTx(tx, studentDetails)
		return err
	})

	if err != nil {
		return openapi.Response(400, ""), err
	}

	return openapi.Response(200, resp), nil
}

func (s *StaffApiServiceImpl) SearchAuctionsTeacher(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}

	if userDetails.Role == UserRoleStudent {
		auctions, err := getStudentAuctions(s.db, userDetails)
		if err != nil {
			return openapi.Response(401, ""), err
		}
		var studentAuctions []openapi.Auction
		for i := range auctions {
			if auctions[i].OwnerId.Id == userDetails.Name {
				auctions[i].Visibility = visibilityToSlice(s.db, userDetails, auctions[i].Visibility)
				studentAuctions = append(studentAuctions, auctions[i])
			}
		}

		return openapi.Response(200, studentAuctions), nil
	}

	var resp []openapi.Auction
	err = s.db.View(func(tx *bolt.Tx) error {
		auctionsBucket, err := getAuctionsTx(tx, userDetails)
		if err != nil {
			return err
		}

		resp, err = getTeacherAuctionsRx(tx, auctionsBucket, userDetails)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return openapi.Response(400, nil), err
	}

	return openapi.Response(200, resp), nil
}

func (s *StaffApiServiceImpl) SearchTransactions(ctx context.Context, teacherId string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, openapi.ResponseAuth{
			IsAuth: false,
			Error:  true,
		}), nil
	}
	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	teacherDetails, err := getUserInLocalStore(s.db, teacherId)
	if err != nil {
		return openapi.Response(401, ""), err
	}

	var resp []openapi.ResponseTransactions
	err = s.db.View(func(tx *bolt.Tx) error {
		if userDetails.Role == UserRoleTeacher {
			resp, err = getTeacherTransactionsTx(tx, teacherDetails)
			if err != nil {
				return err
			}
		} else {
			teacherDetails.Name = teacherId[:1] + "." + teacherId[1:]
			resp, err = getTeacherTransactionsTx(tx, teacherDetails)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return openapi.Response(401, ""), err
	}

	return openapi.Response(200, resp), nil
}

// NewStaffApiServiceImpl creates a default api service
func NewStaffApiServiceImpl(db *bolt.DB, clock Clock) openapi.StaffApiServicer {
	return &StaffApiServiceImpl{
		db:    db,
		clock: clock,
	}
}

func studentsToSlice(students *bolt.Bucket) ([]string, error) {
	Members := make([]string, 0)
	cstudents := students.Cursor()
	for k, v := cstudents.First(); k != nil; k, v = cstudents.Next() {
		if v == nil {
			return nil, fmt.Errorf("members of class not found")
		}
		Members = append(Members, string(k))
	}
	return Members, nil
}

func visibilityToSlice(db *bolt.DB, userDetails UserInfo, classIds []string) (resp []string) {
	_ = db.View(func(tx *bolt.Tx) error {
		resp = visibilityToSliceRx(tx, userDetails, classIds)

		return nil
	})

	return
}

func visibilityToSliceRx(tx *bolt.Tx, userDetails UserInfo, classIds []string) (resp []string) {
	for _, key := range classIds {
		if key == KeyEntireSchool || key == KeyFreshman || key == KeySophomores || key == KeyJuniors || key == KeySeniors {
			resp = append(resp, key)
		} else {
			classBucket, _, err := getClassAtSchoolTx(tx, userDetails.SchoolId, key)
			if err != nil {
				continue
			}
			name := classBucket.Get([]byte(KeyName))
			resp = append(resp, string(name))
		}
	}
	return
}
