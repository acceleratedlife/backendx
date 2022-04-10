package main

import (
	"context"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

type teacherApiServiceImpl struct {
	db *bolt.DB
}

func (s *teacherApiServiceImpl) TeacherAddClass(ctx context.Context, body openapi.RequestAddClass) (openapi.ImplResponse, error) {
	//TODO implement me
	//This endpoint seems to server no purpose. Front end needs to be changed. Teacher should no longer add admin but add school as the addCode is now tied to the school.
	panic("implement me")
}

func NewTeacherApiServiceImpl(db *bolt.DB) openapi.TeacherApiServicer {
	return &teacherApiServiceImpl{
		db: db,
	}
}
