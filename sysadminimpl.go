package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	crand "crypto/rand"

	openapi "github.com/acceleratedlife/backend/go"
	"github.com/go-pkgz/auth/token"
	"github.com/go-pkgz/lgr"
	"github.com/golang-jwt/jwt"
	bolt "go.etcd.io/bbolt"
)

func FindOrCreateSchool(db *bolt.DB, clock Clock, name string, city string, zip int) (id string, err error) {

	err = db.Update(func(tx *bolt.Tx) error {
		schools, err := tx.CreateBucketIfNotExists([]byte("schools"))
		if err != nil {
			return fmt.Errorf("cannot create 'schools': %v", err)
		}

		c := schools.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			school := schools.Bucket(k)

			if string(school.Get([]byte("name"))) != name {
				continue
			}
			if string(school.Get([]byte("city"))) != city {
				continue
			}
			if btoi(school.Get([]byte("zip"))) != zip {
				continue
			}
			id = string(k)
			return nil
		}

		id = RandomString(10)
		school, err := schools.CreateBucket([]byte(id))
		if err != nil {
			return fmt.Errorf("cannot create bucket for new school: %v", err)
		}

		err = school.Put([]byte(KeyName), []byte(name))
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeyCity), []byte(city))
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeyZip), itob(zip))
		if err != nil {
			return err
		}

		addCodes, _ := constSlice()
		addCode := randomWords(2, 100, addCodes)
		free, err := freeAddCodeRx(tx, addCode)
		if err != nil {
			return err
		}
		for !free {
			addCode = randomWords(2, 100, addCodes)
			free, err = freeAddCodeRx(tx, addCode)
			if err != nil {
				return err
			}
		}

		err = school.Put([]byte(KeyAddCode), []byte(addCode))
		if err != nil {
			return err
		}

		endTime := clock.Now().Add(time.Hour * 72).Truncate(time.Second).Format(time.RFC3339)

		err = school.Put([]byte(KeyRegEnd), []byte(endTime))
		if err != nil {
			return err
		}

		settings := openapi.Settings{
			Student2student: false,
			Lottery:         false,
			Odds:            500,
		}

		marshal, err := json.Marshal(settings)
		if err != nil {
			return err
		}

		err = school.Put([]byte(KeySettings), marshal)
		if err != nil {
			return err
		}

		classes, err := school.CreateBucket([]byte(KeyClasses))
		if err != nil {
			return err
		}

		_, err = school.CreateBucket([]byte(KeyAuctions))
		if err != nil {
			return err
		}

		def := []struct {
			name   string
			period int
		}{
			{name: KeyFreshman, period: 9},
			{name: KeySophomores, period: 10},
			{name: KeyJuniors, period: 11},
			{name: KeySeniors, period: 12},
		}

		for _, c := range def {
			_, err = addClassDetailsTx(tx, classes, clock, c.name, c.period, true)
			if err != nil {
				return err
			}
		}

		cb, err := school.CreateBucket([]byte(KeyCB))
		if err != nil {
			return err
		}
		_, err = cb.CreateBucket([]byte(KeyAccounts))
		if err != nil {
			return err
		}

		return nil
	})

	return id, err
}

func schoolsByZip(db *bolt.DB, zip int32) ([]openapi.ResponseSchoolsInner, error) {
	res := make([]openapi.ResponseSchoolsInner, 0)

	err := db.View(func(tx *bolt.Tx) error {
		schools := tx.Bucket([]byte("schools"))
		if schools == nil {
			return nil
		}

		c := schools.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				continue
			}

			school := schools.Bucket(k)

			if btoi32(school.Get([]byte("zip"))) != zip {
				continue
			}
			schoolId := string(k)
			schoolName := string(school.Get([]byte("name")))
			res = append(res, openapi.ResponseSchoolsInner{
				Name: schoolName,
				Id:   schoolId,
			})
			return nil
		}
		return nil
	})

	return res, err
}

