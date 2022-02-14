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

func (s *SchoolAdminServiceImpl) SearchAdminTeacherClass(ctx context.Context, s2 string) (openapi.ImplResponse, error) {
	userData := ctx.Value("user").(token.User)
	userDetails, err := getUserInLocalStore(s.db, userData.Name)
	if err != nil {
		return openapi.Response(404, nil), nil
	}

	members := make([]openapi.ClassWithMembersMembers, 0)
	addCode := ""

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
			Id:      "",
			OwnerId: userDetails.Name,
			Period:  0,
			Name:    "teachers at school",
			AddCode: addCode,
			Members: members,
		}), nil

}

func NewSchoolAdminServiceImpl(db *bolt.DB) openapi.SchoolAdminApiServicer {
	return &SchoolAdminServiceImpl{
		db: db,
	}
}
