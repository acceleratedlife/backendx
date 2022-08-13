package main

import (
	"context"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

type SchoolAdminServiceImpl struct {
	db *bolt.DB
}

func (s *SchoolAdminServiceImpl) SearchAdminTeacherClass(ctx context.Context, Id string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, nil), nil
	}

	members := make([]openapi.ClassWithMembersMembers, 0)
	addCode := ""

	if userDetails.Role != UserRoleAdmin {
		return openapi.Response(401, ""), nil
	}

	_ = s.db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, userDetails.SchoolId)
		if err != nil {
			return err
		}

		addCodeB := school.Get([]byte(KeyAddCode))
		if addCodeB != nil {
			addCode = string(addCodeB)
		}

		teachers := school.Bucket([]byte(KeyTeachers))

		if teachers == nil {
			return nil
		}

		t := teachers.Cursor()
		for k, v := t.First(); k != nil; k, v = t.Next() {
			if v != nil {
				continue
			}
			user, err := getUserInLocalStoreTx(tx, string(k))
			if err != nil {
				lgr.Printf("ERROR teacher not found in users-bucket")
				continue
			}
			members = append(members, openapi.ClassWithMembersMembers{
				LastName: user.LastName,
				Id:       string(k),
			})
		}
		return nil

	})

	return openapi.Response(200,
		openapi.ClassWithMembers{
			Id:      userDetails.SchoolId,
			OwnerId: userDetails.Name,
			Period:  -1,
			Name:    "My Teachers",
			AddCode: addCode,
			Members: members,
		}), nil

}

func (s *SchoolAdminServiceImpl) GetSettings(ctx context.Context) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, nil), nil
	}

	if userDetails.Role == UserRoleStudent {
		return openapi.Response(401, ""), nil
	}

	settings, err := getSettings(s.db, userDetails)
	if err != nil {
		lgr.Printf("ERROR cannot get settings with : %s %v", userDetails.Name, err)
		return openapi.Response(500, "{}"), err
	}

	return openapi.Response(200, settings), nil

}

func (s *SchoolAdminServiceImpl) SetSettings(ctx context.Context, body openapi.Settings) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(400, nil), nil
	}

	if userDetails.Role != UserRoleAdmin {
		return openapi.Response(401, ""), nil
	}

	err = setSettings(s.db, userDetails, body)
	if err != nil {
		lgr.Printf("ERROR cannot set: %v", err)
		return openapi.Response(500, "{}"), err
	}

	return openapi.Response(200, nil), nil

}

func NewSchoolAdminServiceImpl(db *bolt.DB) openapi.SchoolAdminApiServicer {
	return &SchoolAdminServiceImpl{
		db: db,
	}
}
