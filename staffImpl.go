package main

import (
	"fmt"
	openapi "github.com/acceleratedlife/backend/go"
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
		teachers := school.Bucket([]byte("teachers"))
		if teachers == nil {
			return nil
		}
		teacher := teachers.Bucket([]byte(teacherId))
		if teacher == nil {
			return nil
		}
		res = getClasses1Tx(teacher, teacherId)
		return nil
	})
	return
}

func getClassesTx(teacher *bolt.Bucket) []openapi.ResponseMakeClassInner {
	classes := make([]openapi.ResponseMakeClassInner, 0)

	c := teacher.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}
		classBucket := teacher.Bucket(k)
		iClass := openapi.ResponseMakeClassInner{
			Id:      string(k),
			Name:    string(classBucket.Get([]byte(KeyName))),
			Owner:   "",
			Period:  btoi32(classBucket.Get([]byte("period"))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: make([]string, 0),
		}
		classes = append(classes, iClass)
	}

	return classes
}

func getClasses1Tx(teacher *bolt.Bucket, ownerId string) []openapi.Class {
	data := make([]openapi.Class, 0)
	c := teacher.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}
		classBucket := teacher.Bucket(k)
		iClass := openapi.Class{
			Id:      string(k),
			OwnerId: ownerId,
			Period:  btoi32(classBucket.Get([]byte("period"))),
			Name:    string(classBucket.Get([]byte(KeyName))),
			AddCode: string(classBucket.Get([]byte(KeyAddCode))),
			Members: make([]string, 0),
		}
		data = append(data, iClass)
	}
	return data
}

func (s *StaffApiServiceImpl) MakeClassImpl(userDetails UserInfo, request openapi.RequestMakeClass) (classes []openapi.ResponseMakeClassInner, err error) {
	schoolId := userDetails.SchoolId
	teacherId := userDetails.Name
	className := request.Name
	period := request.Period

	_, classes, err = CreateClass(s.db, schoolId, teacherId, className, int(period))

	return classes, err
}

func CreateClass(db *bolt.DB, schoolId, teacherId, className string, period int) (classId string, classes []openapi.ResponseMakeClassInner, err error) {
	err = db.Update(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolId)
		if err != nil {
			return err
		}
		teachers := school.Bucket([]byte("teachers"))
		if teachers == nil {
			return fmt.Errorf("user does not exist")
		}
		teacher := teachers.Bucket([]byte(teacherId))
		if teacher == nil {
			return fmt.Errorf("user does not exist")
		}

		classId, err = addClassDetailsTx(teacher, className, period)
		if err != nil {
			return err
		}

		classes = getClassesTx(teacher)

		return nil
	})
	return
}

func addClassDetailsTx(bucket *bolt.Bucket, className string, period int) (classId string, err error) {
	classId = RandomString(15)

	class, err1 := bucket.CreateBucket([]byte(classId))
	if err1 != nil {
		return "", err1
	}

	err = class.Put([]byte(KeyName), []byte(className))
	if err != nil {
		return "", err
	}
	err = class.Put([]byte("period"), itob32(int32(period)))
	if err != nil {
		return "", err
	}
	addCode := RandomString(6)
	err = class.Put([]byte(KeyAddCode), []byte(addCode))
	return
}