// opens a db.view to pass to getSchoolUsersRx
func getSchoolUsers(db *bolt.DB, schoolId string) (resp []openapi.UserNoHistory, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getSchoolUsersRx(tx, schoolId)
		return err
	})

	return
}

func makeToken(jwtSvc *token.Service, target UserInfo) (string, string, error) {

	tgt := token.User{
		Name: target.Name,
	}

	const ttl = 15 * time.Minute

	// XSRF = random 32-byte URL-safe string
	xsrfBytes := make([]byte, 32)
	if _, err := crand.Read(xsrfBytes); err != nil {
		return "", "", err
	}
	xsrf := base64.RawURLEncoding.EncodeToString(xsrfBytes)

	claims := token.Claims{
		StandardClaims: jwt.StandardClaims{
			Id:        xsrf,                       // XSRF token goes into jti / id
			ExpiresAt: time.Now().Add(ttl).Unix(), // custom TTL, e.g. 15 min
			Issuer:    jwtSvc.Opts.Issuer,         // keep same issuer
		},
		User: &tgt,
	}

	jwtStr, err := jwtSvc.Token(claims) // sign only, no cookies written
	if err != nil {
		return "", "", err
	}
	return jwtStr, xsrf, nil
}

// gets all the users in a school
func getSchoolUsersRx(tx *bolt.Tx, schoolId string) (resp []openapi.UserNoHistory, err error) {
	school, err := schoolByIdTx(tx, schoolId)
	if err != nil {
		return
	}

	users := make([]string, 0)

	adminsBucket := school.Bucket([]byte(KeyAdmins))
	c := adminsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	teacherBucket := school.Bucket([]byte(KeyTeachers))
	c = teacherBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	studentsBucket := school.Bucket([]byte(KeyStudents))
	c = studentsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	//get all the users in the school
	for _, user := range users {
		userInfo, err := getUserInLocalStoreTx(tx, user)
		if err != nil {
			continue
		}

		rank := userInfo.Rank
		if rank == 0 {
			rank = 99999 // put zero-ranked users at the end
		}

		resp = append(resp, openapi.UserNoHistory{
			Id:            userInfo.Name,
			FirstName:     userInfo.FirstName,
			LastName:      userInfo.LastName,
			Role:          userInfo.Role,
			Rank:          rank, //if rank is 0 set it to 99999 so it is at the bottom of the list, this is to fix the sorting issue
			TaxableIncome: userInfo.TaxableIncome,
			Income:        userInfo.Income,
			College:       userInfo.College,
			NetWorth:      userInfo.NetWorth,
			LottoWin:      userInfo.LottoWin,
			LottoPlay:     userInfo.LottoPlay,
		})
	}

	return

}

// opens a db.view to pass to getSchoolsRx
func getSchools(db *bolt.DB) (resp []openapi.ResponseSchools, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		resp, err = getSchoolsRx(tx)
		return err
	})

	return
}

// opens a db.view to pass to getSchoolsRx
func getSchoolsRx(tx *bolt.Tx) (resp []openapi.ResponseSchools, err error) {
	schools := tx.Bucket([]byte("schools"))
	if schools == nil {
		return nil, fmt.Errorf("no schools available")
	}

	c := schools.Cursor()

	for k, v := c.First(); k != nil; k, v = c.Next() {
		if v != nil {
			continue
		}

		school := schools.Bucket(k)
		isPaused := school.Get([]byte(KeyPaused))

		item := openapi.ResponseSchools{
			Id:       string(k),
			Name:     string(school.Get([]byte(KeyName))),
			City:     string(school.Get([]byte(KeyCity))),
			Zip:      btoi32(school.Get([]byte(KeyZip))),
			IsPaused: isPaused != nil,
		}

		//loop through each bucket and count how many keys are in each
		adminsBucket := school.Bucket([]byte(KeyAdmins))
		c := adminsBucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Staff++
		}

		teacherBucket := school.Bucket([]byte(KeyTeachers))
		c = teacherBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Staff++
		}

		studentsBucket := school.Bucket([]byte(KeyStudents))
		c = studentsBucket.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			item.Students++
		}

		resp = append(resp, item)
	}
	return
}

