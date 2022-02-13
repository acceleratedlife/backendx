package main

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"time"
)

const (
	UserRoleStudent = int32(0)
	UserRoleTeacher = int32(1)
	UserRoleAdmin   = int32(2)
)

type UserInfo struct {
	Name           string
	FirstName      string
	LastName       string
	Email          string
	Confirmed      bool
	PasswordSha    string
	SchoolId       string
	Role           int32     // 0 student, 1 teacher, 2 admin
	Salary         float64   `json:",omitempty"`
	LastSalaryPaid time.Time `json:",omitempty"`
}

type PathId struct {
	schoolId  string
	teacherId string
	classId   string
}

func checkUserInLocalStore(db *bolt.DB, user, password string) (ok bool, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		users := tx.Bucket([]byte("users"))
		if users == nil {
			return fmt.Errorf("user does not exist")
		}

		userBuf := users.Get([]byte(user))

		if userBuf == nil {
			return fmt.Errorf("user does not exist")
		}

		var user UserInfo
		err := json.Unmarshal(userBuf, &user)
		if err != nil {
			return err
		}

		ok = EncodePassword(password) == user.PasswordSha
		if !ok {
			return fmt.Errorf("password does not match")
		}
		return nil
	})

	return
}

func getUserInLocalStore(db *bolt.DB, userId string) (user UserInfo, err error) {

	err = db.View(func(tx *bolt.Tx) error {
		user, err = getUserInLocalStoreTx(tx, userId)
		return err
	})

	return
}

func getUserInLocalStoreTx(tx *bolt.Tx, userId string) (user UserInfo, err error) {
	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		err = fmt.Errorf("user does not exist")
		return
	}

	userBuf := users.Get([]byte(userId))

	if userBuf == nil {
		err = fmt.Errorf("user does not exist")
		return
	}

	err = json.Unmarshal(userBuf, &user)
	return
}

// AddUserTx adds register info
func AddUserTx(tx *bolt.Tx, info UserInfo) error {
	users, err := tx.CreateBucketIfNotExists([]byte(KeyUsers))
	if err != nil {
		return err
	}

	user := users.Get(userNameToB(info.Email))

	if user != nil {
		return fmt.Errorf("user already exists - %s", info.Email)
	}

	userBuf, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = users.Put(userNameToB(info.Email), userBuf)
	if err != nil {
		return err
	}

	return nil
}

func getStudentBucketTx(tx *bolt.Tx, userName string) (*bolt.Bucket, error) {
	userInfo, err := getUserInLocalStoreTx(tx, userName)
	if err != nil {
		return nil, err
	}

	schools := tx.Bucket([]byte(KeySchools))
	if schools == nil {
		return nil, fmt.Errorf("schools  does not exist")
	}

	school := schools.Bucket([]byte(userInfo.SchoolId))

	students, err := school.CreateBucketIfNotExists([]byte(KeyStudents))
	if err != nil {
		return nil, err
	}

	student, err := students.CreateBucketIfNotExists([]byte(userName))
	if err != nil {
		return nil, err
	}
	return student, nil
}
