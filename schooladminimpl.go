package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/lgr"
	bolt "go.etcd.io/bbolt"
)

func validateUserData(body UserInfo) (bool, error) {
	if body.SchoolId == "" {
		return false, fmt.Errorf("no school specified")
	}
	if body.Name == "" {
		return false, fmt.Errorf(" username is mandatory")
	}
	if len(body.Email) < 3 {
		return false, fmt.Errorf("wrong email format")
	}

	if body.FirstName == "" {
		return false, fmt.Errorf("first name is mandatory")
	}

	if body.LastName == "" {
		return false, fmt.Errorf("last name is mandatory")
	}

	return true, nil
}

func CreateSchoolAdmin(db *bolt.DB, body UserInfo) (openapi.ImplResponse, error) {

	_, err := validateUserData(body)
	if err != nil {
		return openapi.ImplResponse{}, err
	}

	if body.Role != UserRoleAdmin {
		return openapi.ImplResponse{}, fmt.Errorf("not a school admin")
	}
	v := openapi.ImplResponse{}

	err = db.Update(func(tx *bolt.Tx) error {

		newUser := UserInfo{
			Name:        body.Name,
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			Email:       body.Email,
			Confirmed:   true,
			PasswordSha: body.PasswordSha,
			SchoolId:    body.SchoolId,
			Role:        UserRoleAdmin,
		}

		err = AddUserTx(tx, newUser)
		if err != nil {
			return err
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(body.SchoolId))

		if school == nil {
			return fmt.Errorf("school is not created")
		}

		admins, err := school.CreateBucketIfNotExists([]byte(KeyAdmins))
		if err != nil {
			return err
		}

		_, err = school.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}

		admin := admins.Get([]byte(body.Name))

		if admin != nil {
			lgr.Printf("INFO admin %s is already registered ", body.Email)
			return nil
			//fmt.Errorf("school admin %s is already registered", body.Email)
		}

		err = admins.Put([]byte(body.Name), []byte(""))
		if err != nil {
			return err
		}

		return nil
	})

	return v, err
}

func createTeacher(db *bolt.DB, newUser UserInfo) (err error) {
	err = db.Update(func(tx *bolt.Tx) error {

		err = AddUserTx(tx, newUser)
		if err != nil {
			return err
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(newUser.SchoolId))

		if school == nil {
			return fmt.Errorf("school not found")
		}
		teachers, err := school.CreateBucketIfNotExists([]byte(KeyTeachers))
		if err != nil {
			return err
		}
		teacher, err := teachers.CreateBucket([]byte(newUser.Email))
		if err != nil {
			return err
		}
		_, err = teacher.CreateBucket([]byte(KeyClasses))
		if err != nil {
			return nil
		}
		return nil
	})
	return
}

func createStudent(db *bolt.DB, newUser UserInfo, pathId PathId) (err error) {
	newUser.Income = float32(math.Floor(rand.Float64()*(335.0-104.0) + 104.0))
	newUser.CareerTransition = false
	err = db.Update(func(tx *bolt.Tx) error {
		err = AddUserTx(tx, newUser)
		if err != nil {
			// AddUser can fail if new user exists in the table as a student
			users := tx.Bucket([]byte(KeyUsers))
			if users == nil {
				return err
			}
			user := users.Get(userNameToB(newUser.Email))
			if user == nil {
				return err
			}
			var existingUser UserInfo

			err = json.Unmarshal(user, &existingUser)
			if err != nil {
				return err
			}
			if existingUser.Role != UserRoleStudent {
				return fmt.Errorf("you cannot add existing user with the role %d as a student to a class", existingUser.Role)
			}
		}

		schools, err := tx.CreateBucketIfNotExists([]byte(KeySchools))
		if err != nil {
			return err
		}

		school := schools.Bucket([]byte(newUser.SchoolId))

		if school == nil {
			return fmt.Errorf("school not found")
		}
		schoolStudentsBucket, err := school.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}
		schoolStudentBucket, err := schoolStudentsBucket.CreateBucketIfNotExists([]byte(newUser.Email))
		if err != nil {
			return err
		}
		history := []openapi.History{}
		row1 := openapi.History{
			Date:     time.Now(),
			NetWorth: 0.0,
		}
		history = append(history, row1)
		marshal, _ := json.Marshal(history)
		schoolStudentBucket.Put([]byte(KeyHistory), marshal)
		teachers, err := school.CreateBucketIfNotExists([]byte(KeyTeachers))
		if err != nil {
			return err
		}

		teacher := teachers.Bucket([]byte(pathId.teacherId))
		if teacher == nil {
			return fmt.Errorf("teacher not found")
		}
		classes := teacher.Bucket([]byte(KeyClasses))
		if classes == nil {
			return fmt.Errorf("classes not found")
		}
		class := classes.Bucket([]byte(pathId.classId))
		if class == nil {
			return fmt.Errorf("class not found")
		}
		students, err := class.CreateBucketIfNotExists([]byte(KeyStudents))
		if err != nil {
			return err
		}
		student := students.Get([]byte(newUser.Email))
		if student != nil {
			return fmt.Errorf("student is already in the class")
		}
		return students.Put([]byte(newUser.Email), []byte(""))
	})
	return err
}
