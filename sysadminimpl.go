package main

import (
	"fmt"

	openapi "github.com/acceleratedlife/backend/go"
	bolt "go.etcd.io/bbolt"
)

func FindOrCreateSchool(db *bolt.DB, name string, city string, zip int) (id string, err error) {

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

		err = school.Put([]byte(KeyAddCode), []byte(RandomString(6)))
		if err != nil {
			return err
		}

		classes, err := school.CreateBucket([]byte(KeyClasses))
		if err != nil {
			return err
		}

		def := []struct {
			name   string
			period int
		}{
			{name: "Freshman", period: 9},
			{name: "Sophomore", period: 10},
			{name: "Junior", period: 11},
			{name: "Senior", period: 12},
		}

		for _, c := range def {
			_, err = addClassDetailsTx(classes, c.name, c.period, true)
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
