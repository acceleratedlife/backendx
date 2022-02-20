package main

import (
	"fmt"

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
