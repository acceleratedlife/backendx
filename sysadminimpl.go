package main

import (
	"encoding/json"
	"fmt"
	"time"

	openapi "github.com/acceleratedlife/backend/go"
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
			Odds:            5000,
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

func schoolsByZip(db *bolt.DB, zip int32) ([]openapi.ResponseSchools, error) { //needs to be tested
	res := make([]openapi.ResponseSchools, 0)

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

			schoolZip := btoi32(school.Get([]byte(KeyZip)))

			if zip != 0 && schoolZip != zip {
				continue
			}

			schoolId := string(k)
			schoolName := string(school.Get([]byte(KeyName)))
			schoolCity := string(school.Get([]byte(KeyCity)))
			studentsBucket := school.Bucket([]byte(KeyStudents))

			cursor := studentsBucket.Cursor()
			studentCount := int32(0)
			for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
				studentCount++
			}

			res = append(res, openapi.ResponseSchools{
				Name:     schoolName,
				Id:       schoolId,
				City:     schoolCity,
				Zip:      schoolZip,
				Students: studentCount,
			})
		}
		return nil
	})

	return res, err
}

func CreateSysAdmin(db *bolt.DB, body UserInfo) (openapi.ImplResponse, error) {

	if body.Role != UserRoleSysAdmin {
		return openapi.ImplResponse{}, fmt.Errorf("not a sysAdmin")
	}
	v := openapi.ImplResponse{}

	err := db.Update(func(tx *bolt.Tx) error {

		newUser := UserInfo{
			Name:        body.Name,
			FirstName:   body.FirstName,
			LastName:    body.LastName,
			Email:       body.Email,
			Confirmed:   true,
			PasswordSha: body.PasswordSha,
			SchoolId:    body.SchoolId,
			Role:        UserRoleSysAdmin,
		}

		err := AddUserTx(tx, newUser)
		if err != nil {
			return err
		}

		return nil
	})

	return v, err
}
