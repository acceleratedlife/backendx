package main

import (
	"encoding/json"
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

func getClassAtSchoolTx(tx *bolt.Tx, schoolId, classId string) (classBucket *bolt.Bucket, err error) {

	school, err := SchoolByIdTx(tx, schoolId)
	if err != nil {
		return nil, err
	}
	teachers := school.Bucket([]byte(KeyTeachers))
	if teachers == nil {
		return nil, fmt.Errorf("no teachers at school")
	}
	cTeachers := teachers.Cursor()
	for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
		if v != nil {
			continue
		}
		teacher := teachers.Bucket(k)
		classBucket = teacher.Bucket([]byte(classId)) //found the class
		if classBucket == nil {
			continue
		}
		return classBucket, nil
	}
	return nil, fmt.Errorf("class not found")
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

func getClassMembership(db *bolt.DB, schoolID, userId string) (classes []openapi.Class, err error) {
	classes = make([]openapi.Class, 0)
	err = db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolID)
		if err != nil {
			return err
		}
		cSchool := school.Cursor()
		for kSchool, v := cSchool.First(); kSchool != nil; kSchool, v = cSchool.Next() { //iterate school bucket
			if v != nil {
				continue
			}
			if string(kSchool) == KeyClasses { //find the school classes bucket
				admins := school.Bucket([]byte(KeyAdmins))
				cAdmins := admins.Cursor()
				ownerKey, _ := cAdmins.First()
				schoolClass, err := getClassMembershipHelper(school, kSchool, userId, ownerKey)
				if err != nil {
					return err
				}
				if schoolClass.OwnerId == "" {
					continue
				}
				classes = append(classes, schoolClass)
			}
			if string(kSchool) == KeyTeachers { //find the teachers bucket
				teachers := school.Bucket(kSchool)
				cTeachers := teachers.Cursor()
				for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
					if v != nil {
						continue
					}
					teacher := teachers.Bucket(k)
					cTeacher := teacher.Cursor()
					for kTeacher, v := cTeacher.First(); kTeacher != nil; kTeacher, v = cTeacher.Last() {
						if v != nil {
							continue
						}
						teacherClass, err := getClassMembershipHelper(teacher, kTeacher, userId, k)
						if err != nil {
							return err
						}
						if teacherClass.OwnerId == "" {
							continue
						}
						classes = append(classes, teacherClass)
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return classes, nil
}

func getClassMembershipHelper(schoolOrTeacher *bolt.Bucket, classIdKey []byte, userId string, ownerIdKey []byte) (class openapi.Class, err error) {
	classBucket := schoolOrTeacher.Bucket(classIdKey)
	studentsBucket := classBucket.Bucket([]byte(KeyStudents))
	if studentsBucket == nil {
		return class, nil
	}
	student := studentsBucket.Get([]byte(userId))
	if student == nil {
		return class, nil
	}
	classByte := schoolOrTeacher.Get(classIdKey)
	err = json.Unmarshal(classByte, &class)
	if err != nil {
		return class, nil
	}
	class.OwnerId = string(ownerIdKey)
	// class := openapi.Class {
	// 	Id: string(kClass),
	// 	// OwnerId: string(class.Get([]byte("ownerId"))),
	// 	Period: btoi32(classBucket.Get([]byte(KeyPeriod))),
	// 	Name: string(classBucket.Get([]byte(KeyName))),
	// 	AddCode: string(classBucket.Get([]byte(KeyAddCode))),
	// 	Members: []string(classBucket.Get([]byte(KeyStudents))),
	// }

	return class, nil
}

func getClassOwner(db *bolt.DB, classId, schoolId string) (user UserInfo, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		school, err := SchoolByIdTx(tx, schoolId)
		if err != nil {
			return err
		}
		cSchool := school.Cursor()
		for kSchool, v := cSchool.First(); kSchool != nil; kSchool, v = cSchool.Next() { //iterate school bucket
			if v != nil {
				continue
			}
			if string(kSchool) == KeyClasses { //find the school classes bucket
				classes := school.Bucket(kSchool)
				class := classes.Get([]byte(classId))
				if class == nil {
					continue
				}
				admins := school.Bucket([]byte(KeyAdmins))
				cAdmins := admins.Cursor()
				ownerKey, _ := cAdmins.First()
				ownerId := string(ownerKey)
				user, err = getUserInLocalStore(db, ownerId)
				if err != nil {
					return err
				}
				return nil
			}
			if string(kSchool) == KeyTeachers { //find the teachers bucket
				teachers := school.Bucket(kSchool)
				cTeachers := teachers.Cursor()
				for k, v := cTeachers.First(); k != nil; k, v = cTeachers.Next() { //iterate the teachers
					if v != nil {
						continue
					}
					teacher := teachers.Bucket(k)
					class := teacher.Get([]byte(classId))
					if class == nil {
						continue
					}
					user, err = getUserInLocalStore(db, string(k))
					if err != nil {
						return err
					}
					return nil
				}
			}
		}
		return nil
	})
	if err != nil {
		return user, err
	}
	return user, nil
}