func message(db *bolt.DB, message, userId, schoolId *string, students, staff bool) error {
	return db.Update(func(tx *bolt.Tx) error {
		return messageTx(tx, message, userId, schoolId, students, staff)
	})
}

func messageTx(tx *bolt.Tx, message, userId, schoolId *string, students, staff bool) error {
	users := tx.Bucket([]byte(KeyUsers))
	if users == nil {
		return fmt.Errorf("users not found")
	}

	//a message to a single user
	if userId != nil {
		userData := users.Get([]byte(*userId))
		if userData == nil {
			return fmt.Errorf("user not found")
		}
		var user UserInfo
		err := json.Unmarshal(userData, &user)
		if err != nil {
			return err
		}

		user.Messages = append(user.Messages, *message)
		marshal, err := json.Marshal(user)
		if err != nil {
			return err
		}

		err = users.Put([]byte(*userId), marshal)
		if err != nil {
			return err
		}

		return nil
	}

	c := users.Cursor()

	if schoolId == nil {

		//a message to all users
		if students && staff {

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				userData := users.Get([]byte(k))
				var user UserInfo
				err := json.Unmarshal(userData, &user)
				if err != nil {
					lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
					continue
				}

				if user.Role == UserRoleSysAdmin {
					continue
				}

				//I might need to prevent admin teacher accounts from getting messages here

				user.Messages = append(user.Messages, *message)
				marshal, err := json.Marshal(user)
				if err != nil {
					return err
				}

				err = users.Put([]byte(k), marshal)
				if err != nil {
					return err
				}
			}
			return nil
		}

		//a message to all students
		if students {

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				userData := users.Get([]byte(k))
				var user UserInfo
				err := json.Unmarshal(userData, &user)
				if err != nil {
					lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
					continue
				}

				if user.Role != UserRoleStudent {
					continue
				}

				user.Messages = append(user.Messages, *message)
				marshal, err := json.Marshal(user)
				if err != nil {
					return err
				}

				err = users.Put([]byte(k), marshal)
				if err != nil {
					return err
				}
			}
			return nil
		}

		//a message to all staff
		if staff {

			for k, _ := c.First(); k != nil; k, _ = c.Next() {
				userData := users.Get([]byte(k))
				var user UserInfo
				err := json.Unmarshal(userData, &user)
				if err != nil {
					lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
					continue
				}

				if user.Role == UserRoleSysAdmin {
					continue
				}

				if user.Role == UserRoleStudent {
					continue
				}

				user.Messages = append(user.Messages, *message)
				marshal, err := json.Marshal(user)
				if err != nil {
					return err
				}

				err = users.Put([]byte(k), marshal)
				if err != nil {
					return err
				}
			}
			return nil
		}
	}

	//a message to all users of a school
	if students && staff {

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			userData := users.Get([]byte(k))
			var user UserInfo
			err := json.Unmarshal(userData, &user)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}

			if user.Role == UserRoleSysAdmin {
				continue
			}

			if user.SchoolId != *schoolId {
				continue
			}

			user.Messages = append(user.Messages, *message)
			marshal, err := json.Marshal(user)
			if err != nil {
				return err
			}

			err = users.Put([]byte(k), marshal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	//a message to all students of a school
	if students {

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			userData := users.Get([]byte(k))
			var user UserInfo
			err := json.Unmarshal(userData, &user)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}

			if user.Role == UserRoleSysAdmin {
				continue
			}

			if user.SchoolId != *schoolId {
				continue
			}

			if user.Role != UserRoleStudent {
				continue
			}

			user.Messages = append(user.Messages, *message)
			marshal, err := json.Marshal(user)
			if err != nil {
				return err
			}

			err = users.Put([]byte(k), marshal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	//a message to all staff of a school
	if staff {

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			userData := users.Get([]byte(k))
			var user UserInfo
			err := json.Unmarshal(userData, &user)
			if err != nil {
				lgr.Printf("ERROR cannot unmarshal userInfo for %s", k)
				continue
			}

			if user.Role == UserRoleSysAdmin {
				continue
			}

			if user.SchoolId != *schoolId {
				continue
			}

			if user.Role == UserRoleStudent {
				continue
			}

			user.Messages = append(user.Messages, *message)
			marshal, err := json.Marshal(user)
			if err != nil {
				return err
			}

			err = users.Put([]byte(k), marshal)
			if err != nil {
				return err
			}
		}
		return nil
	}

	return fmt.Errorf("no message sent")
}

