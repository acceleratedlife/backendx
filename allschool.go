package main

import (
	"context"
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type AllSchoolApiServiceImpl struct {
	db *bolt.DB
}

func (a *AllSchoolApiServiceImpl) AddCodeClass(ctx context.Context, body openapi.RequestUser) (openapi.ImplResponse, error) {
	//type contextKey string
	//userData := ctx.Value("user").(token.User)
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

	if userDetails.Role == UserRoleTeacher {
		// println(body)
		return openapi.Response(404, nil), fmt.Errorf("not implemented")
	}

	if userDetails.Role == UserRoleAdmin {
		newCode := RandomString(6)

		class := openapi.ClassWithMembers{
			Id:      userDetails.SchoolId,
			OwnerId: userDetails.Email,
			Period:  0,
			Name:    "",
			AddCode: newCode,
			Members: nil,
		}
		err = a.db.Update(func(tx *bolt.Tx) error {
			//populate class
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

	resp := make([]openapi.ResponseMemberClass, 0)
	err = a.db.Update(func(tx *bolt.Tx) error {
		classBucket, err := ClassForAll(tx, body.Id)
		if err != nil {
			return err
		}
		studentsBucket := classBucket.Bucket([]byte(KeyStudents))
		err = studentsBucket.Delete([]byte(body.KickId))
		if err != nil {
			return err
		}
		classes, err := classesWithOwnerDetails(a.db, userDetails.SchoolId, userDetails.Email)
		if err != nil {
			return err
		}
		resp = classes
		return nil
	})
	if err != nil {
		lgr.Printf("ERROR cannot edit class with Id: %s %v", body.Id, err)
		return openapi.Response(500, "{}"), nil
	}
	return openapi.Response(200, resp), nil
}

func (a AllSchoolApiServiceImpl) SearchAuctions(ctx context.Context, s string) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AllSchoolApiServiceImpl) SearchMyClasses(ctx context.Context, query openapi.RequestUser) (openapi.ImplResponse, error) {
	//TODO implement me
	panic("implement me")
}

// NewAllSchoolApiServiceImpl creates a default api service
func NewAllSchoolApiServiceImpl(db *bolt.DB) openapi.AllSchoolApiServicer {
	return &AllSchoolApiServiceImpl{
		db: db,
	}
}

func classesWithOwnerDetails(db *bolt.DB, schoolID, userId string) ([]openapi.ResponseMemberClass, error) {
	sClasses, err := getClassMembership(db, schoolID, userId)
	if err != nil {
		return nil, err
	}
	classes := make([]openapi.ResponseMemberClass, 0)
	for _, currentClass := range sClasses {
		ownerDetails, err := getClassOwner(db, currentClass.Id, schoolID)
		if err != nil {
			return nil, err
		}
		owner := openapi.ResponseMemberClassOwner{
			FirstName: ownerDetails.FirstName,
			LastName:  ownerDetails.LastName,
			Id:        ownerDetails.Email,
		}
		class := openapi.ResponseMemberClass{
			Id:     currentClass.Id,
			Owner:  owner,
			Period: currentClass.Period,
		}
		classes = append(classes, class)
	}
	return classes, nil
}
