package main

import (
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

func UserByIdTx(tx *bolt.Tx, userId string) (user *bolt.Bucket, err error) {
	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		return nil, fmt.Errorf("no users available")
	}

	c := users.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {

		if string([]byte(k)) == userId {
			println(k, v)
		}

	}
	return nil, fmt.Errorf("school not found")

}

func ClassForAll(tx *bolt.Tx, classId string) (classBucket *bolt.Bucket, err error) {
	schools := tx.Bucket([]byte(KeySchools))
	cSchools := schools.Cursor()
	for k, v := cSchools.First(); k != nil; k, v = cSchools.Next() { //iterate through all schools
		if v != nil {
			continue
		}
		school := schools.Bucket(k)
		if school == nil {
			continue
		}
		cSchool := school.Cursor()
		for key, v := cSchool.First(); key != nil; key, v = cSchool.Next() { //iterate school bucket
			if v != nil {
				continue
			}
			if string(key) == KeyClasses { //find the school classes bucket
				classes := school.Bucket(key)
				classBucket := classes.Bucket([]byte(classId)) //found the class
				if classBucket == nil {
					continue
				}
				return classBucket, nil
			}
			if string(key) == KeyTeachers { //find the teachers bucket
				teachers := school.Bucket(key)
				cTeachers := teachers.Cursor()
				for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
					if v != nil {
						continue
					}
					teacher := teachers.Bucket(k)
					classBucket := teacher.Bucket([]byte(classId)) //found the class
					if classBucket == nil {
						continue
					}
					return classBucket, nil
				}
			}
		}
	}
	return nil, fmt.Errorf("class not found")
}

func PopulateClassMembers(tx *bolt.Tx, classBucket *bolt.Bucket) (Members []openapi.ClassWithMembersMembers, err error) {
	Members = make([]openapi.ClassWithMembersMembers, 0)
	students := classBucket.Bucket([]byte("students"))
	cStudents := students.Cursor()
	for k, _ := cStudents.First(); k != nil; k, _ = cStudents.Next() { //iterate students bucket
		user, err := getUserInLocalStoreTx(tx, string(k))
		if err != nil {
			return nil, err
		}

		nWorth, _ := StudentNetWorthTx(tx, user.Email).Float64()
		nUser := openapi.ClassWithMembersMembers{
			Id:        user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Rank:      2,
			NetWorth:  float32(nWorth),
		}
		Members = append(Members, nUser)
	}
	return Members, nil
}