func deleteSchool(db *bolt.DB, schoolId string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return deleteSchoolTx(tx, schoolId)
	})
}

func deleteSchoolTx(tx *bolt.Tx, schoolId string) error {
	schools := tx.Bucket([]byte("schools"))
	if schools == nil {
		return fmt.Errorf("schools not found")
	}

	schoolBucket := schools.Bucket([]byte(schoolId))
	if schoolBucket == nil {
		return fmt.Errorf("school not found")
	}

	adminsBucket := schoolBucket.Bucket([]byte(KeyAdmins))
	if adminsBucket == nil {
		return fmt.Errorf("admins not found")
	}

	users := make([]string, 0)

	c := adminsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	teachersBucket := schoolBucket.Bucket([]byte(KeyTeachers))
	if teachersBucket == nil {
		return fmt.Errorf("teachers not found")
	}

	c = teachersBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	studentsBucket := schoolBucket.Bucket([]byte(KeyStudents))
	if studentsBucket == nil {
		return fmt.Errorf("students not found")
	}

	c = studentsBucket.Cursor()
	for k, _ := c.First(); k != nil; k, _ = c.Next() {
		users = append(users, string(k))
	}

	usersBucket := tx.Bucket([]byte(KeyUsers))
	if usersBucket == nil {
		return fmt.Errorf("users not found")
	}

	for _, user := range users {
		err := usersBucket.Delete([]byte(user))
		if err != nil {
			return err
		}
	}

	err := schools.DeleteBucket([]byte(schoolId))
	if err != nil {
		return err
	}

	return nil
}

func togglePause(db *bolt.DB, schoolId string) error {
	return db.Update(func(tx *bolt.Tx) error {
		return togglePauseTx(tx, schoolId)
	})
}

func togglePauseTx(tx *bolt.Tx, schoolId string) (err error) {
	schools := tx.Bucket([]byte("schools"))
	if schools == nil {
		return fmt.Errorf("schools not found")
	}

	schoolBucket := schools.Bucket([]byte(schoolId))
	if schoolBucket == nil {
		return fmt.Errorf("school not found")
	}

	paused := schoolBucket.Get([]byte(KeyPaused))
	if paused == nil {
		err = schoolBucket.Put([]byte(KeyPaused), []byte("true"))
	} else {
		err = schoolBucket.Delete([]byte(KeyPaused))
	}

	return
}

func isSchoolPaused(db *bolt.DB, schoolId string) (isPaused bool, err error) {
	err = db.View(func(tx *bolt.Tx) error {
		isPaused, err = isSchoolPausedTx(tx, schoolId)
		return err
	})
	return isPaused, err
}

func isSchoolPausedTx(tx *bolt.Tx, schoolId string) (isPaused bool, err error) {
	schools := tx.Bucket([]byte("schools"))
	if schools == nil {
		return false, fmt.Errorf("schools not found")
	}

	schoolBucket := schools.Bucket([]byte(schoolId))
	if schoolBucket == nil {
		return false, fmt.Errorf("school not found")
	}

	paused := schoolBucket.Get([]byte(KeyPaused))
	if paused == nil {
		return false, nil
	}
	return true, nil
}
